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

package federation

import (
	"fmt"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	"github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/manager"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("base_federation")

//go:generate counterfeiter -o mocks/update.go -fake-name Update . Update

type Update interface {
	SpecUpdated() bool
	MemberUpdated() bool
}

//go:generate counterfeiter -o mocks/override.go -fake-name Override . Override

type Override interface {
	ClusterRole(v1.Object, *rbacv1.ClusterRole, resources.Action) error
	ClusterRoleBinding(v1.Object, *rbacv1.ClusterRoleBinding, resources.Action) error
}

//go:generate counterfeiter -o mocks/basefederation.go -fake-name Federation . Federation

type Federation interface {
	PreReconcileChecks(instance *current.Federation, update Update) error
	Initialize(instance *current.Federation, update Update) error
	ReconcileManagers(instance *current.Federation, update Update) error
	CheckStates(instance *current.Federation) (common.Result, error)
	Reconcile(instance *current.Federation, update Update) (common.Result, error)
}

var _ Federation = (*BaseFederation)(nil)

const (
	KIND = "FEDERATION"
)

type BaseFederation struct {
	Client controllerclient.Client
	Scheme *runtime.Scheme

	Config *config.Config

	Override Override

	ClusterRoleManager        resources.Manager
	ClusterRoleBindingManager resources.Manager
}

func New(client controllerclient.Client, scheme *runtime.Scheme, config *config.Config, o Override) *BaseFederation {
	base := &BaseFederation{
		Client:   client,
		Scheme:   scheme,
		Config:   config,
		Override: o,
	}

	base.CreateManagers()

	return base
}

// TODO: leave this due to we might need managers in the future
// - configmap manager
func (federation *BaseFederation) CreateManagers() {
	override := federation.Override
	mgr := manager.New(federation.Client, federation.Scheme)

	federation.ClusterRoleManager = mgr.CreateClusterRoleManager("", override.ClusterRole, federation.GetLabels, federation.Config.FederationInitConfig.ClusterRoleFile)
	federation.ClusterRoleBindingManager = mgr.CreateClusterRoleBindingManager("", override.ClusterRoleBinding, federation.GetLabels, federation.Config.FederationInitConfig.ClusterRoleBindingFile)
}

// Reconcile on Federation upon Update
func (federation *BaseFederation) Reconcile(instance *current.Federation, update Update) (common.Result, error) {
	var err error

	if err = federation.PreReconcileChecks(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed on prereconcile checks")
	}

	if err = federation.Initialize(instance, update); err != nil {
		return common.Result{}, operatorerrors.Wrap(err, operatorerrors.FederationInitilizationFailed, "failed to initialize federation")
	}

	// TODO: define managers
	if err = federation.ReconcileManagers(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}

	return federation.CheckStates(instance)
}

// PreReconcileChecks on Federation upon Update
func (federation *BaseFederation) PreReconcileChecks(instance *current.Federation, update Update) error {
	log.Info(fmt.Sprintf("PreReconcileChecks on Federation %s", instance.GetName()))

	if !instance.HasInitiator() {
		return errors.New("federation initiator is empty")
	}

	if instance.HasMultiInitiator() {
		return errors.New("federation only allow one initiator")
	}

	if !instance.HasPolicy() {
		return errors.New("federation policy is empty")
	}

	return nil
}

// Initialize on Federation upon Update
func (federation *BaseFederation) Initialize(instance *current.Federation, update Update) error {
	return nil
}

// ReconcileManagers on Federation upon Update
func (federation *BaseFederation) ReconcileManagers(instance *current.Federation, update Update) error {
	var err error

	// cluster role do not need to update
	if err = federation.ClusterRoleManager.Reconcile(instance, false); err != nil {
		return errors.Wrap(err, "reconcile cluster role")
	}

	if err = federation.ClusterRoleBindingManager.Reconcile(instance, update.MemberUpdated()); err != nil {
		return errors.Wrap(err, "reconcile cluster role binding")
	}

	return nil
}

// CheckStates on Federation
func (federation *BaseFederation) CheckStates(instance *current.Federation) (common.Result, error) {
	status := instance.Status.CRStatus
	if !instance.HasType() {
		status.Type = current.FederationPending
		status.Status = current.True
	}
	return common.Result{
		Status: &status,
	}, nil
}

// GetLabels from instance.GetLabels
func (federation *BaseFederation) GetLabels(instance v1.Object) map[string]string {
	return instance.GetLabels()
}
