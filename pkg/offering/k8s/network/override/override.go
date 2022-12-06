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
	rbacv1 "k8s.io/api/rbac/v1"
)

//go:generate counterfeiter -o mocks/baseOverride.go -fake-name BaseOverride . baseOverride
type baseOverride interface {
	CreateClusterRole(*current.Network, *rbacv1.ClusterRole) error
	UpdateClusterRole(*current.Network, *rbacv1.ClusterRole) error

	CreateClusterRoleBinding(*current.Network, *rbacv1.ClusterRoleBinding) error
	UpdateClusterRoleBinding(*current.Network, *rbacv1.ClusterRoleBinding) error
}

type Override struct {
	BaseOverride baseOverride
}
