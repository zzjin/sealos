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
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/labring/crdbase/tests"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GuessShortName guesses the short name of the given name.
func TestModel2CRD(t *testing.T) {
	crdbObj := NewTestCRDBase()

	testData := []struct {
		testData Model
		wantErr  error
		wantData *apiextv1.CustomResourceDefinition
	}{
		{tests.TestModel{}, nil, &apiextv1.CustomResourceDefinition{}},
		{&tests.TestModel{}, nil, &apiextv1.CustomResourceDefinition{}},
		{"UserController", nil, &apiextv1.CustomResourceDefinition{}},
		{time.Now(), errors.New("asdf"), &apiextv1.CustomResourceDefinition{}}, // test empty exported field
	}

	for _, test := range testData {
		gotData, gotErr := crdbObj.Model2CRD(test.testData)
		if gotErr != test.wantErr {
			t.Errorf("%s differ (-got, +want): %s %s", test.testData, gotErr, test.wantErr)
		} else {
			// y, _ := yaml.Marshal(gotData)
			// fmt.Println(string(y))

			if diff := cmp.Diff(test.wantData, gotData); diff != "" {
				t.Errorf("model-tp-crd mismatch (-want +got):\n%s", diff)
			}
		}
	}
}

func TestCreate(t *testing.T) {
	crdbObj := NewTestCRDBase()
	ctx := context.TODO()

	tests := []struct {
		testData Model
		wantErr  error
		wantID   string
		wantOR   controllerutil.OperationResult
	}{
		{
			tests.TestModel{
				User: "testPrimary",
				Name: "testUser",
				Age:  1,
				Info: &tests.TestInfo{
					Gender: 1,
				},
			},
			nil,
			"testPrimary",
			controllerutil.OperationResultCreated,
		},
	}

	for _, test := range tests {
		gotID, gotOR, gotErr := crdbObj.Model(test.testData).Create(ctx, test.testData)
		if gotErr != test.wantErr {
			t.Errorf("%s differ (-got, +want): %s %s", test.testData, gotErr, test.wantErr)
		} else {
			// y, _ := yaml.Marshal(gotData)
			// fmt.Println(string(y))

			if diff := cmp.Diff(test.wantID, gotID); diff != "" {
				t.Errorf("create model id mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantOR, gotOR); diff != "" {
				t.Errorf("create model id mismatch (-want +got):\n%s", diff)
			}
		}
	}
}
