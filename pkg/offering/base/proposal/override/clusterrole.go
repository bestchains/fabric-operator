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
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) ClusterRole(object v1.Object, clusterRole *rbacv1.ClusterRole, action resources.Action) error {
	instance := object.(*current.Proposal)
	switch action {
	case resources.Create, resources.Update:
		return o.SyncClusterRole(instance, clusterRole)
	}

	return nil
}

func (o *Override) SyncClusterRole(instance *current.Proposal, cr *rbacv1.ClusterRole) error {
	rule := rbacv1.PolicyRule{
		APIGroups:     []string{client.CRDGroup},
		Resources:     []string{"proposals"},
		ResourceNames: []string{instance.GetName()},
		Verbs:         []string{"get"},
	}

	cr.Rules = append(cr.Rules, rule)

	cr.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Proposal",
			APIVersion: client.SchemeGroupVersion.String(),
			Name:       instance.GetName(),
		},
	}

	return nil
}