package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	common "github.com/kubeflow/common/job_controller/api/v1"
	pyv1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	torchjobclient "github.com/kubeflow/pytorch-operator/pkg/client/clientset/versioned"
	"github.com/kubeflow/pytorch-operator/pkg/util"
	"github.com/kubeflow/tf-operator/pkg/common/jobcontroller"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	name       = flag.String("name", "", "The name for the PyTorchJob to create.")
	namespace  = flag.String("namespace", "kubeflow", "The namespace to create the test job in.")
	numJobs    = flag.Int("num_jobs", 1, "The number of jobs to run.")
	timeout    = flag.Duration("timeout", 10*time.Minute, "The timeout for the test")
	image      = flag.String("image", "", "The Test image to run")
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
)

func getReplicaSpec(worker int32) map[pyv1.PyTorchReplicaType]*common.ReplicaSpec {
	spec := make(map[pyv1.PyTorchReplicaType]*common.ReplicaSpec)
	spec[pyv1.PyTorchReplicaTypeMaster] = replicaSpec(1)
	spec[pyv1.PyTorchReplicaTypeWorker] = replicaSpec(worker)
	return spec
}

func replicaSpec(replica int32) *common.ReplicaSpec {
	return &common.ReplicaSpec{
		Replicas: proto.Int32(replica),
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:            "pytorch",
						Image:           *image,
						ImagePullPolicy: "IfNotPresent",
					},
				},
			},
		},
	}
}

func hasCondition(status common.JobStatus, condType common.JobConditionType) bool {
	for _, condition := range status.Conditions {
		if condition.Type == condType && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func isSucceeded(status common.JobStatus) bool {
	return hasCondition(status, common.JobSucceeded)
}

func isFailed(status common.JobStatus) bool {
	return hasCondition(status, common.JobFailed)
}

func run() (string, error) {
	jobName := *name
	if jobName == "" {
		jobName = "nil-cleanpolicy-test-job"
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	if *image == "" {
		log.Fatalf("--image must be provided.")
	}

	client := kubernetes.NewForConfigOrDie(config)

	torchJobClient, err := torchjobclient.NewForConfig(config)
	if err != nil {
		return "", err
	}

	// Create a PyTorchJob WITHOUT setting CleanPodPolicy.
	// This tests that:
	// 1. The controller does not panic when CleanPodPolicy is nil (nil guards).
	// 2. Scheme defaulting sets nil to CleanPodPolicyNone, preserving pods
	//    after job completion.
	original := &pyv1.PyTorchJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: pyv1.PyTorchJobSpec{
			PyTorchReplicaSpecs: getReplicaSpec(3),
		},
	}

	_, err = torchJobClient.KubeflowV1().PyTorchJobs(*namespace).Create(original)
	if err != nil {
		log.Errorf("Creating the job failed: %v", err)
		return jobName, err
	}
	log.Infof("Job created (CleanPodPolicy unset): \n%v", util.Pformat(original))

	var torchJob *pyv1.PyTorchJob
	for endTime := time.Now().Add(*timeout); time.Now().Before(endTime); {
		torchJob, err = torchJobClient.KubeflowV1().PyTorchJobs(*namespace).Get(jobName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("Getting PyTorchJob %q: %v", jobName, err)
			return jobName, err
		}

		if isSucceeded(torchJob.Status) || isFailed(torchJob.Status) {
			log.Infof("job %v finished:\n%v", jobName, util.Pformat(torchJob))
			break
		}
		log.Infof("Waiting for job %v to finish", jobName)
		time.Sleep(5 * time.Second)
	}

	if torchJob == nil {
		return jobName, fmt.Errorf("Failed to get PyTorchJob %v", jobName)
	}

	if !isSucceeded(torchJob.Status) {
		return jobName, fmt.Errorf("PyTorchJob %v did not succeed;\n %v", jobName, util.Pformat(torchJob))
	}

	// Verify pods are preserved after completion.
	// Since scheme defaulting sets nil CleanPodPolicy to CleanPodPolicyNone,
	// the controller should NOT delete any pods.
	for rtype, r := range original.Spec.PyTorchReplicaSpecs {
		for i := 0; i < int(*r.Replicas); i++ {
			podName := jobcontroller.GenGeneralName(torchJob.Name, strings.ToLower(string(rtype)), strconv.Itoa(i))
			_, err := client.CoreV1().Pods(*namespace).Get(podName, metav1.GetOptions{})
			if err != nil {
				return jobName, fmt.Errorf("Pod %v for ReplicaType %v Index %v was deleted despite nil CleanPodPolicy (expected preserved): %v", podName, rtype, i, err)
			}
		}
	}
	log.Infof("All pods preserved after job completion (nil CleanPodPolicy defaulted to None)")

	if err := torchJobClient.KubeflowV1().PyTorchJobs(*namespace).Delete(jobName, &metav1.DeleteOptions{}); err != nil {
		return jobName, fmt.Errorf("Failed to delete PyTorchJob %v: %v", jobName, err)
	}

	deleted := false
	for endTime := time.Now().Add(*timeout); time.Now().Before(endTime); {
		_, err = torchJobClient.KubeflowV1().PyTorchJobs(*namespace).Get(jobName, metav1.GetOptions{})
		if k8s_errors.IsNotFound(err) {
			deleted = true
			break
		} else if err != nil {
			log.Errorf("Getting PyTorchJob %q for deletion check: %v", jobName, err)
		} else {
			log.Infof("Job %v still exists", jobName)
		}
		time.Sleep(5 * time.Second)
	}

	if !deleted {
		return jobName, fmt.Errorf("Deletion of PyTorchJob %v failed", jobName)
	}
	return jobName, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func main() {
	flag.Parse()

	if *kubeconfig == "" {
		if kubeEnv := os.Getenv("KUBECONFIG"); kubeEnv != "" {
			*kubeconfig = kubeEnv
		} else if home := homeDir(); home != "" {
			*kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	type Result struct {
		Error error
		Name  string
	}
	c := make(chan Result, *numJobs)

	for i := 0; i < *numJobs; i++ {
		go func() {
			name, err := run()
			if err != nil {
				log.Errorf("Job %v didn't run successfully: %v", name, err)
			} else {
				log.Infof("Job %v ran successfully", name)
			}
			c <- Result{
				Name:  name,
				Error: err,
			}
		}()
	}

	numSucceeded := 0
	numFailed := 0

	for endTime := time.Now().Add(*timeout); numSucceeded+numFailed < *numJobs && time.Now().Before(endTime); {
		select {
		case res := <-c:
			if res.Error == nil {
				numSucceeded += 1
			} else {
				numFailed += 1
			}
		case <-time.After(time.Until(endTime)):
			log.Errorf("Timeout waiting for PyTorchJob to finish.")
		}
	}

	if numSucceeded+numFailed < *numJobs {
		log.Errorf("Timeout waiting for jobs to finish; only %v of %v PyTorchJobs completed.", numSucceeded+numFailed, *numJobs)
	}

	if numSucceeded == *numJobs {
		fmt.Println("Successfully ran PyTorchJob with nil CleanPodPolicy")
	} else {
		fmt.Printf("Running PyTorchJobs failed \n")
		os.Exit(1)
	}
}
