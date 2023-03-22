package models

import (
	"context"
	"strconv"

	crdb "github.com/labring/crdbase"
	"github.com/labring/crdbase/query"
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
	q := query.Query{}
	res, _ := db.Model(c).Get(context.Background(), q)
	return 0, nil
}
