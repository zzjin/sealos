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

package crdb

import (
	"context"

	"github.com/labring/crdbase/query"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BaseConstructor interface {
	AutoMigrate(ctx context.Context, models ...any) error
	Prune(ctx context.Context, models ...any) error
}

type BaseWriter interface {
	Create(ctx context.Context, model any) (string, controllerutil.OperationResult, error)
	CreateOrUpdate(ctx context.Context, model any, uf MutateFn) (string, controllerutil.OperationResult, error)
	CreateOrUpdateList(ctx context.Context, model any, uf MutateFn) ([]string, controllerutil.OperationResult, error)

	Delete(ctx context.Context, name string) error
	DeleteAllOf(ctx context.Context, query query.Query) error
}

type BaseReader interface {
	Get(ctx context.Context, query query.Query, out any) error
	List(ctx context.Context, query query.Query, out any) error
}

type Base interface {
	BaseConstructor
	BaseWriter
	BaseReader
}

type BaseModel interface {
	BaseWriter
	BaseReader
}
