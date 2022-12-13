/*
 * Copyright contributors to the Hyperledger Fabric Operator project
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 * 	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package override

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) AdminClusterRoleBinding(object v1.Object, crb *rbacv1.ClusterRoleBinding, action resources.Action) error {
	instance := object.(*current.Organization)
	switch action {
	case resources.Create, resources.Update:
		return o.CreateAdminClusterRoleBinding(instance, crb)
	}

	return nil
}

func (o *Override) CreateAdminClusterRoleBinding(instance *current.Organization, crb *rbacv1.ClusterRoleBinding) error {
	crb.Name = instance.GetRoleBinding(instance.GetAdminClusterRole())

	crb.Subjects = []rbacv1.Subject{
		common.GetDefaultSubject(instance.Spec.Admin, instance.GetUserNamespace(), o.SubjectKind),
	}

	crb.RoleRef = rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     instance.GetAdminClusterRole(),
	}

	crb.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Organization",
			APIVersion: "ibp.com/v1beta1",
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
		},
	}

	return nil
}
