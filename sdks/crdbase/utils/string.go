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
	"strings"

	"github.com/jaevor/go-nanoid"
)

// GuessShortNames guesses the short name of the given name.
func GuessShortNames(name string) string {
	var ret strings.Builder

	for _, r := range name {
		if r >= 'A' && r <= 'Z' {
			ret.WriteRune(rune(r + 32))
		}
	}

	return ret.String()
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// GenerateMetaName generates a lowercase RFC 1123-compliant nanoid with specified length
func GenerateMetaName(length ...int) string {
	nanoLength := 12
	if len(length) > 0 {
		if length[0] >= 2 || length[0] <= 255 {
			nanoLength = length[0]
		}
	}
	// Use the go-nanoid package to generate a custom nanoid with length
	gener, _ := nanoid.Custom(alphabet, nanoLength)
	label := gener()
	return label
}
