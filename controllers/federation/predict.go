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
	"reflect"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/go-test/deep"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileFederation) CreateFunc(e event.CreateEvent) bool {
	var reconcile bool

	switch e.Object.(type) {
	case *current.Federation:
		federation := e.Object.(*current.Federation)
		log.Info(fmt.Sprintf("Create event detected for federation '%s'", federation.GetNamespacedName()))
		reconcile = r.PredictFederationCreate(federation)

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

// Watch Federation & Proposal
func (r *ReconcileFederation) UpdateFunc(e event.UpdateEvent) bool {
	var reconcile bool

	switch e.ObjectOld.(type) {
	case *current.Federation:
		oldFed := e.ObjectOld.(*current.Federation)
		newFed := e.ObjectNew.(*current.Federation)
		log.Info(fmt.Sprintf("Update event detected for federation '%s'", oldFed.GetNamespacedName()))

		reconcile = r.PredicFederationUpdate(oldFed, newFed)

		// TODO: watch proposal status
		// case *current.Proposal:

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
