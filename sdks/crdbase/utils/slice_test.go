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
	"reflect"
	"strings"
	"testing"
)

func TestUniqSlice(t *testing.T) {
	got := UniqSlice([]string{"a", "b", "c", "c"})
	if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("UniqSlice() = %v", got)
	}

	got2 := UniqSlice([]int64{1, 2, 3, 4, 4})
	if !reflect.DeepEqual(got2, []int64{1, 2, 3, 4}) {
		t.Errorf("UniqSlice() = %v", got2)
	}
}

func TestUniqSliceSlice(t *testing.T) {
	got := UniqSliceSlice(
		[][]string{
			{"a", "b", "c", "c"},
			{"a", "b", "c"},
		},
		func(ss []string) ([]string, string) {
			ss = UniqSlice(ss)
			return ss, strings.Join(ss, ",")
		},
	)
	if !reflect.DeepEqual(got, [][]string{{"a", "b", "c"}}) {
		t.Errorf("UniqSliceSlice() = %v", got)
	}

	got2 := UniqSliceSlice(
		[][]int64{
			{1, 2, 3, 4, 4},
			{1, 2, 3, 4},
		},
		func(ss []int64) ([]int64, string) {
			ss = UniqSlice(ss)
			return ss, fmt.Sprint(ss)
		},
	)
	if !reflect.DeepEqual(got2, [][]int64{{1, 2, 3, 4}}) {
		t.Errorf("UniqSliceSlice() = %v", got2)
	}
}
