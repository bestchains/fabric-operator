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

package network

import (
	"fmt"
	"reflect"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/go-test/deep"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileNetwork) CreateFunc(e event.CreateEvent) bool {
	var reconcile bool

	switch e.Object.(type) {
	case *current.Network:
		network := e.Object.(*current.Network)
		log.Info(fmt.Sprintf("Create event detected for network '%s'", network.GetNamespacedName()))
		reconcile = r.PredictNetworkCreate(network)

	}

	return reconcile
}

func (r *ReconcileNetwork) PredictNetworkCreate(network *current.Network) bool {
	update := Update{}

	if network.HasType() {
		log.Info(fmt.Sprintf("Operator restart detected, running update flow on existing network '%s'", network.GetNamespacedName()))

		// Get the spec state of the resource before the operator went down, this
		// will be used to compare to see if the spec of resources has changed
		cm, err := r.GetSpecState(network)
		if err != nil {
			log.Info(fmt.Sprintf("Failed getting saved fedeation spec '%s', triggering create: %s", network.GetNamespacedName(), err.Error()))
			return true
		}

		specBytes := cm.BinaryData["spec"]
		existingNet := &current.Network{}
		err = yaml.Unmarshal(specBytes, &existingNet.Spec)
		if err != nil {
			log.Info(fmt.Sprintf("Unmarshal failed for saved network spec '%s', triggering create: %s", network.GetNamespacedName(), err.Error()))
			return true
		}

		diff := deep.Equal(network.Spec, existingNet.Spec)
		if diff != nil {
			log.Info(fmt.Sprintf("Network '%s' spec was updated while operator was down", network.GetNamespacedName()))
			log.Info(fmt.Sprintf("Difference detected: %v", diff))
			update.specUpdated = true
		}

		added, removed := current.DifferMembers(network.Spec.Members, existingNet.Spec.Members)
		if len(added) != 0 || len(removed) != 0 {
			log.Info(fmt.Sprintf("Network '%s' members was updated while operator was down", network.GetNamespacedName()))
			log.Info(fmt.Sprintf("Difference detected: added members %v", added))
			log.Info(fmt.Sprintf("Difference detected: removed members %v", removed))
			update.memberUpdated = true
		}

		log.Info(fmt.Sprintf("Create event triggering reconcile for updating Network '%s'", network.GetNamespacedName()))
		r.PushUpdate(network.GetNamespacedName(), update)
		return true
	}

	update.specUpdated = true
	update.memberUpdated = true

	r.PushUpdate(network.GetNamespacedName(), update)

	return true
}

// Watch Network & Proposal
func (r *ReconcileNetwork) UpdateFunc(e event.UpdateEvent) bool {
	var reconcile bool

	switch e.ObjectOld.(type) {
	case *current.Network:
		oldNet := e.ObjectOld.(*current.Network)
		newNet := e.ObjectNew.(*current.Network)
		log.Info(fmt.Sprintf("Update event detected for network '%s'", oldNet.GetNamespacedName()))

		reconcile = r.PredicNetworkUpdate(oldNet, newNet)
	}
	return reconcile
}

func (r *ReconcileNetwork) PredicNetworkUpdate(oldNet *current.Network, newNet *current.Network) bool {
	update := Update{}

	if reflect.DeepEqual(oldNet.Spec, newNet.Spec) {
		return false
	}

	update.specUpdated = true

	added, removed := current.DifferMembers(oldNet.GetMembers(), newNet.GetMembers())
	if len(added) != 0 || len(removed) != 0 {
		log.Info(fmt.Sprintf("Difference detected: added members %v", added))
		log.Info(fmt.Sprintf("Difference detected: removed members %v", removed))
		update.memberUpdated = true
	}

	r.PushUpdate(oldNet.GetNamespacedName(), update)

	log.Info(fmt.Sprintf("Spec update triggering reconcile on Network custom resource %s: update [ %+v ]", oldNet.Name, update.GetUpdateStackWithTrues()))

	return true
}
