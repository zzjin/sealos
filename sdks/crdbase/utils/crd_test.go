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

package utils

import (
	"testing"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestNewCRDJSONSchemaProps(t *testing.T) {
	newObj := NewCRDJSONSchemaProps()
	if v, ok := newObj["spec"]; ok {
		v.Properties = map[string]apiextv1.JSONSchemaProps{
			"test": {},
		}
	}
	if v, ok := _baseSchemaProps["spec"]; ok {
		if v.Properties != nil {
			t.Error("base field changed with error")
		}
	}
}
