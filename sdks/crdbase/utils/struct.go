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

	"github.com/invopop/jsonschema"
	"github.com/mitchellh/mapstructure"
)

const StructTagName = "json"

// EnsureStructs Ensure that the type of the variable is a struct
func EnsureStructs(is ...any) ([]any, error) {
	ret := []any{}
	for _, i := range is {
		if i, err := EnsureStruct(i); err != nil {
			return nil, err
		} else {
			ret = append(ret, i)
		}
	}
	return ret, nil
}

// EnsureStruct Ensure that the type of the variable is a struct
func EnsureStruct(i any) (any, error) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() != reflect.Struct {
			return nil, errors.New("not a pointer to a struct")
		}
		return reflect.New(t.Elem()).Interface(), nil
	} else if t.Kind() == reflect.Struct {
		return i, nil
	}
	return nil, errors.New("not a struct or a pointer to a struct")
}

// GetStructName returns the name of the struct
func GetStructName(i any) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() == reflect.Struct {
			return t.Elem().Name()
		}
	} else if t.Kind() == reflect.Struct {
		return t.Name()
	}

	return ""
}

// GetStructID returns the name of the struct
func GetStructID(i any) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() == reflect.Struct {
			return t.Elem().String()
		}
	} else if t.Kind() == reflect.Struct {
		return t.String()
	}

	return ""
}

// GetStructExportedFields returns all exported fields of the struct
func GetStructExportedFields(i any) map[string]reflect.StructField {
	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() == reflect.Struct {
			t = t.Elem()
		} else {
			return nil
		}
	} else if t.Kind() != reflect.Struct {
		return nil
	}

	ret := map[string]reflect.StructField{}

	fLen := t.NumField()
	for i := 0; i < fLen; i++ {
		field := t.Field(i)
		if field.IsExported() {
			ret[field.Name] = field
		}
	}

	return ret
}

// StructJSON2Map converts a struct to a map
func StructJSON2Map(in any) (map[string]any, error) {
	out := map[string]any{}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &out,
		TagName: StructTagName,
	})
	if err != nil {
		return nil, err
	}
	if err := dec.Decode(in); err != nil {
		return nil, err
	}

	return out, nil
}

// Map2JSONStruct converts a map to a struct
func Map2JSONStruct(in map[string]any, out any) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &out,
		TagName: StructTagName,
	})
	if err != nil {
		return err
	}

	return dec.Decode(in)
}

// Struct2JSONSchema currently implemented by jsonschema.Reflect
func Struct2JSONSchema(in any) *jsonschema.Schema {
	reflect := &jsonschema.Reflector{DoNotReference: true, ExpandedStruct: true}
	return reflect.Reflect(in)
}
