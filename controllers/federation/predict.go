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

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileFederation) CreateFunc(e event.CreateEvent) bool {
	var reconcile bool

	switch e.Object.(type) {
	case *current.Federation:
		federation := e.Object.(*current.Federation)
		log.Info(fmt.Sprintf("Create event detected for federation '%s'", federation.GetNamespacedName()))
		reconcile = r.PredictFederationCreate(federation)

	case *current.Network:
		network := e.Object.(*current.Network)
		log.Info(fmt.Sprintf("Create event detected for network '%s'", network.GetNamespacedName()))
		reconcile = r.PredictNetworkCreate(network)
	}

	return reconcile
}

func (r *ReconcileFederation) PredictFederationCreate(federation *current.Federation) bool {
	update := Update{}

	if federation.HasType() {
		log.Info(fmt.Sprintf("Operator restart detected, running update flow on existing federation '%s'", federation.GetNamespacedName()))

		// Get the spec state of the resource before the operator went down, this
		// will be used to compare to see if the spec of resources has changed
		cm, err := r.GetSpecState(federation)
		if err != nil {
			log.Info(fmt.Sprintf("Failed getting saved fedeation spec '%s', triggering create: %s", federation.GetNamespacedName(), err.Error()))
			return true
		}

		specBytes := cm.BinaryData["spec"]
		existingFed := &current.Federation{}
		err = yaml.Unmarshal(specBytes, &existingFed.Spec)
		if err != nil {
			log.Info(fmt.Sprintf("Unmarshal failed for saved federation spec '%s', triggering create: %s", federation.GetNamespacedName(), err.Error()))
			return true
		}

		diff := deep.Equal(federation.Spec, existingFed.Spec)
		if diff != nil {
			log.Info(fmt.Sprintf("Federation '%s' spec was updated while operator was down", federation.GetNamespacedName()))
			log.Info(fmt.Sprintf("Difference detected: %v", diff))
			update.specUpdated = true
		}

		added, removed := current.DifferMembers(federation.Spec.Members, existingFed.Spec.Members)
		if len(added) != 0 || len(removed) != 0 {
			log.Info(fmt.Sprintf("Federation '%s' members was updated while operator was down", federation.GetNamespacedName()))
			log.Info(fmt.Sprintf("Difference detected: added members %v", added))
			log.Info(fmt.Sprintf("Difference detected: removed members %v", removed))
			update.memberUpdated = true
		}

		log.Info(fmt.Sprintf("Create event triggering reconcile for updating Federation '%s'", federation.GetNamespacedName()))
		r.PushUpdate(federation.GetNamespacedName(), update)
		return true
	}

	update.specUpdated = true
	update.memberUpdated = true
	r.PushUpdate(federation.GetNamespacedName(), update)

	return true
}

func (r *ReconcileFederation) PredictNetworkCreate(network *current.Network) bool {
	err := r.AddNetwork(network.Spec.Federation, network.NamespacedName())
	if err != nil {
		log.Error(err, fmt.Sprintf("Network %s in Federation %s", network.GetNamespacedName(), network.Spec.Federation.String()))
	}
	return false
}

