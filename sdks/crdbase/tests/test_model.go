package tests

import (
	"reflect"

	"github.com/labring/crdbase/fork/github.com/fatih/structtag"
)

// TestModel for all test s model define
type TestModel struct {
	User string    `json:"user" crdb:"user,primary"`
	Name string    `json:"name" crdb:"name,index"`
	Age  int64     `json:"age"`
	Info *TestInfo `json:"info"`

	self int64
}

// TestInfo for sub struct
type TestInfo struct {
	Gender int8 `json:"gender"`
}

var (
	TestModelStructField = map[string]reflect.StructField{
		"User": {
			Name:      "User",
			PkgPath:   "",
			Type:      reflect.TypeOf(string("")),
			Tag:       reflect.StructTag(`json:"user" crdb:"user,primary"`),
			Offset:    0,
			Index:     []int{0},
			Anonymous: false,
		},
		"Name": {
			Name:      "Name",
			PkgPath:   "",
			Type:      reflect.TypeOf(string("")),
			Tag:       reflect.StructTag(`json:"name" crdb:"name,index"`),
			Offset:    16,
			Index:     []int{1},
			Anonymous: false,
		},
		"Age": {
			Name:      "Age",
			PkgPath:   "",
			Type:      reflect.TypeOf(int64(0)),
			Tag:       reflect.StructTag(`json:"age"`),
			Offset:    32,
			Index:     []int{2},
			Anonymous: false,
		},
		"Info": {
			Name:      "Info",
			PkgPath:   "",
			Type:      reflect.TypeOf(&TestInfo{}),
			Tag:       reflect.StructTag(`json:"info"`),
			Offset:    40,
			Index:     []int{3},
			Anonymous: false,
		},
	}

	TestModelStructTag = map[string]*structtag.Tag{
		"User": {
			Key:     "crdb",
			Name:    "user",
			Options: []string{"primary"},
		},
		"Name": {
			Key:     "crdb",
			Name:    "name",
			Options: []string{"index"},
		},
		"Age":  nil,
		"Info": nil,
	}
)
