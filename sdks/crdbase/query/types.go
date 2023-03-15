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

package query

// JP for short of `jsonpath`
const (
	JPCreationTimestamp    = "metadata.creationTimestamp"
	JPLastUpdatedTimestamp = "spec.lastUpdatedTimestamp"
	JPName                 = "metadata.name"
)

// IsSortable checks if the field is sortable
func IsSortable(s string) bool {
	return s == JPCreationTimestamp || s == JPLastUpdatedTimestamp || s == JPName
}

// DS for short of default selectable
const (
	DSName = "metadata.name"
	DSKind = "kind"
	DSStatus = "status"
)

// IsDefaultSelectable checks if the field is selectable
func IsDefaultSelectable(s string) bool {
	return s == DSName || s == DSKind || s == DSStatus
}
