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

func (o *Override) ClusterRoleBinding(object v1.Object, crb *rbacv1.ClusterRoleBinding, action resources.Action) error {
	instance := object.(*current.Federation)
	switch action {
	case resources.Create:
		return o.CreateClusterRoleBinding(instance, crb)
	case resources.Update:
		return o.UpdateClusterRoleBinding(instance, crb)
	}

	return nil
}

func (o *Override) CreateClusterRoleBinding(instance *current.Federation, crb *rbacv1.ClusterRoleBinding) error {
	subjects := make([]rbacv1.Subject, 0, len(instance.GetMembers()))

	for _, member := range instance.GetMembers() {
		organization, err := o.GetOrganization(member)
		if err != nil {
			return err
		}
		subjects = append(subjects, common.GetDefaultSubject(organization.Spec.Admin, organization.Namespace, o.SubjectKind))
	}

	crb.Subjects = subjects

	crb.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Federation",
			APIVersion: "ibp.com/v1beta1",
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
		},
	}

	return nil
}

// TODO: same as Create for now
func (o *Override) UpdateClusterRoleBinding(instance *current.Federation, crb *rbacv1.ClusterRoleBinding) error {
	subjects := make([]rbacv1.Subject, 0, len(instance.GetMembers()))

	for _, member := range instance.GetMembers() {
		organization, err := o.GetOrganization(member)
		if err != nil {
			return err
		}
		subjects = append(subjects, common.GetDefaultSubject(organization.Spec.Admin, organization.Namespace, o.SubjectKind))
	}

	crb.Subjects = subjects

	crb.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Federation",
			APIVersion: "ibp.com/v1beta1",
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
		},
	}

	return nil
}
