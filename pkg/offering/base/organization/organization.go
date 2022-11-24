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

package organization

import (
	"fmt"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	"github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/IBM-Blockchain/fabric-operator/version"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("base_organization")

//go:generate counterfeiter -o mocks/update.go -fake-name Update . Update

type Update interface {
	AdminOrCAUpdated() bool
}

type Override interface {
}

//go:generate counterfeiter -o mocks/initializer.go -fake-name InitializerOrganization . InitializerOrganization

type InitializerOrganization interface {
	CreateOrUpdateOrgMSPSecret(instance *current.Organization) error
}

type Organization interface {
	PreReconcileChecks(instance *current.Organization, update Update) error
	Initialize(instance *current.Organization, update Update) error
	ReconcileManagers(instance *current.Organization, update Update) error
	CheckStates(instance *current.Organization) (common.Result, error)
	Reconcile(instance *current.Organization, update Update) (common.Result, error)
}

var _ Organization = (*BaseOrganization)(nil)

const (
	KIND = "ORGANIZATION"
)

type BaseOrganization struct {
	Client controllerclient.Client
	Scheme *runtime.Scheme

	Config *config.Config

	Initializer InitializerOrganization
}

func New(client controllerclient.Client, scheme *runtime.Scheme, config *config.Config) *BaseOrganization {
	base := &BaseOrganization{
		Client: client,
		Scheme: scheme,
		Config: config,
	}

	base.Initializer = NewInitializer(config.OrganizationInitConfig, scheme, client, base.GetLabels)

	base.CreateManagers()

	return base
}

// TODO: leave this due to we might need managers in the future
// - configmap manager
func (organization *BaseOrganization) CreateManagers() {}

// Reconcile on Organization upon Update
func (organization *BaseOrganization) Reconcile(instance *current.Organization, update Update) (common.Result, error) {
	var err error

	if err = organization.PreReconcileChecks(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed on prereconcile checks")
	}

	if err = organization.Initialize(instance, update); err != nil {
		return common.Result{}, operatorerrors.Wrap(err, operatorerrors.OrganizationInitilizationFailed, "failed to initialize organization")
	}

	// TODO: define managers
	if err = organization.ReconcileManagers(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}

	return organization.CheckStates(instance)
}

// PreReconcileChecks on Organization upon Update
func (organization *BaseOrganization) PreReconcileChecks(instance *current.Organization, update Update) error {
	log.Info(fmt.Sprintf("PreReconcileChecks on Organization %s", instance.GetName()))

	if !instance.HasCARef() {
		return errors.New("organization caRef is empty")
	}

	if !instance.HasAdmin() {
		return errors.New("organization admin is empty")
	}

	return nil
}

// Initialize on Organization upon Update
func (organization *BaseOrganization) Initialize(instance *current.Organization, update Update) error {
	log.Info(fmt.Sprintf("Checking if organization '%s' needs initialization", instance.GetName()))

	if update.AdminOrCAUpdated() {
		err := organization.Initializer.CreateOrUpdateOrgMSPSecret(instance)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReconcileManagers on Organization upon Update
func (organization *BaseOrganization) ReconcileManagers(instance *current.Organization, update Update) error {
	return nil
}

// CheckStates on Organization
func (organization *BaseOrganization) CheckStates(instance *current.Organization) (common.Result, error) {
	return common.Result{
		Status: &current.CRStatus{
			Type:    current.Created,
			Version: version.Operator,
		},
	}, nil
}

// GetLabels from instance.GetLabels
func (organization *BaseOrganization) GetLabels(instance v1.Object) map[string]string {
	return instance.GetLabels()
}
