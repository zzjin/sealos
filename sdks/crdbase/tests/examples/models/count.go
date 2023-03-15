package models

type CountType int

const (
	DownloadCount CountType = iota
	UpdateCount
)

type Count struct {
	Name      string    `json:"name" crdb:"primary"`
	CountType CountType `json:"type" crdb:"index"`
	Counter   int64     `json:"count"`
}

func (c *Count) Add(name string, countType CountType, step int) (int64, error) {
	// crdbtest.DB.Get()

	return 0, nil
}
