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
	"fmt"
	"reflect"

	"github.com/labring/crdbase/query"
	"github.com/labring/crdbase/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type CRDModelAction struct {
	CRDBase
	CRDModelSchema

	gvk schema.GroupVersionKind
}

// FIXME: real Impl.
type UpdateFunc func(cr Model) Model

func (crdb CRDBase) Model(m Model) CRDModelAction {
	mSchema := GetCRDModelSchema(m)

	return CRDModelAction{
		CRDBase:        crdb,
		CRDModelSchema: *mSchema,

		gvk: schema.GroupVersionKind{
			Group:   crdb.GroupVersion.Group,
			Version: crdb.GroupVersion.Version,
			Kind:    mSchema.Kind(),
		},
	}
}

func (crdms CRDModelAction) Create(ctx context.Context, model any) (string, controllerutil.OperationResult, error) {
	uniqName := crdms.GetPrimaryFieldValue(model)
	if uniqName == "" {
		uniqName = utils.GenNanoID()
	}

	// Unstructured For create CR
	cr, err := crdms.model2UnstructuredCR(uniqName, model)
	if err != nil {
		return "", controllerutil.OperationResultNone, err
	}

	if err := crdms.client.Create(ctx, cr); err != nil {
		return "", controllerutil.OperationResultNone, err
	}
	return uniqName, controllerutil.OperationResultCreated, nil
}

func (crdms CRDModelAction) CreateOrUpdate(ctx context.Context, model any, um UpdateFunc) (string, controllerutil.OperationResult, error) {
	if reflect.TypeOf(model).Kind() == reflect.Slice {
		return "", controllerutil.OperationResultNone, fmt.Errorf("model must be a pointer to a struct, not slice")
	}

	uniqName := crdms.GetPrimaryFieldValue(model)
	if uniqName == "" {
		return crdms.Create(ctx, model)
	}

	// Unstructured For create CR
	cr, err := crdms.model2UnstructuredCR(uniqName, model)
	if err != nil {
		return "", controllerutil.OperationResultNone, err
	}

	getCR := crdms.NewGetUnstructured()

	if err := crdms.client.Get(ctx, crdms.NamespacedName(uniqName), getCR); err != nil {
		if !apierrors.IsNotFound(err) {
			return "", controllerutil.OperationResultNone, err
		}
		if err := crdms.client.Create(ctx, cr); err != nil {
			return "", controllerutil.OperationResultNone, err
		}
		return uniqName, controllerutil.OperationResultCreated, nil
	}

	cr.SetResourceVersion(getCR.GetResourceVersion())

	// TODO: transfer data to model and back to unstructured.
	updatedCR := &unstructured.Unstructured{}

	if err := crdms.client.Update(ctx, updatedCR); err != nil {
		return "", controllerutil.OperationResultNone, err
	}

	return uniqName, controllerutil.OperationResultUpdated, nil
}

func (crdms CRDModelAction) CreateOrUpdateList(ctx context.Context, model any, um UpdateFunc) (string, controllerutil.OperationResult, error) {
	if reflect.TypeOf(model).Kind() != reflect.Slice {
		return "", controllerutil.OperationResultNone, fmt.Errorf("model must be a pointer to a struct, not slice")
	}

	return "", controllerutil.OperationResultNone, nil
}

func (crdms CRDModelAction) model2UnstructuredCR(name string, m Model) (*unstructured.Unstructured, error) {
	modelMap, err := utils.StructJSON2Map(m)
	if err != nil {
		return nil, fmt.Errorf("convert model to Unstructured fail: %w", err)
	}

	// Unstructured For create CR
	mcr := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": crdms.ApiVersion(),
			"kind":       crdms.Names.Kind,
			"metadata": map[string]any{
				"name":      name,
				"namespace": crdms.Namespace,
				// "labels": map[string]any{
				// 	crdBaseURL + "/managed-by": providerName,
				// },
			},
			"spec": modelMap,
		},
	}

	y, _ := yaml.Marshal(mcr)
	crdms.log.V(1).Info("model2UnstructuredCR", "unstructured", string(y))

	return mcr, nil
}

func (crdms CRDModelAction) NamespacedName(name string) types.NamespacedName {
	return types.NamespacedName{Namespace: crdms.Namespace, Name: name}
}

func (crdms CRDModelAction) NewGetUnstructured() *unstructured.Unstructured {
	un := &unstructured.Unstructured{}
	un.SetGroupVersionKind(crdms.gvk)
	return un
}

func (crdms CRDModelAction) NewGetUnstructuredList() *unstructured.UnstructuredList {
	unl := &unstructured.UnstructuredList{}
	unl.SetGroupVersionKind(crdms.gvk)
	return unl
}

// Delete deletes the given object by name from datastore.
func (crdms CRDModelAction) Delete(ctx context.Context, name string) error {
	deleteObj := crdms.NewGetUnstructured()
	deleteObj.SetNamespace(crdms.Namespace)
	deleteObj.SetName(name)

	return crdms.client.Delete(ctx, deleteObj)
}

func (crdms CRDModelAction) DeleteAllOf(ctx context.Context, query query.Query) error {
	return nil
}

func (crdms CRDModelAction) Get(ctx context.Context, q query.Query, out any) error {
	// TODO: real options
	opts := q.ToListOptions()

	gotList := crdms.NewGetUnstructuredList()

	err := crdms.client.List(ctx, gotList, opts...)
	if err != nil {
		return err
	}

	// TODO: Always filter response by query
	q.PostFilter(gotList)

	if len(gotList.Items) == 0 {
		return nil
	}

	got := gotList.Items[0]

	if err := utils.Map2JSONStruct(got.UnstructuredContent(), &out); err != nil {
		return fmt.Errorf("failed to convert map to struct: %w", err)
	}

	return nil

}
func (crdms CRDModelAction) List(ctx context.Context, query query.Query, out any) error {
	if _, _, err := utils.EnsureStructSlice(out); err != nil {
		return fmt.Errorf("out must be a slice to struct")
	}

	return nil
}
