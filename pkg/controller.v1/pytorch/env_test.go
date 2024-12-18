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
	"testing"

	v1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	"github.com/kubeflow/pytorch-operator/pkg/common/util/v1/testutil"
)

func TestSetPodTemplateEnv(t *testing.T) {
	testCases := []struct {
		job          *v1.PyTorchJob
		rt           v1.PyTorchReplicaType
		index        string
		expectedEnvs map[string]string
	}{
		{
			job:   testutil.NewPyTorchJobWithMaster(0),
			rt:    v1.PyTorchReplicaTypeMaster,
			index: "0",
			expectedEnvs: map[string]string{
				EnvMasterAddr:       "test-pytorchjob-master-0",
				EnvMasterPort:       "23456",
				EnvWorldSize:        "1",
				EnvRank:             "0",
				EnvPythonUnbuffered: "1",
				EnvPetMasterAddr:    "test-pytorchjob-master-0",
				EnvPetMasterPort:    "23456",
				EnvPetNnodes:        "1",
			},
		},
		{
			job:   testutil.NewPyTorchJobWithMaster(1),
			rt:    v1.PyTorchReplicaTypeMaster,
			index: "0",
			expectedEnvs: map[string]string{
				EnvMasterAddr:       "test-pytorchjob-master-0",
				EnvMasterPort:       "23456",
				EnvWorldSize:        "2",
				EnvRank:             "0",
				EnvPythonUnbuffered: "1",
				EnvPetMasterAddr:    "test-pytorchjob-master-0",
				EnvPetMasterPort:    "23456",
				EnvPetNnodes:        "2",
			},
		},
		{
			job:   testutil.NewPyTorchJobWithMaster(1),
			rt:    v1.PyTorchReplicaTypeWorker,
			index: "0",
			expectedEnvs: map[string]string{
				EnvMasterAddr:       "test-pytorchjob-master-0",
				EnvMasterPort:       "23456",
				EnvWorldSize:        "2",
				EnvRank:             "1",
				EnvPythonUnbuffered: "1",
				EnvPetMasterAddr:    "test-pytorchjob-master-0",
				EnvPetMasterPort:    "23456",
				EnvPetNnodes:        "2",
			},
		},
		{
			job:   testutil.NewPyTorchJobWithMaster(2),
			rt:    v1.PyTorchReplicaTypeMaster,
			index: "0",
			expectedEnvs: map[string]string{
				EnvMasterAddr:       "test-pytorchjob-master-0",
				EnvMasterPort:       "23456",
				EnvWorldSize:        "3",
				EnvRank:             "0",
				EnvPythonUnbuffered: "1",
				EnvPetMasterAddr:    "test-pytorchjob-master-0",
				EnvPetMasterPort:    "23456",
				EnvPetNnodes:        "3",
			},
		},
		{
			job:   testutil.NewPyTorchJobWithMaster(2),
			rt:    v1.PyTorchReplicaTypeWorker,
			index: "0",
			expectedEnvs: map[string]string{
				EnvMasterAddr:       "test-pytorchjob-master-0",
				EnvMasterPort:       "23456",
				EnvWorldSize:        "3",
				EnvRank:             "1",
				EnvPythonUnbuffered: "1",
				EnvPetMasterAddr:    "test-pytorchjob-master-0",
				EnvPetMasterPort:    "23456",
				EnvPetNnodes:        "3",
			},
		},
		{
			job:   testutil.NewPyTorchJobWithMaster(2),
			rt:    v1.PyTorchReplicaTypeWorker,
			index: "1",
			expectedEnvs: map[string]string{
				EnvMasterAddr:       "test-pytorchjob-master-0",
				EnvMasterPort:       "23456",
				EnvWorldSize:        "3",
				EnvRank:             "2",
				EnvPythonUnbuffered: "1",
				EnvPetMasterAddr:    "test-pytorchjob-master-0",
				EnvPetMasterPort:    "23456",
				EnvPetNnodes:        "3",
			},
		},
	}

	for _, tc := range testCases {
		templateSpec := tc.job.Spec.PyTorchReplicaSpecs[tc.rt].Template
		if err := setPodTemplateEnv(&templateSpec, tc.job, tc.rt, tc.index); err != nil {
			t.Errorf("Failed to set pod template env: %v", err)
		}

		actualEnvs := map[string]string{}
		for _, env := range templateSpec.Spec.Containers[0].Env {
			actualEnvs[env.Name] = env.Value
		}

		for key, expectedVal := range tc.expectedEnvs {
			actualVal, ok := actualEnvs[key]
			if !ok {
				t.Errorf("Environment variable \"%s\" not found in pod template", key)
			}
			if actualVal != expectedVal {
				t.Errorf("Environment variable \"%s\" not equal to \"%s\", actual: \"%s\"", key, expectedVal, actualVal)
			}
		}
	}
}
