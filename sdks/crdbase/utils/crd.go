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
