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
	"errors"

	"github.com/labring/crdbase/utils"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Model any
type Data any

type ModelKindNameInterface interface {
	KindName() string
}

// Model2CRD Unstructured parse model and convert it to crd struct to apply
func (crdb *CrdBase) Model2CRD(m Model) (*apiextv1.CustomResourceDefinition, error) {
	schema := GetCrdModelSchema(m)

	if schema.IsEmpty() {
		return nil, errors.New("model cannot converted to crd")
	}

	crdJSONSchema := utils.NewCRDJSONSchemaProps()
	// since we control the new props func, so just do it
	if entry, ok := crdJSONSchema["spec"]; ok {
		entry.Properties = schema.Spec
		crdJSONSchema["spec"] = entry
	}

	return crdb.NewCustomResourceDefinition(schema.Names, crdJSONSchema), nil
}

func Model2KindName(m Model) apiextv1.CustomResourceDefinitionNames {
	var name string
	if mn, ok := m.(ModelKindNameInterface); ok {
		name = mn.KindName()
	} else {
		name = utils.GetStructName(m)
	}
	return utils.Name2KindNames(name)
}
