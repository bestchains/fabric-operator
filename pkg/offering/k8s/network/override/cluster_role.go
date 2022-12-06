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
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) ClusterRole(object v1.Object, clusterRole *rbacv1.ClusterRole, action resources.Action) error {
	instance := object.(*current.Network)
	switch action {
	case resources.Create:
		return o.CreateClusterRole(instance, clusterRole)
	case resources.Update:
		return o.UpdateClusterRole(instance, clusterRole)
	}

	return nil
}

func (o *Override) CreateClusterRole(instance *current.Network, cr *rbacv1.ClusterRole) error {
	// No need to override `BaseOverride`
	return o.BaseOverride.CreateClusterRole(instance, cr)
}

func (o *Override) UpdateClusterRole(instance *current.Network, cr *rbacv1.ClusterRole) error {
	// No need to override `BaseOverride`
	return o.BaseOverride.UpdateClusterRole(instance, cr)
}
