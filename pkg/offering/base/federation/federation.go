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
	"context"
	"fmt"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	"github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	bcrbac "github.com/IBM-Blockchain/fabric-operator/pkg/rbac"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("base_federation")

//go:generate counterfeiter -o mocks/update.go -fake-name Update . Update

type Update interface {
	SpecUpdated() bool
	MemberUpdated() bool
	ProposalActivated() bool
	ProposalFailed() bool
	ProposalDissolved() bool
}

//go:generate counterfeiter -o mocks/override.go -fake-name Override . Override

type Override interface{}

//go:generate counterfeiter -o mocks/basefederation.go -fake-name Federation . Federation

type Federation interface {
	PreReconcileChecks(instance *current.Federation, update Update) error
	Initialize(instance *current.Federation, update Update) error
	ReconcileManagers(instance *current.Federation, update Update) error
	CheckStates(instance *current.Federation, update Update) (common.Result, error)
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

	RBACManager *bcrbac.Manager
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

func (federation *BaseFederation) CreateManagers() {
	federation.RBACManager = bcrbac.NewRBACManager(federation.Client, nil)
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

	return federation.CheckStates(instance, update)
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
	if update.MemberUpdated() {
		err := federation.RBACManager.Reconcile(bcrbac.Federation, instance, bcrbac.ResourceUpdate)
		if err != nil {
			return err
		}

		networkList := &current.NetworkList{}
		if err := federation.Client.List(context.TODO(), networkList); err != nil {
			log.Error(err, fmt.Sprintf("failed to list networks by selector %s=%s",
				current.NETWORK_FEDERATION_LABEL, instance.GetName()))
			return err
		}
		log.Info(fmt.Sprintf("Listnetwork1 by selector %s=%s, get %+v",
			current.NETWORK_FEDERATION_LABEL, instance.GetName(), networkList.Items))
		log.Info(fmt.Sprintf("sync federation %s members", instance.GetName()))
		for i, n := range networkList.Items {
			if n.Labels == nil || n.Labels[current.NETWORK_FEDERATION_LABEL] != instance.GetName() {
				continue
			}
			networkList.Items[i].MergeMembers(instance.Spec.Members)
			log.Info(fmt.Sprintf("merge network %s' members", n.GetName()))
			if err = federation.Client.Patch(context.TODO(), &networkList.Items[i], nil, controllerclient.PatchOption{
				Resilient: &controllerclient.ResilientPatch{
					Retry:    3,
					Into:     &current.Network{},
					Strategy: client.MergeFrom,
				}}); err != nil {
				log.Error(err, fmt.Sprintf("federtion %s member update, failed to update network %s's member", instance.GetName(), n.GetName()))
			}
		}
	}
	return nil
}

// ReconcileRBAC will sync current federation to every member's AdminClusterRole
func (federation *BaseFederation) ReconcileRBAC() error {
	return nil
}

// CheckStates on Federation
func (federation *BaseFederation) CheckStates(instance *current.Federation, update Update) (common.Result, error) {
	status := instance.Status.CRStatus
	if !instance.HasType() {
		status.Type = current.FederationPending
		status.Status = current.True
	}
	if update.ProposalActivated() {
		status.Type = current.FederationActivated
		status.Status = current.True
	} else if update.ProposalFailed() {
		status.Type = current.FederationFailed
		status.Status = current.True
	} else if update.ProposalDissolved() {
		status.Type = current.FederationDissolved
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
