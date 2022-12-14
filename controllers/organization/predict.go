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
	"time"

	iam "github.com/IBM-Blockchain/fabric-operator/api/iam/v1alpha1"
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileOrganization) CreateFunc(e event.CreateEvent) bool {
	var reconcile bool
	switch e.Object.(type) {
	case *current.Organization:
		organization := e.Object.(*current.Organization)
		log.Info(fmt.Sprintf("Create event detected for organization '%s'", organization.GetName()))
		reconcile = r.PredictOrganizationCreate(organization)
		if reconcile {
			log.Info(fmt.Sprintf("Create event triggering reconcile for creating organization '%s'", organization.GetName()))
		}

	case *current.Federation:
		federation := e.Object.(*current.Federation)
		log.Info(fmt.Sprintf("Create event detected for federation '%s'", federation.GetName()))
		reconcile = r.PredictFederationCreate(federation)
	case *current.IBPCA:
		reconcile = false
	}
	return reconcile
}

func (r *ReconcileOrganization) PredictOrganizationCreate(organization *current.Organization) bool {
	update := Update{}
	if organization.HasType() {
		log.Info(fmt.Sprintf("Operator restart detected, running update flow on existing organization '%s'", organization.GetName()))

		cm, err := r.GetSpecState(organization)
		if err != nil {
			log.Info(fmt.Sprintf("Failed getting saved organization spec '%s', triggering create: %s", organization.GetName(), err.Error()))
			return true
		}

		specBytes := cm.BinaryData["spec"]
		existingOrg := &current.Organization{}
		err = yaml.Unmarshal(specBytes, &existingOrg.Spec)
		if err != nil {
			log.Info(fmt.Sprintf("Unmarshal failed for saved organization spec '%s', triggering create: %s", organization.GetName(), err.Error()))
			return true
		}

		diff := deep.Equal(organization.Spec, existingOrg.Spec)
		if diff != nil {
			log.Info(fmt.Sprintf("Organization '%s' spec was updated while operator was down", organization.GetName()))
			log.Info(fmt.Sprintf("Difference detected: %s", diff))
			update.specUpdated = true
		}
		if organization.Spec.Admin != existingOrg.Spec.Admin {
			update.adminUpdated = true
		}
		r.PushUpdate(organization.GetName(), update)
		return true
	}

	update.adminUpdated = true
	r.PushUpdate(organization.GetName(), update)
	return true
}

func (r *ReconcileOrganization) PredictFederationCreate(federation *current.Federation) bool {
	var err error

	for _, m := range federation.Spec.Members {
		err = r.AddFed(m, federation)
		if err != nil {
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", m.GetNamespacedName(), federation.GetName()))
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
		log.Info(fmt.Sprintf("Update event detected for organization '%s'", oldOrg.GetName()))

		reconcile = r.PredictOrganizationUpdate(oldOrg, newOrg)

	case *current.Federation:
		oldFed := e.ObjectOld.(*current.Federation)
		newFed := e.ObjectNew.(*current.Federation)
		log.Info(fmt.Sprintf("Update event detected for fedeartion '%s'", oldFed.GetName()))

		reconcile = r.PredictFederationUpdate(oldFed, newFed)
	case *current.IBPCA:
		oldCA := e.ObjectOld.(*current.IBPCA)
		newCA := e.ObjectNew.(*current.IBPCA)
		log.Info(fmt.Sprintf("Update event detected for ibpca '%s'", oldCA.GetName()))

		reconcile = r.PredictCAUpdate(oldCA, newCA)
	}
	return reconcile
}

func (r *ReconcileOrganization) PredictOrganizationUpdate(oldOrg *current.Organization, newOrg *current.Organization) bool {
	update := Update{}

	if reflect.DeepEqual(oldOrg.Spec, newOrg.Spec) {
		return false
	}

	if oldOrg.Spec.Admin != newOrg.Spec.Admin {
		update.adminUpdated = true
		if r.Config.OrganizationInitConfig.IAMEnabled {
			// delete annotations from previous Admin
			oldAdminUser, err := r.GetIAMUser(oldOrg.Spec.Admin)
			if err != nil {
				log.Error(err, fmt.Sprintf("failed to get iam user %s", oldOrg.Spec.Admin))
			} else {
				err = r.DeleteBlockchainAnnotations(oldOrg, oldAdminUser)
				if err != nil {
					log.Error(err, fmt.Sprintf("failed to delete annotation %s for %s", oldOrg.GetAnnotationKey(), oldOrg.Spec.Admin))
				}
			}
		}
	}

	r.PushUpdate(oldOrg.Name, update)

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
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", am.GetNamespacedName(), newFed.GetName()))
		}
	}

	for _, rm := range removed {
		err = r.DeleteFed(rm, newFed)
		if err != nil {
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", rm.GetNamespacedName(), newFed.GetName()))
		}
	}

	return false
}

