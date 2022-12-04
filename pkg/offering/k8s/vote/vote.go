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

package k8svote

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	basevote "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/vote"
	basevoteoverride "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/vote/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/vote/override"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("k8s_vote")

type Override interface {
	basevote.Override
}

type Vote struct {
	*basevote.BaseVote
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config) *Vote {
	o := &override.Override{
		Override: &basevoteoverride.Override{
			Client: client,
		},
	}
	vote := &Vote{
		BaseVote: basevote.New(client, scheme, config, o),
	}

	return vote
}

func (c *Vote) Reconcile(instance *current.Vote) (common.Result, error) {
	return c.BaseVote.Reconcile(instance)
}
