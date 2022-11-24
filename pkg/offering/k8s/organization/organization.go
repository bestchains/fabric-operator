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

package k8sorg

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	baseorg "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ baseorg.Organization = &Organization{}

type Organization struct {
	*baseorg.BaseOrganization
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config) *Organization {
	organization := &Organization{
		BaseOrganization: baseorg.New(client, scheme, config),
	}
	return organization
}

func (organization *Organization) Reconcile(instance *current.Organization, update baseorg.Update) (common.Result, error) {
	var err error

	if err = organization.PreReconcileChecks(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed on prereconcile checks")
	}

	if err = organization.Initialize(instance, update); err != nil {
		return common.Result{}, operatorerrors.Wrap(err, operatorerrors.OrganizationInitilizationFailed, "failed to initialize organization")
	}

	if err = organization.ReconcileManagers(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}

	return organization.CheckStates(instance)
}

// TODO: customize for kubernetes

// PreReconcileChecks on Organization
func (organization *Organization) PreReconcileChecks(instance *current.Organization, update baseorg.Update) error {
	return organization.BaseOrganization.PreReconcileChecks(instance, update)
}

// Initialize on Organization after PreReconcileChecks
func (organization *Organization) Initialize(instance *current.Organization, update baseorg.Update) error {
	return organization.BaseOrganization.Initialize(instance, update)
}

// ReconcileManagers on Organization after Initialize
func (organization *Organization) ReconcileManagers(instance *current.Organization, update baseorg.Update) error {
	return organization.BaseOrganization.ReconcileManagers(instance, update)
}

// CheckStates on Organization after ReconcileManagers
func (organization *Organization) CheckStates(instance *current.Organization) (common.Result, error) {
	return organization.BaseOrganization.CheckStates(instance)
}
