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

	"github.com/labring/crdbase/utils"

	"golang.org/x/sync/singleflight"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type ModelSchema struct {
	ID           string
	PrimaryField string

	Indexes [][]string

	Names apiextv1.CustomResourceDefinitionNames
	Spec  map[string]apiextv1.JSONSchemaProps
}

var (
	modelSchemaGroup = singleflight.Group{}
)

func GetCRDModelsSchemas(models ...Model) []*ModelSchema {
	var ret []*ModelSchema
	for _, m := range models {
		schema := GetCRDModelSchema(m)
		if !schema.IsEmpty() {
			ret = append(ret, schema)
		}
	}
	return ret
}

func GetCRDModelSchema(m Model) *ModelSchema {
	// use struct id as key
	res, _, _ := modelSchemaGroup.Do(utils.GetStructID(m), func() (interface{}, error) {
		obj, err := newCRDModelSchema(m)
		if err != nil {
			return &ModelSchema{}, err
		}
		return obj, nil
	})
	return res.(*ModelSchema)
}

// newCRDModelSchema
func newCRDModelSchema(m Model) (*ModelSchema, error) {
	names := Model2KindName(m)

	// use struct id as key, model should be struct or a struct pointer
	id := utils.GetStructID(m)
	if id == "" {
		return nil, fmt.Errorf("model %s id is empty", names.Kind)
	}

	crdJSONSchema := utils.Struct2JSONSchemaProps(m)

	fields, tags, err := utils.ParseFieldsTagsByStruct(m, crdBaseTagKey)
	if err != nil {
		return nil, err
	}

	var primaryField string
	var indexes [][]string

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

	return &ModelSchema{
		ID:           id,
		Names:        names,
		Spec:         crdJSONSchema,
		PrimaryField: primaryField,
		Indexes:      indexes,
	}, nil
}

func (ms *ModelSchema) IsEmpty() bool {
	return (ms.Names.Plural == "" || ms.Names.Singular == "") || len(ms.Spec) == 0
}

func (ms *ModelSchema) ResourceName() string {
	return ms.Names.Plural
}

func (ms *ModelSchema) Kind() string {
	return ms.Names.Kind
}

// GetPrimaryFieldValue get primary field value from a model
func (ms *ModelSchema) GetPrimaryFieldValue(m Data) string {
	if ms.PrimaryField != "" {
		v := reflect.ValueOf(m)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		v = v.FieldByName(ms.PrimaryField)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.String {
			str := v.String()
			if str != "" {
				return str
			}
		}
	}

	return ""
}
