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
	"testing"

	"github.com/jaevor/go-nanoid"
)

// GuessShortName guesses the short name of the given name.
func TestGuessShortName(t *testing.T) {
	tests := []struct {
		testStr string
		wantStr string
	}{
		{"UserName", "un"},
		{"User-controller", "u"},
		{"access", ""},
		{"Not Found", "nf"},
		{"not found", ""},
	}

	for _, test := range tests {
		gotStr := GuessShortNames(test.testStr)
		if test.wantStr != gotStr {
			t.Errorf("%s differ (-got, +want): %s %s", test.testStr, gotStr, test.wantStr)
		}
	}
}

func TestGenerateID(t *testing.T) {
	m := map[string]struct{}{}

	for i := 0; i < 1000; i++ {
		id := GenerateID()
		if _, ok := m[id]; ok || id == "" {
			t.Errorf("nano id duplicated: %s", id)
		} else {
			m[id] = struct{}{}
		}
	}
}

func TestGenerateIDRFC1123(t *testing.T) {
	for i := 0; i < 10000000; i++ {
		id := GenerateID()
		if id != "" && id[0] >= 30 && id[0] <= 39 {
			fmt.Printf("not spec: %s\n", id)
		}
	}
}

func BenchmarkSonyFlakeID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idInt, _ := sf.NextID()

		id := enc34.Encode(idInt)
		if id == "" {
			b.FailNow()
		}
	}
}
func BenchmarkNanoID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idFunc, _ := nanoid.Standard(12)
		id := idFunc()
		if id == "" {
			b.FailNow()
		}
	}
}
