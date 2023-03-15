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

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (crdb *CRDBase) Create(ctx context.Context, model any) (string, controllerutil.OperationResult, error) {
	return crdb.Model(model).Create(ctx, model)
}

func (crdb *CRDBase) CreateOrUpdate(ctx context.Context, model any) (controllerutil.OperationResult, error) {
	return crdb.Model(model).CreateOrUpdate(ctx, model)
}

func (crdb *CRDBase) CreateOrPatch(ctx context.Context, model any) (controllerutil.OperationResult, error) {
	return crdb.Model(model).CreateOrPatch(ctx, model)
}

// // Delete deletes the given object by name from datastore.
// func (crdb *CRDBase) Delete(ctx context.Context, name string) error {
// 	return crdb.Model(model).Delete(ctx, name)
// }

// func (crdb *CRDBase) Get(ctx context.Context, query query.Query, out any) error {
// 	return crdb.Model(model).Get(ctx, model)

// }
// func (crdb *CRDBase) List(ctx context.Context, query query.Query, out any) error {
// 	return crdb.Model(model).List(ctx, model)
// }
// func (crdb *CRDBase) DeleteAllOf(ctx context.Context, query query.Query) error {
// 	return crdb.Model(model).DeleteAllOf(ctx, model)
// }
