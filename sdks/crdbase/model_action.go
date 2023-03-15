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

	"github.com/labring/crdbase/query"
	"github.com/labring/crdbase/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type CRDModelAction struct {
	CRDBase
	CRDModelSchema

	rdc dynamic.ResourceInterface
}

func (crdb CRDBase) Model(m Model) CRDModelAction {
	mSchema := GetCRDModelSchema(m)

	gvr := schema.GroupVersionResource{
		Group:    crdb.GroupVersion.Group,
		Version:  crdb.GroupVersion.Version,
		Resource: mSchema.ResourceName(),
	}

	return CRDModelAction{
		CRDBase:        crdb,
		CRDModelSchema: *mSchema,

		rdc: crdb.dynamicClient.Resource(gvr).Namespace(crdb.Namespace),
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

	if _, err := crdms.rdc.Create(ctx, cr, metav1.CreateOptions{}); err != nil {
		return "", controllerutil.OperationResultNone, err
	}
	return uniqName, controllerutil.OperationResultCreated, nil
}

func (crdms CRDModelAction) CreateOrUpdate(ctx context.Context, model any) (string, controllerutil.OperationResult, error) {
	uniqName := crdms.GetPrimaryFieldValue(model)
	if uniqName == "" {
		return crdms.Create(ctx, model)
	}

	// Unstructured For create CR
	cr, err := crdms.model2UnstructuredCR(uniqName, model)
	if err != nil {
		return "", controllerutil.OperationResultNone, err
	}

	if _, err := crdms.rdc.Get(ctx, uniqName, metav1.GetOptions{}); err != nil {
		if !apierrors.IsNotFound(err) {
			return "", controllerutil.OperationResultNone, err
		}
		if _, err := crdms.rdc.Create(ctx, cr, metav1.CreateOptions{}); err != nil {
			return "", controllerutil.OperationResultNone, err
		}
		return uniqName, controllerutil.OperationResultCreated, nil
	}

	if _, err := crdms.rdc.Update(ctx, cr, metav1.UpdateOptions{}); err != nil {
		return "", controllerutil.OperationResultNone, err
	}

	return uniqName, controllerutil.OperationResultUpdated, nil
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
				"name": name,
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

// Delete deletes the given object by name from datastore.
func (crdms CRDModelAction) Delete(ctx context.Context, name string) error {
	return crdms.rdc.Delete(ctx, name, metav1.DeleteOptions{})
}

func (crdms CRDModelAction) DeleteAllOf(ctx context.Context, query query.Query) error {
	return nil
}

func (crdms CRDModelAction) Get(ctx context.Context, query query.Query, out any) error {
	opt := metav1.ListOptions{}

	gotList, err := crdms.rdc.List(ctx, opt)
	if err != nil {
		return err
	}

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
