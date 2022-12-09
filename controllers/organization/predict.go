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
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/go-test/deep"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileOrganization) CreateFunc(e event.CreateEvent) bool {
	var reconcile bool
	switch e.Object.(type) {
	case *current.Organization:
		organization := e.Object.(*current.Organization)
		log.Info(fmt.Sprintf("Create event detected for organization '%s'", organization.GetNamespacedName()))
		reconcile = r.PredictOrganizationCreate(organization)
		if reconcile {
			log.Info(fmt.Sprintf("Create event triggering reconcile for creating organization '%s'", organization.GetNamespacedName()))
		}

	case *current.Federation:
		federation := e.Object.(*current.Federation)
		log.Info(fmt.Sprintf("Create event detected for federation '%s'", federation.GetNamespacedName()))
		reconcile = r.PredictFederationCreate(federation)
	}
	return reconcile
}

func (r *ReconcileOrganization) PredictOrganizationCreate(organization *current.Organization) bool {
	update := Update{}
	if organization.HasType() {
		log.Info(fmt.Sprintf("Operator restart detected, running update flow on existing organization '%s'", organization.GetNamespacedName()))

		cm, err := r.GetSpecState(organization)
		if err != nil {
			log.Info(fmt.Sprintf("Failed getting saved organization spec '%s', triggering create: %s", organization.GetNamespacedName(), err.Error()))
			return true
		}

		specBytes := cm.BinaryData["spec"]
		existingOrg := &current.Organization{}
		err = yaml.Unmarshal(specBytes, &existingOrg.Spec)
		if err != nil {
			log.Info(fmt.Sprintf("Unmarshal failed for saved organization spec '%s', triggering create: %s", organization.GetNamespacedName(), err.Error()))
			return true
		}

		diff := deep.Equal(organization.Spec, existingOrg.Spec)
		if diff != nil {
			log.Info(fmt.Sprintf("Organization '%s' spec was updated while operator was down", organization.GetNamespacedName()))
			log.Info(fmt.Sprintf("Difference detected: %s", diff))
			update.specUpdated = true
		}
		if organization.Spec.Admin != existingOrg.Spec.Admin || organization.Spec.CAReference.Name != existingOrg.Spec.CAReference.Name {
			update.adminOrCAUpdated = true
		}
		r.PushUpdate(organization.GetNamespacedName(), update)
		return true
	}

	update.adminOrCAUpdated = true
	r.PushUpdate(organization.GetNamespacedName(), update)
	return true
}

func (r *ReconcileOrganization) PredictFederationCreate(federation *current.Federation) bool {
	var err error

	for _, m := range federation.Spec.Members {
		err = r.AddFed(m, federation)
		if err != nil {
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", m.GetNamespacedName(), federation.GetNamespacedName()))
		}
	}

	return false
}

func (r *ReconcileOrganization) UpdateFunc(e event.UpdateEvent) bool {
	var reconcile bool

	switch e.ObjectOld.(type) {
	case *current.Organization:
		oldOrg := e.ObjectOld.(*current.Organization)
		newOrg := e.ObjectNew.(*current.Organization)
		log.Info(fmt.Sprintf("Update event detected for organization '%s'", oldOrg.GetNamespacedName()))

		reconcile = r.PredictOrganizationUpdate(oldOrg, newOrg)

	case *current.Federation:
		oldFed := e.ObjectOld.(*current.Federation)
		newFed := e.ObjectNew.(*current.Federation)
		log.Info(fmt.Sprintf("Update event detected for fedeartion '%s'", oldFed.GetNamespacedName()))

		reconcile = r.PredictFederationUpdate(oldFed, newFed)

	}
	return reconcile
}

func (r *ReconcileOrganization) PredictOrganizationUpdate(oldOrg *current.Organization, newOrg *current.Organization) bool {
	update := Update{}

	if reflect.DeepEqual(oldOrg.Spec, newOrg.Spec) {
		return false
	}

	if oldOrg.Spec.Admin != newOrg.Spec.Admin || oldOrg.Spec.CAReference.Name != newOrg.Spec.CAReference.Name {
		update.adminOrCAUpdated = true
	}

	r.PushUpdate(oldOrg.GetNamespacedName(), update)

	log.Info(fmt.Sprintf("Spec update triggering reconcile on Organization custom resource %s: update [ %+v ]", oldOrg.Name, update.GetUpdateStackWithTrues()))

	return true
}

func (r *ReconcileOrganization) PredictFederationUpdate(oldFed *current.Federation, newFed *current.Federation) bool {
	var err error

	oldMembers := oldFed.Spec.Members
	newMembers := newFed.Spec.Members

	added, removed := current.DifferMembers(oldMembers, newMembers)

	for _, am := range added {
		err = r.AddFed(am, newFed)
		if err != nil {
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", am.GetNamespacedName(), newFed.GetNamespacedName()))
		}
	}

	for _, rm := range removed {
		err = r.DeleteFed(rm, newFed)
		if err != nil {
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", rm.GetNamespacedName(), newFed.GetNamespacedName()))
		}
	}

	return false
}

func (r *ReconcileOrganization) DeleteFunc(e event.DeleteEvent) bool {
	var reconcile bool
	switch e.Object.(type) {
	case *current.Federation:
		federation := e.Object.(*current.Federation)
		log.Info(fmt.Sprintf("Delete event detected for federation '%s'", federation.GetNamespacedName()))
		reconcile = r.PredictFederationDelete(federation)
	}
	return reconcile
}

func (r *ReconcileOrganization) PredictFederationDelete(federation *current.Federation) bool {
	var err error

	for _, m := range federation.Spec.Members {
		err = r.DeleteFed(m, federation)
		if err != nil {
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", m.GetNamespacedName(), federation.GetNamespacedName()))
		}
	}

	return false
}

func (r *ReconcileOrganization) AddFed(m current.Member, federation *current.Federation) error {
	var err error
	organization := &current.Organization{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      m.Name,
		Namespace: m.Namespace,
	}, organization)
	if err != nil {
		return err
	}

	conflict := organization.Status.AddFederation(current.NamespacedName{
		Name:      federation.Name,
		Namespace: federation.Namespace,
	})
	// conflict detected,do not need to PatchStatus
	if conflict {
		return errors.Errorf("federation %s already exist in organization %s", federation.GetNamespacedName(), m.GetNamespacedName())
	}

	err = r.client.PatchStatus(context.TODO(), organization, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    2,
			Into:     &current.Organization{},
			Strategy: client.MergeFrom,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileOrganization) DeleteFed(m current.Member, federation *current.Federation) error {
	var err error

	organization := &current.Organization{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      m.Name,
		Namespace: m.Namespace,
	}, organization)
	if err != nil {
		return err
	}

	exist := organization.Status.DeleteFederation(current.NamespacedName{
		Name:      federation.Name,
		Namespace: federation.Namespace,
	})

	// federation do not exist in this organization ,do not need to PatchStatus
	if !exist {
		return errors.Errorf("federation %s not exist in organization %s", federation.GetNamespacedName(), m.GetNamespacedName())
	}

	err = r.client.PatchStatus(context.TODO(), organization, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    2,
			Into:     &current.Organization{},
			Strategy: client.MergeFrom,
		},
	})
	if err != nil {
		return err
	}
	return nil
}