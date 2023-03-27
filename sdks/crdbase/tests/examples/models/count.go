package models

import (
	"strconv"

	crdb "github.com/labring/crdbase"
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
	Name      string    `json:"name" crdb:"name,unique"`
	CountType CountType `json:"type" crdb:"type,index"`
	Counter   int64     `json:"count"`
}

func (c *Count) Get(db *crdb.CRDBase, name string, countType CountType) (int64, error) {
	return 0, nil
}
