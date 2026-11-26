module github.com/kubeflow/pytorch-operator

go 1.26.0

require (
	github.com/go-openapi/spec v0.19.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.4
	github.com/kubeflow/common v0.1.0
	github.com/kubeflow/tf-operator v1.0.0-rc.0.0.20190916040037-5adee6f30c86
	github.com/kubernetes-sigs/kube-batch v0.5.0
	github.com/kubernetes-sigs/yaml v1.1.0
	github.com/onrik/logrus v0.4.0
	github.com/prometheus/client_golang v1.25.0
	github.com/sirupsen/logrus v1.10.2
	k8s.io/api v0.15.9
	k8s.io/apimachinery v0.15.9
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v0.4.0
	k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/kubernetes v1.15.9
)

require (
	cloud.google.com/go v0.44.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/emicklei/go-restful v2.9.6+incompatible // indirect
	github.com/evanphx/json-patch v0.0.0-20190203023257-5858425f7550 // indirect
	github.com/go-openapi/jsonpointer v0.19.2 // indirect
	github.com/go-openapi/jsonreference v0.19.2 // indirect
	github.com/go-openapi/swag v0.19.4 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.72.0 // indirect
	github.com/prometheus/procfs v0.22.0 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/oauth2 v0.37.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
	k8s.io/apiextensions-apiserver v0.0.0 // indirect
	k8s.io/apiserver v0.15.9 // indirect
	k8s.io/component-base v0.15.9 // indirect
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.15.9
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.15.9
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.10-beta.0
	k8s.io/apiserver => k8s.io/apiserver v0.15.9
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.15.9
	k8s.io/client-go => k8s.io/client-go v0.15.9
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.15.9
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.15.9
	k8s.io/code-generator => k8s.io/code-generator v0.15.13-beta.0
	k8s.io/component-base => k8s.io/component-base v0.15.9
	k8s.io/cri-api => k8s.io/cri-api v0.15.13-beta.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.15.9
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.15.9
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.15.9
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.15.9
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.15.9
	k8s.io/kubectl => k8s.io/kubectl v0.15.13-beta.0
	k8s.io/kubelet => k8s.io/kubelet v0.15.9
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.15.9
	k8s.io/metrics => k8s.io/metrics v0.15.9
	k8s.io/node-api => k8s.io/node-api v0.15.9
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.15.9
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.15.9
	k8s.io/sample-controller => k8s.io/sample-controller v0.15.9
)
