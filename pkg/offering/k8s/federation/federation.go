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

package k8sfed

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	basefed "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/federation"
	baseoverride "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/federation/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/federation/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ basefed.Federation = &Federation{}

type Federation struct {
	BaseFederation basefed.Federation
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config) *Federation {
	o := &override.Override{
		BaseOverride: &baseoverride.Override{
			Client: client,
		},
	}
	federation := &Federation{
		BaseFederation: basefed.New(client, scheme, config, o),
	}
	return federation
}

func (federation *Federation) Reconcile(instance *current.Federation, update basefed.Update) (common.Result, error) {
	var err error

	if err = federation.PreReconcileChecks(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed on prereconcile checks")
	}

	if err = federation.Initialize(instance, update); err != nil {
		return common.Result{}, operatorerrors.Wrap(err, operatorerrors.FederationInitilizationFailed, "failed to initialize federation")
	}

	if err = federation.ReconcileManagers(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}

	return federation.CheckStates(instance)
}

// TODO: customize for kubernetes

// PreReconcileChecks on Federation
func (federation *Federation) PreReconcileChecks(instance *current.Federation, update basefed.Update) error {
	return federation.BaseFederation.PreReconcileChecks(instance, update)
}

// Initialize on Federation after PreReconcileChecks
func (federation *Federation) Initialize(instance *current.Federation, update basefed.Update) error {
	return federation.BaseFederation.Initialize(instance, update)
}

// ReconcileManagers on Federation after Initialize
func (federation *Federation) ReconcileManagers(instance *current.Federation, update basefed.Update) error {
	return federation.BaseFederation.ReconcileManagers(instance, update)
}

// CheckStates on Federation after ReconcileManagers
func (federation *Federation) CheckStates(instance *current.Federation) (common.Result, error) {
	return federation.BaseFederation.CheckStates(instance)
}
