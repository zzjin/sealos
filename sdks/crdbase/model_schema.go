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
	"strings"

	"github.com/labring/crdbase/utils"
	"golang.org/x/sync/singleflight"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ModelSchema struct {
	IdentifyKey string // key for unique model structure

	Indexes [][]string
	Uniques [][]string

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
	identifyKey := utils.GetStructID(m)
	if identifyKey == "" {
		return nil, fmt.Errorf("model %s id is empty", names.Kind)
	}

	crdJSONSchema := utils.Struct2JSONSchemaProps(m)

	fields, tags, err := utils.ParseFieldsTagsByStruct(m, crdBaseTagKey)
	if err != nil {
		return nil, err
	}

	var indexes [][]string
	var uniques [][]string

	for name, tag := range tags {
		if tag == nil {
			continue
		}
		for _, opt := range tag.Options {
			field, ok := fields[name]
			if !ok {
				return nil, fmt.Errorf("field %s not found", field.Name)
			}
			if opt == "index" {
				indexes = append(indexes, []string{field.Name})
				continue
			}
			if opt == "unique" {
				// unique ~= index
				indexes = append(indexes, []string{field.Name})
				uniques = append(uniques, []string{field.Name})
				break
			}
		}
	}

	indexes = utils.UniqSliceSlice(indexes)
	uniques = utils.UniqSliceSlice(uniques)

	return &ModelSchema{
		IdentifyKey: identifyKey,
		Names:       names,
		Spec:        crdJSONSchema,
		Indexes:     indexes,
		Uniques:     uniques,
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

func (ms *ModelSchema) GetIndexesFunc(indexes []string) (string, client.IndexerFunc) {
	refName := fmt.Sprintf("%sRef", strings.Join(indexes, ""))
	indexerFunc := func(obj client.Object) []string {
		return indexes
	}

	return refName, indexerFunc
}
