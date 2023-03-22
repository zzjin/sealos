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
	"time"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AutoMigrate auto apply crd(s) to spec models
func (crdb *CrdBase) AutoMigrate(ctx context.Context, models ...any) error {
	crds, err := crdb.generateCRDs(models)
	if err != nil {
		return fmt.Errorf("unable to generate crds: %w", err)
	}

	// crdb.addToSchemes(crds)

	if err := crdb.installCRDs(ctx, crds); err != nil {
		return fmt.Errorf("unable to install crds: %w", err)
	}

	if err := crdb.waitCRDs(ctx, crds); err != nil {
		return fmt.Errorf("unable to ensure crds: %w", err)
	}

	namess := crdb.getNamesByCRDs(crds)
	if err := crdb.applyRBAC(ctx, namess); err != nil {
		return fmt.Errorf("unable to apply rbac: %w", err)
	}

	// TODO: add indexes

	return nil
}

func (crdb *CrdBase) generateCRDs(models []any) ([]*apiextv1.CustomResourceDefinition, error) {
	var crds []*apiextv1.CustomResourceDefinition

	for _, model := range models {
		crd, err := crdb.Model2CRD(model)
		if err != nil {
			return nil, err
		}
		crds = append(crds, crd)
	}

	return crds, nil
}

// func (crdb *CrdBase) addToSchemes(crds []*apiextv1.CustomResourceDefinition) {
// 	sch := crdb.Manager.GetScheme()

// 	for _, crd := range crds {
// 		gv := schema.GroupVersion{
// 			Group:   crd.Spec.Group,
// 			Version: crd.Spec.Versions[0].Name,
// 		}
// 		schemeBuilder := &scheme.Builder{GroupVersion: gv}
// 		schemeBuilder.AddToScheme(sch)
// 	}
// }

func (crdb *CrdBase) getNamesByCRDs(crds []*apiextv1.CustomResourceDefinition) []apiextv1.CustomResourceDefinitionNames {
	var namess []apiextv1.CustomResourceDefinitionNames

	for _, crd := range crds {
		namess = append(namess, crd.Spec.Names)
	}

	return namess
}

// Prune Delete all crd(s) from spec models along with all the stored crs
func (crdb *CrdBase) Prune(ctx context.Context, models ...any) error {
	crds, err := crdb.generateCRDs(models)
	if err != nil {
		return fmt.Errorf("unable to generate crds: %w", err)
	}

	//Since we do not have any finalizers, just delete the CRDs
	if err := crdb.deleteCRDs(ctx, crds); err != nil {
		return fmt.Errorf("unable to delete crds: %w", err)
	}

	namess := crdb.getNamesByCRDs(crds)
	if err := crdb.deleteRBAC(ctx, namess); err != nil {
		return fmt.Errorf("unable to delete rbac: %w", err)
	}
	return nil
}

