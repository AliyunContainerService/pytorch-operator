// Copyright 2024 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pytorch

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kubeflow/tf-operator/pkg/common/jobcontroller"
	corev1 "k8s.io/api/core/v1"

	v1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
)

const (
	EnvMasterAddr       = "MASTER_ADDR"
	EnvMasterPort       = "MASTER_PORT"
	EnvWorldSize        = "WORLD_SIZE"
	EnvRank             = "RANK"
	EnvPythonUnbuffered = "PYTHONUNBUFFERED"

	// EnvPetMasterAddr is the environment variable name for the master address.
	EnvPetMasterAddr = "PET_MASTER_ADDR"
	// EnvPetMasterPort is the environment variable name for the master port.
	EnvPetMasterPort = "PET_MASTER_PORT"
	// EnvPetNprocPerNode is the environment variable name for the number of processes per node.
	EnvPetNprocPerNode = "PET_NPROC_PER_NODE"
	// EnvPetNnodes is the environment variable name for the number of nodes.
	EnvPetNnodes = "PET_NNODES"
	// EnvPetNodeRank is the environment variable name for the rank of nodes.
	EnvPetNodeRank = "PET_NODE_RANK"
)

// setPodTemplateEnv sets the environment variables for the podTemplateSpec.
func setPodTemplateEnv(podTemplateSpec *corev1.PodTemplateSpec, job *v1.PyTorchJob, rtype v1.PyTorchReplicaType, index string) error {
	rank, err := strconv.Atoi(index)
	if err != nil {
		return err
	}

	masterAddr := jobcontroller.GenGeneralName(job.Name, strings.ToLower(string(v1.PyTorchReplicaTypeMaster)), strconv.Itoa(0))
	if rtype == v1.PyTorchReplicaTypeMaster {
		if rank != 0 {
			return errors.New("invalid config: There should be only a single master with index=0")
		}
	} else {
		rank = rank + 1
	}

	masterPort, err := GetPortFromPyTorchJob(job, v1.PyTorchReplicaTypeMaster)
	if err != nil {
		return err
	}

	worldSize := getWorldSize(job)
	nnodes := int(getTotalReplicas(job))

	for i := range podTemplateSpec.Spec.Containers {
		container := &podTemplateSpec.Spec.Containers[i]
		setContainerEnvIfNotFound(container, EnvMasterAddr, masterAddr)
		setContainerEnvIfNotFound(container, EnvMasterPort, strconv.Itoa(int(masterPort)))
		setContainerEnvIfNotFound(container, EnvWorldSize, strconv.Itoa(worldSize))
		setContainerEnv(container, EnvRank, strconv.Itoa(rank))
		setContainerEnvIfNotFound(container, EnvPythonUnbuffered, "1")
		setContainerEnvIfNotFound(container, EnvPetMasterAddr, masterAddr)
		setContainerEnvIfNotFound(container, EnvPetMasterPort, strconv.Itoa(int(masterPort)))
		setContainerEnvIfNotFound(container, EnvPetNnodes, strconv.Itoa(nnodes))
		setContainerEnvIfNotFound(container, EnvPetNodeRank, strconv.Itoa(rank))
	}

	return nil
}

// setContainerEnv will add the specified environment variable to the container.
func setContainerEnv(container *corev1.Container, name string, value string) {
	if container == nil {
		return
	}

	container.Env = append(container.Env, corev1.EnvVar{Name: name, Value: value})
}

// setContainerEnvIfNotFound will add the specified environment variable to the container if not found.
func setContainerEnvIfNotFound(container *corev1.Container, name string, value string) {
	if container == nil {
		return
	}

	found := false
	for _, env := range container.Env {
		if env.Name == name {
			found = true
			break
		}
	}

	if !found {
		setContainerEnv(container, name, value)
	}
}
