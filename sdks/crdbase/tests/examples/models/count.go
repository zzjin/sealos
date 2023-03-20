package models

import (
	"context"
	"strconv"

	crdb "github.com/labring/crdbase"
	"github.com/labring/crdbase/query"
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
	Name      string    `json:"name" crdb:"name,primaryKey"`
	CountType CountType `json:"type" crdb:"type,index"`
	Counter   int64     `json:"count"`
}

func (c *Count) Get(db *crdb.CrdBase, name string, countType CountType) (int64, error) {
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

	if err := db.Model(c).Get(context.TODO(), q, &ret); err != nil {
		return 0, err
	}

	return 0, nil
}
