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
	"github.com/iancoleman/strcase"
	"github.com/invopop/jsonschema"
	"github.com/jinzhu/inflection"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	_baseSchemaProps = map[string]apiextv1.JSONSchemaProps{
		"apiVersion": {
			Description: `APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources`,
			Type: "string",
		},
		"kind": {
			Description: `Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated. In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds`,
			Type: "string",
		},
		"metadata": {
			Type: "object",
		},
		"spec": {
			Description: "crd spec defines the real saved data",
			Type:        "object",
		},
		"status": {
			Description: "status defines the observed state of saved data",
			Type:        "object",
		},
	}
)

// NewCRDJSONSchemaProps use deep-copy to new object, cannot change _baseSchemaProps
func NewCRDJSONSchemaProps() map[string]apiextv1.JSONSchemaProps {
	return map[string]apiextv1.JSONSchemaProps{
		"apiVersion": {Description: _baseSchemaProps["apiVersion"].Description, Type: _baseSchemaProps["apiVersion"].Type},
		"kind":       {Description: _baseSchemaProps["kind"].Description, Type: _baseSchemaProps["kind"].Type},
		"metadata":   {Type: _baseSchemaProps["metadata"].Type},
		"spec":       {Description: _baseSchemaProps["spec"].Description, Type: _baseSchemaProps["spec"].Type},
		"status":     {Description: _baseSchemaProps["status"].Description, Type: _baseSchemaProps["status"].Type},
	}
}

func Name2KindNames(name string) apiextv1.CustomResourceDefinitionNames {
	kebab := strcase.ToKebab(name)

	plural := inflection.Plural(kebab)
	singular := inflection.Singular(kebab)

	ret := apiextv1.CustomResourceDefinitionNames{
		Plural:   plural,
		Singular: singular,
		Kind:     name,
		ListKind: name + "List",
	}

	shortNames := GuessShortNames(name)
	if shortNames != "" {
		ret.ShortNames = []string{shortNames}
	}

	return ret
}

// Struct2JSONSchemaProps returns struct to apiextensions.k8s.io/v1.JSONSchemaProps
func Struct2JSONSchemaProps(in any) map[string]apiextv1.JSONSchemaProps {
	inJSONSchema := Struct2JSONSchema(in)
	ret := buildJSONSchemaProps(inJSONSchema)
	return ret
}

func buildJSONSchemaProps(from *jsonschema.Schema) map[string]apiextv1.JSONSchemaProps {
	ret := map[string]apiextv1.JSONSchemaProps{}

	fieldKeys := from.Properties.Keys()
	for _, key := range fieldKeys {
		if schemaI, ok := from.Properties.Get(key); ok {
			if schema, ok2 := schemaI.(*jsonschema.Schema); ok2 {
				oneSchema := apiextv1.JSONSchemaProps{
					Description: schema.Description,
					Type:        schema.Type,
					Format:      schema.Format,
				}

				if schema.Properties != nil {
					oneSchema.Properties = buildJSONSchemaProps(schema)
				}

				ret[key] = oneSchema
			}
		}
	}

	return ret
}