func (crdb *CrdBase) installCRDs(ctx context.Context, crds []*apiextv1.CustomResourceDefinition) error {
	for _, crd := range crds {
		crdb.log.V(1).Info("installing CRD", "crd", crd.GetName())
		existingCrd := crd.DeepCopy()
		errGet := crdb.client.Get(ctx, client.ObjectKey{Name: crd.GetName()}, existingCrd)
		switch {
		case apierrors.IsNotFound(errGet):
			if err := crdb.client.Create(ctx, crd); err != nil {
				return fmt.Errorf("unable to create CRD %q: %w", crd.GetName(), err)
			}
		case errGet != nil:
			return fmt.Errorf("unable to get CRD %q to check if it exists: %w", crd.GetName(), errGet)
		default:
			crdb.log.V(1).Info("CRD already exists, updating", "crd", crd.GetName())
			if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
				if err := crdb.client.Get(ctx, client.ObjectKey{Name: crd.GetName()}, existingCrd); err != nil {
					return err
				}
				crd.SetResourceVersion(existingCrd.GetResourceVersion())
				return crdb.client.Update(ctx, crd)
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

var (
	PollInterval = 1 * time.Second
	MaxWait      = 10 * time.Second
)

// poller checks if all the resources have been found in discovery, and returns false if not.
type poller struct {
	// config is used to get discovery
	config *rest.Config

	// waitingFor is the map of resources keyed by group version that have not yet been found in discovery
	waitingFor map[schema.GroupVersion]*sets.String
}

// poll checks if all the resources have been found in discovery, and returns false if not.
func (p *poller) poll() (done bool, err error) {
	// Create a new clientset to avoid any client caching of discovery
	cs, err := clientset.NewForConfig(p.config)
	if err != nil {
		return false, err
	}

	allFound := true
	for gv, resources := range p.waitingFor {
		// All resources found, do nothing
		if resources.Len() == 0 {
			delete(p.waitingFor, gv)
			continue
		}

		// Get the Resources for this GroupVersion
		resourceList, err := cs.Discovery().ServerResourcesForGroupVersion(gv.String())
		if err != nil {
			return false, nil //nolint:nilerr
		}

		// Remove each found resource from the resources set that we are waiting for
		for _, resource := range resourceList.APIResources {
			resources.Delete(resource.Name)
		}

		// Still waiting on some resources in this group version
		if resources.Len() != 0 {
			allFound = false
		}
	}
	return allFound, nil
}

func (crdb *CrdBase) waitCRDs(ctx context.Context, crds []*apiextv1.CustomResourceDefinition) error {
	// Add each CRD to a map of GroupVersion to Resource
	waitingFor := map[schema.GroupVersion]*sets.String{}
	for _, crd := range crds {
		var gvs []schema.GroupVersion
		for _, version := range crd.Spec.Versions {
			if version.Served {
				gvs = append(gvs, schema.GroupVersion{Group: crd.Spec.Group, Version: version.Name})
			}
		}

		for _, gv := range gvs {
			crdb.log.V(1).Info("adding API in waitlist", "GV", gv)
			if _, found := waitingFor[gv]; !found {
				// Initialize the set
				waitingFor[gv] = &sets.String{}
			}
			// Add the Resource
			waitingFor[gv].Insert(crd.Spec.Names.Plural)
		}
	}

	// Poll until all resources are found in discovery
	p := &poller{config: crdb.Manager.GetConfig(), waitingFor: waitingFor}
	return wait.PollImmediate(PollInterval, MaxWait, p.poll)
}

func (crdb *CrdBase) deleteCRDs(ctx context.Context, crds []*apiextv1.CustomResourceDefinition) error {
	// Uninstall each CRD
	for _, crd := range crds {
		crd := crd
		crdb.log.V(1).Info("uninstalling CRD", "crd", crd.GetName())
		if err := crdb.client.Delete(ctx, crd); err != nil {
			// If CRD is not found, we can consider success
			if !apierrors.IsNotFound(err) {
				return err
			}
		}
	}

	return nil
}

func (crdb *CrdBase) applyRBAC(ctx context.Context, namess []apiextv1.CustomResourceDefinitionNames) error {
	var err error

	roleClient := crdb.clientSet.RbacV1().Roles(crdb.Namespace)
	roleBidingClient := crdb.clientSet.RbacV1().RoleBindings(crdb.Namespace)

	// Uninstall each CRD Names role and rolebindings
	for _, names := range namess {
		roles, roleBindings := crdb.NewRBACRolesAndBindings(names)
		crdb.log.V(1).Info("uninstalling Role and RoleBinding", "name", names.Singular)

		for _, role := range roles {
			_, errGet := roleClient.Get(ctx, role.Name, metav1.GetOptions{})
			switch {
			case apierrors.IsNotFound(errGet):
				if _, err = roleClient.Create(ctx, role, metav1.CreateOptions{}); err != nil {
					return fmt.Errorf("unable to create Role %q: %w", role.Name, err)
				}
			case errGet != nil:
				return fmt.Errorf("unable to get Role %q: %w", role.Name, errGet)
			default:
				crdb.log.V(1).Info("Role already exists, updating", "role", role.Name)
				if _, err := roleClient.Update(ctx, role, metav1.UpdateOptions{}); err != nil {
					return err
				}
			}
		}

		for _, roleBinding := range roleBindings {
			_, errGet := roleBidingClient.Get(ctx, roleBinding.Name, metav1.GetOptions{})
			switch {
			case apierrors.IsNotFound(errGet):
				if _, err = roleBidingClient.Create(ctx, roleBinding, metav1.CreateOptions{}); err != nil {
					return fmt.Errorf("unable to create RoleBiding %q: %w", roleBinding.Name, err)
				}
			case errGet != nil:
				return fmt.Errorf("unable to get RoleBiding %q: %w", roleBinding.Name, errGet)
			default:
				crdb.log.V(1).Info("RoleBiding already exists, updating", "rolebiding", roleBinding.Name)
				if _, err := roleBidingClient.Update(ctx, roleBinding, metav1.UpdateOptions{}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (crdb *CrdBase) deleteRBAC(ctx context.Context, namess []apiextv1.CustomResourceDefinitionNames) error {
	roleClient := crdb.clientSet.RbacV1().Roles(crdb.Namespace)
	roleBidingClient := crdb.clientSet.RbacV1().RoleBindings(crdb.Namespace)

	// Uninstall each CRD Names role and rolebindings
	for _, names := range namess {
		roles, roleBindings := crdb.NewRBACRolesAndBindings(names)
		crdb.log.V(1).Info("uninstalling Role and RoleBinding", "name", names.Singular)

		for _, role := range roles {
			role := role
			if err := roleClient.Delete(ctx, role.Name, metav1.DeleteOptions{}); err != nil {
				// If is not found, we can consider success
				if !apierrors.IsNotFound(err) {
					return err
				}
			}
		}

		for _, roleBinding := range roleBindings {
			roleBinding := roleBinding
			if err := roleBidingClient.Delete(ctx, roleBinding.Name, metav1.DeleteOptions{}); err != nil {
				// If is not found, we can consider success
				if !apierrors.IsNotFound(err) {
					return err
				}
			}
		}
	}

	return nil
}