func (r *ReconcileOrganization) PredictCAUpdate(oldCA *current.IBPCA, newCA *current.IBPCA) bool {
	if newCA.Status.CRStatus.Type == current.Deployed {
		organization := newCA.GetOrganization()
		err := r.SetStatusToDeployed(organization)
		if err != nil {
			log.Error(err, fmt.Sprintf("set organization %s to `Deployed`", organization.Name))
		}
	}
	return false
}
func (r *ReconcileOrganization) DeleteFunc(e event.DeleteEvent) bool {
	var reconcile bool
	switch e.Object.(type) {
	case *current.Organization:
		organiation := e.Object.(*current.Organization)
		log.Info(fmt.Sprintf("Delete event detected for organization '%s'", organiation.GetName()))
		reconcile = r.PredictOrganizationDelete(organiation)
	case *current.Federation:
		federation := e.Object.(*current.Federation)
		log.Info(fmt.Sprintf("Delete event detected for federation '%s'", federation.GetName()))
		reconcile = r.PredictFederationDelete(federation)
	}
	return reconcile
}

func (r *ReconcileOrganization) PredictOrganizationDelete(organization *current.Organization) bool {
	var err error
	if r.Config.OrganizationInitConfig.IAMEnabled {
		userList, err := r.GetIAMUsers(organization.GetName())
		if err != nil {
			log.Error(err, fmt.Sprintf("failed to get iam users with organization annotation key %s", organization.GetAnnotationKey()))
		} else {
			for _, iamuser := range userList.Items {
				err = r.DeleteBlockchainAnnotations(organization, &iamuser)
				if err != nil {
					log.Error(err, fmt.Sprintf("failed to delete annotation %s for %s", organization.GetAnnotationKey(), iamuser.GetName()))
				}
			}
		}
	}
	// delete namespace
	ns := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: organization.GetUserNamespace(),
		},
	}
	err = r.client.Delete(context.TODO(), ns, &client.DeleteOptions{})
	if err != nil {
		log.Error(err, fmt.Sprintf("failed to delete namespace for organiation %s", organization.GetName()))
	}
	return false
}

func (r *ReconcileOrganization) PredictFederationDelete(federation *current.Federation) bool {
	var err error

	for _, m := range federation.Spec.Members {
		err = r.DeleteFed(m, federation)
		if err != nil {
			log.Error(err, fmt.Sprintf("Member %s in Federation %s", m.GetNamespacedName(), federation.GetName()))
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
		return errors.Errorf("federation %s already exist in organization %s", federation.GetName(), m.GetNamespacedName())
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

func (r *ReconcileOrganization) SetStatusToDeployed(organization current.NamespacedName) error {
	var err error
	org := &current.Organization{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: organization.Name}, org)
	if err != nil {
		return err
	}

	status := org.Status.CRStatus
	status.Type = current.Deployed
	status.Status = current.True
	status.Reason = "IBPCA Deployed"
	status.LastHeartbeatTime = time.Now().String()

	org.Status = current.OrganizationStatus{
		CRStatus: status,
	}

	err = r.client.PatchStatus(context.TODO(), org, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    2,
			Into:     &current.Organization{},
			Strategy: client.MergeFrom,
		},
	})

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
		return errors.Errorf("federation %s not exist in organization %s", federation.GetName(), m.GetNamespacedName())
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

func (r *ReconcileOrganization) GetIAMUser(username string) (*iam.User, error) {
	var err error
	iamuser := &iam.User{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: username}, iamuser)
	if err != nil {
		return nil, err
	}
	return iamuser, nil
}

func (r *ReconcileOrganization) GetIAMUsers(annotationKey string) (*iam.UserList, error) {
	userList := &iam.UserList{}
	err := r.client.List(context.TODO(), userList, &client.ListOptions{})
	if err != nil {
		return nil, err
	}
	return userList, nil
}

func (r *ReconcileOrganization) DeleteBlockchainAnnotations(instance *current.Organization, iamuser *iam.User) error {
	var err error

	annotationList := &current.BlockchainAnnotationList{
		List: make(map[string]current.BlockchainAnnotation),
	}

	err = annotationList.Unmarshal([]byte(iamuser.Annotations[current.BlockchainAnnotationKey]))
	if err != nil {
		return err
	}

	err = annotationList.DeleteAnnotation(instance.GetAnnotationKey())
	if err != nil {
		return err
	}

	raw, err := annotationList.Marshal()
	if err != nil {
		return err
	}

	iamuser.Annotations[current.BlockchainAnnotationKey] = string(raw)

	err = r.client.Patch(context.TODO(), iamuser, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    2,
			Into:     &iam.User{},
			Strategy: client.MergeFrom,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
