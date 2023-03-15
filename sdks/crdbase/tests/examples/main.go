package main

import (
	"context"

	"github.com/labring/crdbase/tests/examples/crdbtest"
	"github.com/labring/crdbase/tests/examples/models"
)

func main() {
	crdbtest.InitCRDB()

	ctx := context.TODO()

	crdbtest.DB.AutoMigrate(ctx, models.Count{})

}
