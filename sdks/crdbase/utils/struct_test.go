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
	"testing"
)

func TestGetStructName(t *testing.T) {
	type tStruct struct{}

	tests := []struct {
		testStr any
		wantStr string
	}{
		{tStruct{}, "tStruct"},
		{&tStruct{}, "tStruct"},
		{"string type", ""},
		{1, ""},
	}

	for _, test := range tests {
		gotStr := GetStructName(test.testStr)
		if test.wantStr != gotStr {
			t.Errorf("%s differ (-got, +want): %s %s", test.testStr, gotStr, test.wantStr)
		}
	}
}

func TestGetStructID(t *testing.T) {
	type tStruct struct{}

	tests := []struct {
		testStr any
		wantStr string
	}{
		{tStruct{}, "utils.tStruct"},
		{&tStruct{}, "utils.tStruct"},
		{"string type", ""},
		{1, ""},
	}

	for _, test := range tests {
		gotStr := GetStructID(test.testStr)
		if test.wantStr != gotStr {
			t.Errorf("%s differ (-got, +want): %s %s", test.testStr, gotStr, test.wantStr)
		}
	}
}

func TestEnsureStructSlice(t *testing.T) {
	type test struct {
		k int
		v string
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				v: []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				v: &[]string{"1", "2", "3"},
			},
			want: true,
		},
		{
			name: "3",
			args: args{
				v: &test{
					k: 0,
					v: "1",
				},
			},
			want: false,
		},
		{
			name: "3",
			args: args{
				v: test{
					k: 0,
					v: "1",
				},
			},
			want: false,
		},
		{
			name: "4",
			args: args{
				v: []test{
					{k: 0, v: "0"},
					{k: 1, v: "1"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnsureStructSlice(tt.args.v); got != tt.want {
				t.Errorf("EnsureStructSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
