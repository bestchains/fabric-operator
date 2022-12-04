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

package baseproposal

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	resourcemanager "github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/manager"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	k8sruntime "sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("base_proposal")

type Override interface {
	ClusterRole(v1.Object, *rbacv1.ClusterRole, resources.Action) error
	ClusterRoleBinding(v1.Object, *rbacv1.ClusterRoleBinding, resources.Action) error
}

//go:generate counterfeiter -o mocks/update.go -fake-name Update . Update

type Update interface {
}

type Proposal interface {
	PreReconcileChecks(instance *current.Proposal) (bool, error)
	ReconcileManagers(ctx context.Context, instance *current.Proposal) error
	Reconcile(instance *current.Proposal) (common.Result, error)
}

var _ Proposal = &BaseProposal{}

type BaseProposal struct {
	Client k8sclient.Client
	Scheme *runtime.Scheme
	Config *config.Config

	ClusterRoleManager        resources.Manager
	ClusterRoleBindingManager resources.Manager
	Override                  Override
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config, o Override) *BaseProposal {
	p := &BaseProposal{
		Client:   client,
		Scheme:   scheme,
		Config:   config,
		Override: o,
	}

	p.CreateManagers()
	return p
}

func (p *BaseProposal) CreateManagers() {
	override := p.Override
	resourceManager := resourcemanager.New(p.Client, p.Scheme)
	p.ClusterRoleManager = resourceManager.CreateClusterRoleManager("", override.ClusterRole, p.GetLabels, p.Config.CAInitConfig.RoleFile)
	p.ClusterRoleBindingManager = resourceManager.CreateClusterRoleBindingManager("", override.ClusterRoleBinding, p.GetLabels, p.Config.CAInitConfig.RoleBindingFile)
}

func (p *BaseProposal) PreReconcileChecks(instance *current.Proposal) (bool, error) {
	// todo add
	return false, nil
}

func (p *BaseProposal) Reconcile(instance *current.Proposal) (result common.Result, err error) {
	log.Info("Reconciling...")
	update := false
	if instance.Spec.EndAt.IsZero() {
		instance.Spec.EndAt = v1.NewTime(time.Now().Add(time.Hour * 24))
		update = true
	}
	if instance.Spec.StartAt.IsZero() {
		instance.Spec.StartAt = v1.Now()
		update = true
	}
	if update {
		if err := p.Client.Update(context.TODO(), instance); err != nil {
			return common.Result{}, err
		} else {
			return common.Result{}, nil
		}
	}
	if instance.Status.Phase == current.ProposalPending {

	} else if instance.Status.Phase == current.ProposalVoting {

	} else if instance.Status.Phase == current.ProposalFinished {
	}

	if err = p.ReconcileManagers(context.TODO(), instance); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}

	return common.Result{}, nil
}

func (c *BaseProposal) ReconcileManagers(ctx context.Context, instance *current.Proposal) (err error) {
	if err = c.ReconcileRBAC(instance); err != nil {
		return errors.Wrap(err, "failed RBAC reconciliation")
	}
	if err = c.ReconcileVote(ctx, instance); err != nil {
		return errors.Wrap(err, "failed Vote reconciliation")
	}
	return
}

func (c *BaseProposal) ReconcileRBAC(instance *current.Proposal) (err error) {
	if err = c.ClusterRoleManager.Reconcile(instance, false); err != nil {
		return errors.Wrap(err, "failed ClusterRole reconciliation")
	}

	if err = c.ClusterRoleBindingManager.Reconcile(instance, false); err != nil {
		return errors.Wrap(err, "failed ClusterRoleBinding reconciliation")
	}

	return
}

func (c *BaseProposal) ValidateSpec(instance *current.Proposal) error {
	return nil
}

func (c *BaseProposal) GetLabels(instance v1.Object) map[string]string {
	label := os.Getenv("OPERATOR_LABEL_PREFIX")
	if label == "" {
		label = "fabric"
	}

	return map[string]string{
		"app":                          instance.GetName(),
		"creator":                      label,
		"release":                      "operator",
		"helm.sh/chart":                "ibm-" + label,
		"app.kubernetes.io/name":       label,
		"app.kubernetes.io/instance":   label + "Proposal",
		"app.kubernetes.io/managed-by": label + "-operator",
	}
}

func (c *BaseProposal) ReconcileVote(ctx context.Context, instance *current.Proposal) (err error) {
	return c.CreateVoteIfNotExists(ctx, instance)
}

func (c *BaseProposal) CreateVoteIfNotExists(ctx context.Context, instance *current.Proposal) error {
	organizations, err := instance.GetCandidateOrganizations(ctx, c.Client)
	if err != nil {
		log.Error(err, "cant get organizations for proposal:"+instance.GetName())
		return err
	}
	wg := sync.WaitGroup{}
	for _, org := range organizations {
		wg.Add(1)
		go func(org current.NamespacedName) {
			defer func() {
				wg.Done()
			}()
			vote := &current.Vote{
				ObjectMeta: v1.ObjectMeta{
					Name:      instance.GetVoteName(org.Name),
					Namespace: org.Namespace,
					Labels:    instance.GetVoteLabel(),
				},
				Spec: current.VoteSpec{
					ProposalName:     instance.GetName(),
					OrganizationName: org.Name,
					Decision:         nil,
					Description:      "",
				},
			}
			if err = c.Client.Get(ctx, types.NamespacedName{Namespace: org.Namespace, Name: instance.GetVoteName(org.Name)}, vote); err != nil {
				if k8sruntime.IgnoreNotFound(err) == nil {
					log.Info(fmt.Sprintf("not find vote in org:%s, crate now.", org))
					if err = c.Client.Create(ctx, vote, k8sclient.CreateOption{Owner: instance, Scheme: c.Scheme}); err != nil {
						log.Error(err, "Error create vote")
					}
				} else {
					log.Error(err, fmt.Sprintf("Error getting vote in org:%s", org))
					// todo return error
				}
			} else {
				if org.Name == instance.Spec.Initiator.Name && org.Namespace == instance.Spec.Initiator.Namespace {
					//if err = c.Client.Get(ctx, types.NamespacedName{Namespace: org.Namespace, Name: instance.GetVoteName(org.Name)}, vote); err != nil {
					//	log.Error(err, "Error get vote")
					//}
					if pointer.BoolDeref(vote.Spec.Decision, false) {
						return
					}
					vote.Spec.Decision = pointer.Bool(true)
					if err = c.Client.Patch(ctx, vote, nil, k8sclient.PatchOption{Resilient: &k8sclient.ResilientPatch{Retry: 3, Into: &current.Vote{}, Strategy: k8sruntime.MergeFrom}}); err != nil {
						log.Error(err, "Error patch vote")
					}
				}
			}
		}(org)
	}
	wg.Wait()
	return nil
}
