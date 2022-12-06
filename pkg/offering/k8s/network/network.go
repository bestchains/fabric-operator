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

package k8snet

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	basenet "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/network"
	baseoverride "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/network/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/network/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ basenet.Network = &Network{}

type Network struct {
	BaseNetwork basenet.Network
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config) *Network {
	o := &override.Override{
		BaseOverride: &baseoverride.Override{
			Client: client,
		},
	}
	network := &Network{
		BaseNetwork: basenet.New(client, scheme, config, o),
	}
	return network
}

func (network *Network) Reconcile(instance *current.Network, update basenet.Update) (common.Result, error) {
	var err error

	if err = network.PreReconcileChecks(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed on prereconcile checks")
	}

	if err = network.Initialize(instance, update); err != nil {
		return common.Result{}, operatorerrors.Wrap(err, operatorerrors.NetworkInitializationFailed, "failed to initialize network")
	}

	if err = network.ReconcileManagers(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}

	return network.CheckStates(instance)
}

// TODO: customize for kubernetes

// PreReconcileChecks on Network
func (network *Network) PreReconcileChecks(instance *current.Network, update basenet.Update) error {
	return network.BaseNetwork.PreReconcileChecks(instance, update)
}

// Initialize on Network after PreReconcileChecks
func (network *Network) Initialize(instance *current.Network, update basenet.Update) error {
	return network.BaseNetwork.Initialize(instance, update)
}

// ReconcileManagers on Network after Initialize
func (network *Network) ReconcileManagers(instance *current.Network, update basenet.Update) error {
	return network.BaseNetwork.ReconcileManagers(instance, update)
}

// CheckStates on Network after ReconcileManagers
func (network *Network) CheckStates(instance *current.Network) (common.Result, error) {
	return network.BaseNetwork.CheckStates(instance)
}
