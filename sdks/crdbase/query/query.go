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

import (
	"errors"
	"net/url"

	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/conversion/queryparams"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type operator int

const (
	lessThan operator = iota + 1
	lessEq
	equal
	greaterEq
	greaterThan
	in
	notIn
	notEqual

	keyFieldName = "__key__"
)

// filter is a conditional filter on query results.
type filter struct {
	FieldName string
	Op        operator
	Value     interface{}
}

type sortDirection bool

const (
	ascending  sortDirection = false
	descending sortDirection = true
)

// order is a sort order on query results.
type order struct {
	FieldName string
	Direction sortDirection
}

// Query is the query object for transforming a url query to api-server list call.
type Query struct {
	filter []filter
	order  []order

	// Pagination
	// Page  stands for the page number, default to 1 (start from 1)
	Page int `json:"page,omitempty"`
	// Limit stands for the number of items per page, default to 10 (maximum to 1000)
	Limit int `json:"limit,omitempty"`

	// Sort
	// SortBy sort result in which field, default to FieldCreationTimeStamp
	SortBy string `json:"sort_by,omitempty"`
	// Ascending sort result in ascending or descending order, default to descending
	Ascending int8 `json:"ascending,omitempty"`

	// Filters filters the result by key: jsonpath, value: value
	FieldSelectors fields.Requirements `json:"field_selectors,omitempty"`

	// LabelSelector filters the result by key: label, value: value
	LabelSelectors labels.Requirements `json:"label_selectors,omitempty"`
}

func New(page, limit int, sortBy string, ascending int8) *Query {
	return &Query{
		Page:           1,
		Limit:          10,
		SortBy:         "",
		Ascending:      0,
		FieldSelectors: fields.Requirements{},
		LabelSelectors: labels.Requirements{},
	}
}

func Parse(query string) (*Query, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	ret := &Query{}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  ret,
	})
	if err != nil {
		return nil, err
	}
	if err := dec.Decode(values); err != nil {
		return nil, err
	}

	if err := ret.Validate(); err != nil {
		return nil, err
	}

	return ret, nil
}

// Validate ensure all query parameters are valid
func (q *Query) Validate() error {
	if q.Page < 0 {
		return errors.New("pagination must be greater than 0")
	}
	if q.Limit < 0 || q.Limit > 1000 {
		return errors.New("limit must be greater than 0")
	}

	return nil
}

func (q *Query) String() string {
	// here we use our defined struct, hence never fail
	// nosemgrep
	obj, _ := queryparams.Convert(q)
	return obj.Encode()
}

func (q *Query) ToListOptions() []client.ListOption {
	// TODO: Impl.
	sel := labels.NewSelector()
	sel.Add(q.LabelSelectors...)

	return []client.ListOption{client.MatchingFields{}}
}

func (q *Query) PostFilter(data any) any {
	// TODO: Impl.
	return data
}
