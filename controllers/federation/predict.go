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
	"reflect"
	"time"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	bcrbac "github.com/IBM-Blockchain/fabric-operator/pkg/rbac"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileFederation) CreateFunc(e event.CreateEvent) bool {
	federation := e.Object.(*current.Federation)
	log.Info(fmt.Sprintf("Create event detected for federation '%s'", federation.GetName()))

	update := Update{}

	if federation.HasType() {
		log.Info(fmt.Sprintf("Operator restart detected, running update flow on existing federation '%s'", federation.GetName()))

		// Get the spec state of the resource before the operator went down, this
		// will be used to compare to see if the spec of resources has changed
		cm, err := r.GetSpecState(federation)
		if err != nil {
			log.Info(fmt.Sprintf("Failed getting saved fedeation spec '%s', triggering create: %s", federation.GetName(), err.Error()))
			return true
		}

		specBytes := cm.BinaryData["spec"]
		existingFed := &current.Federation{}
		err = yaml.Unmarshal(specBytes, &existingFed.Spec)
		if err != nil {
			log.Info(fmt.Sprintf("Unmarshal failed for saved federation spec '%s', triggering create: %s", federation.GetName(), err.Error()))
			return true
		}

		diff := deep.Equal(federation.Spec, existingFed.Spec)
		if diff != nil {
			log.Info(fmt.Sprintf("Federation '%s' spec was updated while operator was down", federation.GetName()))
			log.Info(fmt.Sprintf("Difference detected: %v", diff))
			update.specUpdated = true
		}

		added, removed := current.DifferMembers(federation.Spec.Members, existingFed.Spec.Members)
		if len(added) != 0 || len(removed) != 0 {
			log.Info(fmt.Sprintf("Federation '%s' members was updated while operator was down", federation.GetName()))
			log.Info(fmt.Sprintf("Difference detected: added members %v", added))
			log.Info(fmt.Sprintf("Difference detected: removed members %v", removed))
			update.memberUpdated = true
		}

		log.Info(fmt.Sprintf("Create event triggering reconcile for updating Federation '%s'", federation.GetName()))
		r.PushUpdate(federation.GetName(), update)
		return true
	}

	update.specUpdated = true
	update.memberUpdated = true
	r.PushUpdate(federation.GetName(), update)

	return true
}

func (r *ReconcileFederation) NetworkCreateFunc(e event.CreateEvent) bool {
	network := e.Object.(*current.Network)
	log.Info(fmt.Sprintf("Create event detected for network '%s'", network.GetName()))

	err := r.AddNetwork(network.Spec.Federation, network.Name)
	if err != nil {
		log.Error(err, fmt.Sprintf("Network %s in Federation %s", network.GetName(), network.Spec.Federation))
	}
	return false
}

func (r *ReconcileFederation) AddNetwork(fedns, netns string) error {
	var err error
	federation := &current.Federation{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name: fedns,
	}, federation)
	if err != nil {
		return err
	}

	conflict := federation.Status.AddNetwork(netns)
	// conflict detected,do not need to PatchStatus
	if conflict {
		return errors.Errorf("network %s already exist in federation %s", netns, fedns)
	}

	err = r.client.PatchStatus(context.TODO(), federation, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    2,
			Into:     &current.Federation{},
			Strategy: client.MergeFrom,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileFederation) UpdateFunc(e event.UpdateEvent) bool {
	oldFed := e.ObjectOld.(*current.Federation)
	newFed := e.ObjectNew.(*current.Federation)
	log.Info(fmt.Sprintf("Update event detected for federation '%s'", oldFed.GetName()))

	update := Update{}

	if reflect.DeepEqual(oldFed.Spec, newFed.Spec) {
		return false
	}

	update.specUpdated = true

	added, removed := current.DifferMembers(oldFed.GetMembers(), newFed.GetMembers())
	if len(added) != 0 || len(removed) != 0 {
		log.Info(fmt.Sprintf("Difference detected: added members %v", added))
		log.Info(fmt.Sprintf("Difference detected: removed members %v", removed))
		update.memberUpdated = true
	}

	r.PushUpdate(oldFed.GetName(), update)

	log.Info(fmt.Sprintf("Spec update triggering reconcile on Federation custom resource %s: update [ %+v ]", oldFed.Name, update.GetUpdateStackWithTrues()))

	return true
}

