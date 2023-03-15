package models

import (
	"context"
	"strconv"

	"github.com/labring/crdbase/query"
	"github.com/labring/crdbase/tests/examples/crdbtest"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/selection"
)

type CountType int

const (
	DownloadCount CountType = iota
	UpdateCount
)

func (ct CountType) String() string {
	return strconv.Itoa(int(ct))
}

type Count struct {
	Name      string    `json:"name" crdb:"primaryKey"`
	CountType CountType `json:"type" crdb:"index"`
	Counter   int64     `json:"count"`
}

func (c *Count) Add(name string, countType CountType, step int) (int64, error) {
	data := &Count{
		Name:      name,
		CountType: countType,
		Counter:   1,
	}

	if _, _, err := crdbtest.DB.Model(Count{}).CreateOrUpdate(context.TODO(), data); err != nil {
		return 0, err
	}

	return data.Counter, nil
}

func (c *Count) Get(name string, countType CountType) (int64, error) {
	ret := &Count{}

	q := query.Query{
		FieldSelectors: fields.Requirements{
			{
				Field:    ".spec.name",
				Value:    name,
				Operator: selection.Equals,
			},
			{
				Field:    ".spec.type",
				Value:    countType.String(),
				Operator: selection.Equals,
			},
		},
	}

	if err := crdbtest.DB.Model(Count{}).Get(context.TODO(), q, &ret); err != nil {
		return 0, err
	}

	return 0, nil
}
