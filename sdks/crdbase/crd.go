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
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	xPreserveUnknownFields = true
)

// NewCustomResourceDefinition returns a new CRD object.
func (crdb *CrdBase) NewCustomResourceDefinition(
	names apiextv1.CustomResourceDefinitionNames,
	schema map[string]apiextv1.JSONSchemaProps,
) *apiextv1.CustomResourceDefinition {
	return &apiextv1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: names.Plural + "." + crdb.GroupVersion.Group,
			Annotations: map[string]string{
				crdBaseURL + "/version":     crdBasePackage,
				crdBaseURL + "/api-version": crdApiVersion,
			},
		},
		Spec: apiextv1.CustomResourceDefinitionSpec{
			Group: crdb.GroupVersion.Group,
			Names: names,
			Scope: apiextv1.NamespaceScoped,
			Versions: []apiextv1.CustomResourceDefinitionVersion{
				{
					Name:    crdApiVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextv1.JSONSchemaProps{
							Description:            "Auto Generated schema for: " + names.Kind,
							Type:                   "object",
							Properties:             schema,
							XPreserveUnknownFields: &xPreserveUnknownFields,
						},
					},
					Subresources: &apiextv1.CustomResourceSubresources{
						Status: &apiextv1.CustomResourceSubresourceStatus{},
					},
					// AdditionalPrinterColumns: []apiextv1.CustomResourceColumnDefinition{},
				},
			},
			// PreserveUnknownFields: true,
		},
	}
}

// NewRBACRolesAndBindings return roles and its bindings.
func (crdb *CrdBase) NewRBACRolesAndBindings(
	names apiextv1.CustomResourceDefinitionNames,
) ([]*rbacv1.Role, []*rbacv1.RoleBinding) {
	managerRole := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.Singular + "-manager-role",
			Namespace: crdb.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{crdb.GroupVersion.Group},
				Resources: []string{names.Plural},
				Verbs: []string{
					"create",
					"delete",
					"get",
					"list",
					"patch",
					"update",
					"watch",
				},
			}, {
				APIGroups: []string{crdb.GroupVersion.Group},
				Resources: []string{names.Plural + "/finalizers"},
				Verbs: []string{
					"update",
				},
			}, {
				APIGroups: []string{crdb.GroupVersion.Group},
				Resources: []string{names.Plural + "/status"},
				Verbs: []string{
					"get",
					"patch",
					"update",
				},
			},
		},
	}
	managerRoleBinding := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "RoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.Singular + "-manager-rolebinding",
			Namespace: crdb.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/component":  "rbac",
				"app.kubernetes.io/created-by": names.Singular,
				"app.kubernetes.io/instance":   "manager-rolebinding",
				"app.kubernetes.io/managed-by": crdBaseName,
				"app.kubernetes.io/name":       "rolebinding",
				"app.kubernetes.io/part-of":    names.Singular,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     managerRole.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      crdb.ServiceAccount,
				Namespace: crdb.Namespace,
			},
		},
	}

	return []*rbacv1.Role{managerRole}, []*rbacv1.RoleBinding{managerRoleBinding}
}
