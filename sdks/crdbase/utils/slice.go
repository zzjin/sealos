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
	"fmt"

	"golang.org/x/exp/constraints"
)

// UniqSlice UniqSlice
func UniqSlice[T constraints.Ordered](list []T) []T {
	target := make([]T, 0)
	listMap := make(map[T]bool)
	for _, v := range list {
		if _, ok := listMap[v]; !ok {
			listMap[v] = true
			target = append(target, v)
		}
	}
	return target
}

// UniqSliceSlice parse lists and for each list, call the unique function to get the unique value.
func UniqSliceSlice[T constraints.Ordered](list [][]T, uniFunc ...func([]T) ([]T, string)) [][]T {
	uni := DefaultUniqueSliceString[T]
	if len(uniFunc) > 0 {
		uni = uniFunc[0]
	}

	target := make([][]T, 0)
	listMap := make(map[string]bool)

	for _, vs := range list {
		nvs, v := uni(vs)
		if _, ok := listMap[v]; !ok {
			listMap[v] = true
			target = append(target, nvs)
		}
	}
	return target
}

// DefaultUniqueSliceString default unique function for slice, and print slice to unique name string.
func DefaultUniqueSliceString[T constraints.Ordered](list []T) ([]T, string) {
	list = UniqSlice(list)
	return list, fmt.Sprint(list)
}
