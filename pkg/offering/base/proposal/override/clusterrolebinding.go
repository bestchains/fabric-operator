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
	"context"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/client"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) ClusterRoleBinding(object v1.Object, crb *rbacv1.ClusterRoleBinding, action resources.Action) error {
	instance := object.(*current.Proposal)
	switch action {
	case resources.Create, resources.Update:
		return o.SyncClusterRoleBinding(instance, crb)
	}

	return nil
}

func (o *Override) SyncClusterRoleBinding(instance *current.Proposal, crb *rbacv1.ClusterRoleBinding) error {
	orgs, err := instance.GetCandidateOrganizations(context.TODO(), o.Client)
	if err != nil {
		return err
	}
	subjects := make([]rbacv1.Subject, 0)

	for _, org := range orgs {
		organization, err := o.GetOrganization(org)
		if err != nil {
			return err
		}
		subjects = append(subjects, common.GetDefaultSubject(organization.Spec.Admin, organization.Namespace, o.SubjectKind))
	}
	crb.Subjects = subjects

	crb.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Proposal",
			APIVersion: client.SchemeGroupVersion.String(),
			Name:       instance.GetName(),
		},
	}

	return nil
}
