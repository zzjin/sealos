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
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/labring/crdbase/fork/github.com/fatih/structtag"
	"github.com/labring/crdbase/tests"
)

func TestParseFieldsTagsByStruct(t *testing.T) {

	testData := []struct {
		testData      any
		wantErr       error
		wantFieldData map[string]reflect.StructField
		wantTagData   map[string]*structtag.Tag
	}{
		{tests.TestModel{}, nil, tests.TestModelStructField, tests.TestModelStructTag},
		{&tests.TestModel{}, nil, tests.TestModelStructField, tests.TestModelStructTag},
		{"UserController", ErrNoFieldTag, nil, nil},
		{time.Now(), ErrNoFieldTag, nil, nil}, // test empty exported field
	}

	for _, test := range testData {
		gotFieldData, gotTagData, gotErr := ParseFieldsTagsByStruct(test.testData, "crdb")
		if !errors.Is(gotErr, test.wantErr) {
			t.Errorf("%s differ (-got, +want): %s %s", test.testData, gotErr, test.wantErr)
		} else {
			if ok := reflect.DeepEqual(test.wantFieldData, gotFieldData); !ok {
				t.Errorf("parsed field mismatch\n")
			}
			if diff := cmp.Diff(test.wantTagData, gotTagData); diff != "" {
				t.Errorf("parsed tag mismatch (-want +got):\n%s", diff)
			}
		}
	}
}
