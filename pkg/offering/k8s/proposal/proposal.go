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

package k8sproposal

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	baseproposal "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/proposal"
	baseproposaloverride "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/proposal/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/proposal/override"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("k8s_proposal")

type Override interface {
	baseproposal.Override
}

type Proposal struct {
	*baseproposal.BaseProposal

	Override Override
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config) *Proposal {
	o := &override.Override{
		Override: &baseproposaloverride.Override{
			Client: client,
		},
	}

	Proposal := &Proposal{
		BaseProposal: baseproposal.New(client, scheme, config, o),
		Override:     o,
	}

	Proposal.CreateManagers()
	return Proposal
}

func (c *Proposal) Reconcile(instance *current.Proposal) (common.Result, error) {
	return c.BaseProposal.Reconcile(instance)
}
