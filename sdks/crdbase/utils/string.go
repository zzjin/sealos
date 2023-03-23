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
	"math/rand"
	"strings"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/osamingo/base58"
	"github.com/sony/sonyflake"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	sfStart = 1648199220 // start time
)

var (
	st = sonyflake.Settings{
		StartTime: time.Unix(sfStart, 0),
	}
)

var (
	sf    *sonyflake.Sonyflake
	enc58 *base58.Encoder
)

func init() {
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		// if we cannot get private-ip, us another random id to init
		// @see https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go
		st.MachineID = func() (uint16, error) { return uint16(rand.Uint32() % (1<<16 - 1)), nil }
		sf = sonyflake.NewSonyflake(st)
		if sf == nil {
			panic("sonyflake not created")
		}
	}

	// use standard base58 to ensure shortest string id
	enc58, _ = base58.NewEncoder(base58.StandardSource)
}

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

// GenNanoID generates a short-term unique id, with timestamp order.
// When sonyflake returns error(witch almost never happen), return real random string.
// @Note: All random string must follow [RFC 1035 Label Names](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#:~:text=contain%20no%20more%20than%20253%20characters)
func GenNanoID() string {
	id, err := sf.NextID()
	if err != nil {
		// over the max speed, using another time method?
		// Use the go-nanoid package to generate a custom nanoid with length
		randID, _ := nanoid.Custom(alphabet, 12)
		return randID()

	}

	return enc58.Encode(id)
}
