// Copyright © 2023 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crdb

import (
	"flag"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// GuessShortName guesses the short name of the given name.
func NewTestCRDBase() *CrdBase {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)

	testLogger := zap.New(zap.UseFlagOptions(&opts))

	// TODO: add fake or test client config
	var mgr manager.Manager

	conf := CrdBaseConfig{
		Manager: mgr,
		GroupVersion: schema.GroupVersion{
			Group:   "test.crdb.sealos.io",
			Version: "v1",
		},
		ServiceAccount: "sealos-test-manager",
		Namespace:      "sealos-test",
	}

	crdbase, _ := NewCrdBase(conf, testLogger)

	return crdbase
}
