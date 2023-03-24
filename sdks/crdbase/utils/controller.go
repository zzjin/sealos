package utils

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateOrUpdateWithRetry(ctx context.Context, client client.Client, obj *unstructured.Unstructured, f controllerutil.MutateFn) (string, controllerutil.OperationResult, error) {
	var result controllerutil.OperationResult
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var err error
		result, err = controllerutil.CreateOrUpdate(ctx, client, obj, func() error {
			return f()
		})
		return err
	}); err != nil {
		return "", result, err
	}
	return obj.GetName(), result, nil
}

func UpdateWithRetry(ctx context.Context, c client.Client, obj *unstructured.Unstructured) (string, controllerutil.OperationResult, error) {
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var exObj unstructured.Unstructured
		exObj.SetGroupVersionKind(obj.GroupVersionKind())

		if err := c.Get(ctx, client.ObjectKeyFromObject(obj), &exObj); err != nil {
			return err
		}
		obj.SetResourceVersion(exObj.GetResourceVersion())
		if err := c.Update(ctx, obj); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return obj.GetName(), controllerutil.OperationResultNone, err
	}
	return obj.GetName(), controllerutil.OperationResultUpdated, nil
}
