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

	"github.com/labring/crdbase/fork/github.com/fatih/structtag"
)

var (
	ErrNoFieldTag = errors.New("not found exported fields")
)

// ParseFieldsTagsByStruct parse struct fields and (passed-in key) tags
func ParseFieldsTagsByStruct(a any, tagKey string) (map[string]reflect.StructField, map[string]*structtag.Tag, error) {
	fieldsMap := GetStructExportedFields(a)

	if len(fieldsMap) == 0 {
		return nil, nil, ErrNoFieldTag
	}

	keyTags := map[string]*structtag.Tag{}

	for name, field := range fieldsMap {
		allTags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return nil, nil, err
		}

		tag, err := allTags.Get(tagKey)
		if err != nil && !errors.Is(err, structtag.ErrTagNotExist) {
			return nil, nil, err
		}

		keyTags[name] = tag
	}

	return fieldsMap, keyTags, nil
}
