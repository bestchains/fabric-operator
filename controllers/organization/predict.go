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
	"reflect"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileOrganization) CreateFunc(e event.CreateEvent) bool {
	update := Update{}

	switch e.Object.(type) {
	case *current.Organization:
		organization := e.Object.(*current.Organization)
		log.Info(fmt.Sprintf("Create event detected for organization '%s'", organization.GetName()))

		update.adminOrCAUpdated = true
		r.PushUpdate(organization.GetName(), update)

		log.Info(fmt.Sprintf("Create event triggering reconcile for creating organization '%s'", organization.GetName()))
	case *corev1.Secret:
		// TODO: add owner reference to admin-secret and organization-secret
		return false
	}
	return true
}

func (r *ReconcileOrganization) UpdateFunc(e event.UpdateEvent) bool {
	update := Update{}

	switch e.ObjectOld.(type) {
	case *current.Organization:
		oldOrg := e.ObjectOld.(*current.Organization)
		newOrg := e.ObjectNew.(*current.Organization)
		log.Info(fmt.Sprintf("Update event detected for organization '%s'", oldOrg.GetName()))

		if reflect.DeepEqual(oldOrg.Spec, newOrg.Spec) {
			return false
		}
		if oldOrg.Spec.Admin != newOrg.Spec.Admin || oldOrg.Spec.CAReference.Name != newOrg.Spec.CAReference.Name {
			update.adminOrCAUpdated = true
		}
		r.PushUpdate(oldOrg.GetName(), update)

		log.Info(fmt.Sprintf("Spec update triggering reconcile on Organization custom resource %s: update [ %+v ]", oldOrg.Name, update.GetUpdateStackWithTrues()))
	case *corev1.Secret:

	case *corev1.ConfigMap:
		return false
	}
	return true
}

// GetUpdateStatus with index 0
func (r *ReconcileOrganization) GetUpdateStatus(instance *current.Organization) *Update {
	return r.GetUpdateStatusAtElement(instance, 0)
}

func (r *ReconcileOrganization) GetUpdateStatusAtElement(instance *current.Organization, index int) *Update {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	update := Update{}
	_, ok := r.update[instance.GetName()]
	if !ok {
		return &update
	}

	if len(r.update[instance.GetName()]) >= 1 {
		update = r.update[instance.GetName()][index]
	}

	return &update
}

func (r *ReconcileOrganization) PushUpdate(instance string, update Update) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.update[instance] = AppendUpdateIfMissing(r.update[instance], update)
}

func (r *ReconcileOrganization) PopUpdate(instance string) *Update {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	update := Update{}
	if len(r.update[instance]) >= 1 {
		update = r.update[instance][0]
		if len(r.update[instance]) == 1 {
			r.update[instance] = []Update{}
		} else {
			r.update[instance] = r.update[instance][1:]
		}
	}

	return &update
}

func AppendUpdateIfMissing(updates []Update, update Update) []Update {
	for _, u := range updates {
		if u == update {
			return updates
		}
	}
	return append(updates, update)
}

func GetUpdateStack(allUpdates map[string][]Update) string {
	stack := ""

	for orderer, updates := range allUpdates {
		currentStack := ""
		for index, update := range updates {
			currentStack += fmt.Sprintf("{ %s}", update.GetUpdateStackWithTrues())
			if index != len(updates)-1 {
				currentStack += " , "
			}
		}
		stack += fmt.Sprintf("%s: [ %s ] ", orderer, currentStack)
	}

	return stack
}
