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

package operatorconfig

import (
	cainit "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/ca"
	fedinit "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/federation"
	ordererinit "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/orderer"
	orginit "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/organization"
	peerinit "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/peer"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering"
	"github.com/go-logr/logr"
)

type Config struct {
	CAInitConfig           *cainit.Config
	PeerInitConfig         *peerinit.Config
	OrdererInitConfig      *ordererinit.Config
	ConsoleInitConfig      *ConsoleConfig
	OrganizationInitConfig *orginit.Config
	FederationInitConfig   *fedinit.Config
	Offering               offering.Type
	Operator               Operator
	Logger                 *logr.Logger
}

type ConsoleConfig struct {
	DeploymentFile           string
	NetworkPolicyIngressFile string
	NetworkPolicyDenyAllFile string
	ServiceFile              string
	DeployerServiceFile      string
	PVCFile                  string
	CMFile                   string
	ConsoleCMFile            string
	DeployerCMFile           string
	RoleFile                 string
	RoleBindingFile          string
	ServiceAccountFile       string
	IngressFile              string
	Ingressv1beta1File       string
	RouteFile                string
}