func (r *ReconcileFederation) AddNetwork(fedns current.NamespacedName, netns current.NamespacedName) error {
	var err error
	federation := &current.Federation{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      fedns.Name,
		Namespace: fedns.Namespace,
	}, federation)
	if err != nil {
		return err
	}

	conflict := federation.Status.AddNetwork(current.NamespacedName{
		Name:      netns.Name,
		Namespace: netns.Namespace,
	})
	// conflict detected,do not need to PatchStatus
	if conflict {
		return errors.Errorf("network %s already exist in federation %s", netns.String(), fedns.String())
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

// Watch Federation & Proposal
func (r *ReconcileFederation) UpdateFunc(e event.UpdateEvent) bool {
	var reconcile bool

	switch e.ObjectOld.(type) {
	case *current.Federation:
		oldFed := e.ObjectOld.(*current.Federation)
		newFed := e.ObjectNew.(*current.Federation)
		log.Info(fmt.Sprintf("Update event detected for federation '%s'", oldFed.GetNamespacedName()))
		reconcile = r.PredicFederationUpdate(oldFed, newFed)
	case *current.Proposal:
		oldProposal := e.ObjectOld.(*current.Proposal)
		newProposal := e.ObjectNew.(*current.Proposal)
		log.Info(fmt.Sprintf("Update event detected for proposal '%s'", oldProposal.Spec.Federation.String()))
		reconcile = r.PredicProposalUpdate(oldProposal, newProposal)
	}
	return reconcile
}

func (r *ReconcileFederation) PredicFederationUpdate(oldFed *current.Federation, newFed *current.Federation) bool {
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

	r.PushUpdate(oldFed.GetNamespacedName(), update)

	log.Info(fmt.Sprintf("Spec update triggering reconcile on Federation custom resource %s: update [ %+v ]", oldFed.Name, update.GetUpdateStackWithTrues()))

	return true
}

func (r *ReconcileFederation) DeleteFunc(e event.DeleteEvent) bool {
	var reconcile bool
	switch e.Object.(type) {
	case *current.Network:
		network := e.Object.(*current.Network)
		log.Info(fmt.Sprintf("Delete event detected for network '%s'", network.GetNamespacedName()))
		reconcile = r.PredictNetworkDelete(network)
	}
	return reconcile
}

func (r *ReconcileFederation) PredictNetworkDelete(network *current.Network) bool {
	fedns := network.Spec.Federation
	netns := network.NamespacedName()
	err := r.DeleteNetwork(fedns, netns)
	if err != nil {
		log.Error(err, fmt.Sprintf("Network %s in Federation %s", netns.String(), fedns.String()))
	}
	return false
}

func (r *ReconcileFederation) DeleteNetwork(fedns current.NamespacedName, netns current.NamespacedName) error {
	var err error
	federation := &current.Federation{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      fedns.Name,
		Namespace: fedns.Namespace,
	}, federation)
	if err != nil {
		return err
	}

	exist := federation.Status.DeleteNetwork(current.NamespacedName{
		Name:      netns.Name,
		Namespace: netns.Namespace,
	})

	// network do not exist in this federation ,do not need to PatchStatus
	if !exist {
		return errors.Errorf("network %s not exist in federation %s", netns.String(), fedns.String())
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

func (r *ReconcileFederation) PredicProposalUpdate(oldProposal *current.Proposal, newProposal *current.Proposal) bool {
	update := Update{}

	if reflect.DeepEqual(oldProposal.Spec, newProposal.Spec) && reflect.DeepEqual(oldProposal.Status, newProposal.Status) {
		return false
	}

	if newProposal.Status.Phase == current.ProposalFinished {
		for _, c := range newProposal.Status.Conditions {
			if c.Status == metav1.ConditionTrue {
				if newProposal.Spec.CreateFederation != nil {
					switch c.Type {
					case current.ProposalFailed:
						update.proposalFailed = true
					case current.ProposalSucceeded:
						update.proposalActivated = true
					}
				} else if newProposal.Spec.AddMember != nil {
					if c.Type == current.ProposalSucceeded {
						fed := &current.Federation{}
						if err := r.client.Get(context.TODO(), types.NamespacedName{Name: newProposal.Spec.Federation.Name, Namespace: newProposal.Spec.Federation.Namespace}, fed); err != nil {
							log.Error(err, fmt.Sprintf("cant find federation %s", newProposal.Spec.Federation.String()))
							return false
						}
						for _, m := range newProposal.Spec.AddMember.Members {
							fed.Spec.Members = append(fed.Spec.Members, current.Member{
								NamespacedName: m,
								Initiator:      false,
							})
						}
						if err := r.client.Update(context.TODO(), fed); err != nil {
							log.Error(err, fmt.Sprintf("cant update federation %s", newProposal.Spec.Federation.String()))
							return false
						}
					}
				} else if newProposal.Spec.DeleteMember != nil {
					if c.Type == current.ProposalSucceeded {
						fed := &current.Federation{}
						if err := r.client.Get(context.TODO(), types.NamespacedName{Name: newProposal.Spec.Federation.Name, Namespace: newProposal.Spec.Federation.Namespace}, fed); err != nil {
							log.Error(err, fmt.Sprintf("cant find federation %s", newProposal.Spec.Federation.String()))
							return false
						}
						newMember := make([]current.Member, 0)
						for _, m := range fed.Spec.Members {
							if m.String() == newProposal.Spec.DeleteMember.Member.String() {
								continue
							}
							newMember = append(newMember, m)
						}
						fed.Spec.Members = newMember
						if err := r.client.Update(context.TODO(), fed); err != nil {
							log.Error(err, fmt.Sprintf("cant update federation %s", newProposal.Spec.Federation.String()))
							return false
						}
					}
				} else if newProposal.Spec.DissolveFederation != nil {
					if c.Type == current.ProposalSucceeded {
						update.proposalDissolved = true
						fed := &current.Federation{}
						if err := r.client.Get(context.TODO(), types.NamespacedName{Name: newProposal.Spec.Federation.Name, Namespace: newProposal.Spec.Federation.Namespace}, fed); err != nil {
							log.Error(err, fmt.Sprintf("cant find federation %s", newProposal.Spec.Federation.String()))
							return false
						}
						newMember := make([]current.Member, 0)
						for _, m := range fed.Spec.Members {
							if !m.Initiator {
								continue
							}
							newMember = append(newMember, m)
						}
						fed.Spec.Members = newMember
						if err := r.client.Update(context.TODO(), fed); err != nil {
							log.Error(err, fmt.Sprintf("cant update federation %s", newProposal.Spec.Federation.String()))
							return false
						}
					}
				}
			}
		}
	}
	if !(update.proposalDissolved || update.proposalActivated || update.proposalFailed) {
		return false
	}
	r.PushUpdate(newProposal.Spec.Federation.String(), update)
	log.Info(fmt.Sprintf("Proposal Status update triggering reconcile on Federation custom resource %s: update [ %+v ]", newProposal.Spec.Federation.String(), update.GetUpdateStackWithTrues()))
	return true
}