func (r *ReconcileFederation) ProposalUpdateFunc(e event.UpdateEvent) bool {
	oldProposal := e.ObjectOld.(*current.Proposal)
	newProposal := e.ObjectNew.(*current.Proposal)
	log.Info(fmt.Sprintf("Update event detected for proposal '%s'", oldProposal.GetName()))
	update := Update{}

	if reflect.DeepEqual(oldProposal.Spec, newProposal.Spec) && reflect.DeepEqual(oldProposal.Status, newProposal.Status) {
		return false
	}

	fed := &current.Federation{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: newProposal.Spec.Federation}, fed); err != nil {
		log.Error(err, fmt.Sprintf("cant find federation %s", newProposal.Spec.Federation))
		return false
	}
	newMember := make([]current.Member, 0)
	now := v1.Now()
	if !newProposal.Spec.EndAt.IsZero() && newProposal.Spec.EndAt.Before(&now) {
		fed.Status.CRStatus = current.CRStatus{
			Type:              current.FederationFailed,
			Status:            current.False,
			Reason:            "vote expired",
			Message:           fmt.Sprintf("vote end date is %s, current is %s", newProposal.Spec.EndAt.Format(time.RFC3339), now.Format(time.RFC3339)),
			LastHeartbeatTime: metav1.Now(),
		}
		if err := r.client.PatchStatus(context.TODO(), fed, nil, k8sclient.PatchOption{
			Resilient: &k8sclient.ResilientPatch{
				Retry:    3,
				Into:     &current.Federation{},
				Strategy: client.MergeFrom,
			},
		}); err != nil {
			log.Error(err, "vote expired, failed to update federation status")
		}
		return false
	}

	if newProposal.Status.Phase == current.ProposalFinished {
		for _, c := range newProposal.Status.Conditions {
			switch c.Type {
			case current.ProposalSucceeded:
				switch newProposal.GetPurpose() {
				case current.CreateFederationProposal:
					update.proposalActivated = true
					for _, m := range fed.Spec.Members {
						if m.JoinedBy == "" {
							m.JoinedBy = newProposal.GetName()
							m.JoinedAt = &now
						}
						newMember = append(newMember, m)
					}
				case current.AddMemberProposal:
					newMember = append(newMember, fed.Spec.Members...)
					hasAddMemberName := make(map[string]bool, len(newMember)+len(newProposal.Spec.AddMember.Members))
					for _, m := range fed.Spec.Members {
						hasAddMemberName[m.Name] = true
					}
					for _, m := range newProposal.Spec.AddMember.Members {
						if hasAddMemberName[m] {
							continue
						}
						newMember = append(newMember, current.Member{
							Name:      m,
							Initiator: false,
							JoinedBy:  newProposal.GetName(),
							JoinedAt:  &now,
						})
						hasAddMemberName[m] = true
					}
				case current.DeleteMemberProposal:
					for _, m := range fed.Spec.Members {
						if m.Name == newProposal.Spec.DeleteMember.Member {
							continue
						}
						newMember = append(newMember, m)
					}
				case current.DissolveFederationProposal:
					update.proposalDissolved = true
					for _, m := range fed.Spec.Members {
						if !m.Initiator {
							continue
						}
						newMember = append(newMember, m)
					}
				case current.DissolveNetworkProposal:
					network := &current.Network{}
					if err := r.client.Get(context.TODO(), types.NamespacedName{Name: newProposal.Spec.DissolveNetwork.Name}, network); err != nil {
						if apierrors.IsNotFound(err) {
							return false
						}
						log.Error(err, fmt.Sprintf("cant find network %s", newProposal.Spec.DissolveNetwork.Name))
						return false
					}
					if network.Status.CRStatus.Type != current.NetworkDissoleved {
						network.Status.CRStatus.Type = current.NetworkDissoleved
						network.Status.Status = current.True
						network.Status.Reason = "DissolveNetworkProposal"
						network.Status.Message = "Proposal to dissolve this network succeeded"
						network.Status.LastHeartbeatTime = metav1.Now()
						if err := r.client.PatchStatus(context.TODO(), network, nil, k8sclient.PatchOption{
							Resilient: &k8sclient.ResilientPatch{
								Retry:    2,
								Into:     &current.Network{},
								Strategy: client.MergeFrom,
							},
						}); err != nil {
							log.Error(err, fmt.Sprintf("cant patch network %s status to %s", newProposal.Spec.DissolveNetwork.Name, current.NetworkDissoleved))
							return false
						}
					}
				}
			case current.ProposalFailed:
				switch newProposal.GetPurpose() {
				case current.CreateFederationProposal:
					update.proposalFailed = true
				}
			}
		}
	}
	// Update if member changes
	if len(newMember) != 0 {
		fed.Spec.Members = newMember
		if err := r.client.Patch(context.TODO(), fed, nil, k8sclient.PatchOption{
			Resilient: &k8sclient.ResilientPatch{
				Retry:    3,
				Into:     &current.Federation{},
				Strategy: client.MergeFrom,
			},
		}); err != nil {
			log.Error(err, fmt.Sprintf("cant update federation %s", newProposal.Spec.Federation))
			return false
		}
	}
	if !(update.proposalDissolved || update.proposalActivated || update.proposalFailed) {
		return false
	}
	r.PushUpdate(newProposal.Spec.Federation, update)
	log.Info(fmt.Sprintf("Proposal Status update triggering reconcile on Federation custom resource %s: update [ %+v ]", newProposal.Spec.Federation, update.GetUpdateStackWithTrues()))
	return true
}

func (r *ReconcileFederation) DeleteFunc(e event.DeleteEvent) bool {
	federation := e.Object.(*current.Federation)
	log.Info(fmt.Sprintf("Delete event detected for federation '%s'", federation.GetName()))
	err := r.rbacManager.Reconcile(bcrbac.Federation, federation, bcrbac.ResourceDelete)
	if err != nil {
		log.Error(err, fmt.Sprintf("Federation %s", federation.Name))
	}
	return false
}

func (r *ReconcileFederation) NetworkDeleteFunc(e event.DeleteEvent) bool {
	network := e.Object.(*current.Network)
	fedns := network.Spec.Federation
	netns := network.Name
	err := r.DeleteNetwork(fedns, netns)
	if err != nil {
		log.Error(err, fmt.Sprintf("Network %s in Federation %s", netns, fedns))
	}
	return false
}

func (r *ReconcileFederation) DeleteNetwork(fedns, netns string) error {
	var err error
	federation := &current.Federation{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name: fedns,
	}, federation)
	if err != nil {
		return err
	}

	exist := federation.Status.DeleteNetwork(netns)

	// network do not exist in this federation ,do not need to PatchStatus
	if !exist {
		return errors.Errorf("network %s not exist in federation %s", netns, fedns)
	}

	err = r.client.PatchStatus(context.TODO(), federation, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    2,
			Into:     &current.Federation{},
			Strategy: client.MergeFrom,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
