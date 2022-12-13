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
	"github.com/IBM-Blockchain/fabric-operator/pkg/client"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) RoleBinding(object v1.Object, rb *rbacv1.RoleBinding, action resources.Action) error {
	instance := object.(*current.Vote)
	switch action {
	case resources.Create, resources.Update:
		return o.SyncRoleBinding(instance, rb)
	}

	return nil
}

func (o *Override) SyncRoleBinding(instance *current.Vote, rb *rbacv1.RoleBinding) error {
	org := instance.GetOrganization()
	organization, err := o.GetOrganization(org)
	if err != nil {
		return err
	}
	rb.Subjects = append(rb.Subjects, common.GetDefaultSubject(organization.Spec.Admin, organization.GetUserNamespace(), o.SubjectKind))
	rb.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Vote",
			APIVersion: client.SchemeGroupVersion.String(),
			Name:       instance.GetName(),
		},
	}

	return nil
}
