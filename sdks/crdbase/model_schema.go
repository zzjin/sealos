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
	"fmt"
	"reflect"
	"sync"

	"github.com/labring/crdbase/utils"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type CRDModelSchema struct {
	ID           string
	Names        apiextv1.CustomResourceDefinitionNames
	Spec         map[string]apiextv1.JSONSchemaProps
	PrimaryField string
	Indexes      [][]string
}

var (
	modelSchemaMap = map[string]*CRDModelSchema{}
	mLock          sync.Mutex
)

func GetCRDModelsSchemas(models ...Model) []*CRDModelSchema {
	ret := []*CRDModelSchema{}
	for _, m := range models {
		schema := GetCRDModelSchema(m)
		if !schema.IsEmpty() {
			ret = append(ret, schema)
		}
	}
	return ret
}

func GetCRDModelSchema(m Model) *CRDModelSchema {
	emptyObj := &CRDModelSchema{
		ID:           "",
		Names:        apiextv1.CustomResourceDefinitionNames{},
		Spec:         map[string]apiextv1.JSONSchemaProps{},
		PrimaryField: "",
		Indexes:      [][]string{},
	}

	id := utils.GetStructID(m)
	if id == "" {
		return emptyObj
	}

	mLock.Lock()
	defer mLock.Unlock()

	if ret, ok := modelSchemaMap[id]; ok {
		return ret
	}

	obj, err := newCRDModelSchema(id, m)
	if err != nil {
		return emptyObj
	}

	modelSchemaMap[id] = obj

	return obj
}

func newCRDModelSchema(id string, m Model) (*CRDModelSchema, error) {
	names := Model2KindName(m)
	crdJSONSchema := utils.Struct2JSONSchemaProps(m)

	fields, tags, err := utils.ParseFieldsTagsByStruct(m, crdbaseTagKey)
	if err != nil {
		return nil, err
	}

	primaryField := ""
	indexes := [][]string{}
	for name, tag := range tags {
		if tag == nil {
			continue
		}

		for _, opt := range tag.Options {
			field, ok := fields[name]
			if !ok {
				return nil, fmt.Errorf("field %s not found", field.Name)
			}

			if opt == "primaryKey" {
				if primaryField != "" {
					return nil, fmt.Errorf("duplicate primary field %s and %s", primaryField, field.Name)
				}

				if field.Type.Kind() != reflect.String {
					return nil, fmt.Errorf("primary field %s must be string", field.Name)
				}

				primaryField = field.Name

				// primary ~= index
				indexes = append(indexes, []string{field.Name})

				break
			}

			if opt == "index" {
				indexes = append(indexes, []string{field.Name})
				continue
			}
		}
	}

	return &CRDModelSchema{
		ID:           id,
		Names:        names,
		Spec:         crdJSONSchema,
		PrimaryField: primaryField,
		Indexes:      indexes,
	}, nil
}

func (crdms *CRDModelSchema) IsEmpty() bool {
	return (crdms.Names.Plural == "" || crdms.Names.Singular == "") || len(crdms.Spec) == 0
}

func (crdms *CRDModelSchema) ResourceName() string {
	return crdms.Names.Plural
}

func (crdms *CRDModelSchema) GetPrimaryFieldValue(m Model) string {
	if crdms.PrimaryField != "" {
		t := reflect.TypeOf(m)

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t.Kind() == reflect.Struct {
			str := reflect.ValueOf(m).FieldByName(crdms.PrimaryField).String()
			if str != "" {
				return str
			}
		}
	}

	return ""
}
