package crdbtest

import (
	"flag"

	crdb "github.com/labring/crdbase"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	DB *crdb.CRDBase
)

func InitCRDB() {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)

	log := zap.New(zap.UseFlagOptions(&opts))

	conf := crdb.CRDBaseConfig{
		Manager: nil,
		GroupVersion: schema.GroupVersion{
			Group:   "crdb.sealos.io",
			Version: "v1",
		},
		ServiceAccount: "crdb-test",
		Namespace:      "crdb-test",
	}

	var err error

	DB, err = crdb.NewCRDBase(conf, log)
	if err != nil {
		panic(err)
	}

}
