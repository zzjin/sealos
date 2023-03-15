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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type CRDModelAction struct {
	CRDBase
	CRDModelSchema
}

func (crdb CRDBase) Model(m Model) CRDModelAction {
	schema := GetCRDModelSchema(m)

	return CRDModelAction{
		CRDBase:        crdb,
		CRDModelSchema: *schema,
	}
}

func (crdms CRDModelAction) getResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    crdms.GroupVersion.Group,
		Version:  crdms.GroupVersion.Version,
		Resource: crdms.ResourceName(),
	}
}

func (crdms CRDModelAction) Create(ctx context.Context, model any) (string, controllerutil.OperationResult, error) {
	uniqName := crdms.GetPrimaryFieldValue(model)

	// Unstructured For create CR
	cr, err := crdms.model2UnstructuredCR(uniqName, model)
	if err != nil {
		return "", controllerutil.OperationResultNone, err
	}

	rs := crdms.dynamicClient.Resource(crdms.getResource())

	if _, err := rs.Get(ctx, uniqName, metav1.GetOptions{}); err != nil {
		if !apierrors.IsNotFound(err) {
			return "", controllerutil.OperationResultNone, err
		}
		if _, err := rs.Create(ctx, cr, metav1.CreateOptions{}); err != nil {
			return "", controllerutil.OperationResultNone, err
		}
		return uniqName, controllerutil.OperationResultCreated, nil
	}

	if _, err := rs.Update(ctx, cr, metav1.UpdateOptions{}); err != nil {
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
				"labels": map[string]any{
					crdBaseURL + "/managed-by": providerName,
				},
			},
			"spec": modelMap,
		},
	}

	y, _ := yaml.Marshal(mcr)
	fmt.Println(string(y))

	return mcr, nil
}

func (crdms CRDModelAction) CreateOrUpdate(ctx context.Context, model any) (controllerutil.OperationResult, error) {
	return controllerutil.OperationResultNone, nil
}

func (crdms CRDModelAction) CreateOrPatch(ctx context.Context, model any) (controllerutil.OperationResult, error) {
	return controllerutil.OperationResultNone, nil
}

// Delete deletes the given object by name from datastore.
func (crdms CRDModelAction) Delete(ctx context.Context, name string) error {
	resource := crdms.dynamicClient.Resource(crdms.getResource()).Namespace(crdms.Namespace)
	return resource.Delete(ctx, name, metav1.DeleteOptions{})
}

func (crdms CRDModelAction) Get(ctx context.Context, query query.Query, out any) error {
	// names := Model2KindName(out)
	// gvr := schema.GroupVersionResource{
	// 	Group:    crdb.GroupVersion.Group,
	// 	Version:  crdb.GroupVersion.Version,
	// 	Resource: names.Plural,
	// }

	// resource := crdb.dynamicClient.Resource(gvr).Namespace(crdb.Namespace)

	// opt := metav1.ListOptions{}

	// gots, err := resource.List(ctx, opt)
	// if err != nil {
	// 	return err
	// }

	// if len(gots.Items) == 0 {
	// 	return nil
	// }

	// got := gots.Items[0]

	// if err := utils.Map2JSONStruct(got.UnstructuredContent(), &out); err != nil {
	// 	return fmt.Errorf("failed to convert map to struct: %w", err)
	// }

	return nil

}
func (crdms CRDModelAction) List(ctx context.Context, query query.Query, out any) error {
	return nil
}
func (crdms CRDModelAction) DeleteAllOf(ctx context.Context, query query.Query) error {
	return nil
}
