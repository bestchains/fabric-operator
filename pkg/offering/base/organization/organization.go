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

	iam "github.com/IBM-Blockchain/fabric-operator/api/iam/v1alpha1"
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	"github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/manager"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/IBM-Blockchain/fabric-operator/version"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("base_organization")

//go:generate counterfeiter -o mocks/update.go -fake-name Update . Update

type Update interface {
	SpecUpdated() bool
	AdminUpdated() bool
}

type Override interface {
	AdminRole(v1.Object, *rbacv1.Role, resources.Action) error
	AdminRoleBinding(v1.Object, *rbacv1.RoleBinding, resources.Action) error
	AdminClusterRoleBinding(v1.Object, *rbacv1.ClusterRoleBinding, resources.Action) error

	ClientRole(v1.Object, *rbacv1.Role, resources.Action) error
}

//go:generate counterfeiter -o mocks/initializer.go -fake-name InitializerOrganization . InitializerOrganization

type InitializerOrganization interface{}

type Organization interface {
	PreReconcileChecks(instance *current.Organization, update Update) error
	Initialize(instance *current.Organization, update Update) error
	ReconcileManagers(instance *current.Organization, update Update) error
	CheckStates(instance *current.Organization) (common.Result, error)
	Reconcile(instance *current.Organization, update Update) (common.Result, error)
}

var _ Organization = (*BaseOrganization)(nil)

const (
	KIND = "ORGANIZATION"
)

type BaseOrganization struct {
	Client controllerclient.Client
	Scheme *runtime.Scheme

	Config *config.Config

	Initializer InitializerOrganization

	Override Override

	AdminRoleManager        resources.Manager
	AdminRoleBindingManager resources.Manager

	AdminClusterRoleBindingManager resources.Manager

	ClientRoleManager resources.Manager
}

func New(client controllerclient.Client, scheme *runtime.Scheme, config *config.Config) *BaseOrganization {
	o := &override.Override{
		Client: client,
	}
	base := &BaseOrganization{
		Client:   client,
		Scheme:   scheme,
		Config:   config,
		Override: o,
	}

	base.Initializer = NewInitializer(config.OrganizationInitConfig, scheme, client, base.GetLabels)

	base.CreateManagers()

	return base
}

func (organization *BaseOrganization) CreateManagers() {
	override := organization.Override
	mgr := manager.New(organization.Client, organization.Scheme)

	organization.AdminRoleManager = mgr.CreateRoleManager("", override.AdminRole, organization.GetLabels, organization.Config.OrganizationInitConfig.AdminRoleFile)
	organization.AdminRoleBindingManager = mgr.CreateRoleBindingManager("", override.AdminRoleBinding, organization.GetLabels, organization.Config.OrganizationInitConfig.AdminRoleBindingFile)
	organization.AdminClusterRoleBindingManager = mgr.CreateClusterRoleBindingManager("", override.AdminClusterRoleBinding, organization.GetLabels, organization.Config.OrganizationInitConfig.AdminClusterRoleBindingFile)

	organization.ClientRoleManager = mgr.CreateRoleManager("", override.ClientRole, organization.GetLabels, organization.Config.OrganizationInitConfig.ClientRoleFile)
}

// Reconcile on Organization upon Update
func (organization *BaseOrganization) Reconcile(instance *current.Organization, update Update) (common.Result, error) {
	var err error

	if err = organization.PreReconcileChecks(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed on prereconcile checks")
	}

	if err = organization.Initialize(instance, update); err != nil {
		return common.Result{}, operatorerrors.Wrap(err, operatorerrors.OrganizationInitilizationFailed, "failed to initialize organization")
	}

	// TODO: define managers
	if err = organization.ReconcileManagers(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}

	return organization.CheckStates(instance)
}

// PreReconcileChecks on Organization upon Update
func (organization *BaseOrganization) PreReconcileChecks(instance *current.Organization, update Update) error {
	log.Info(fmt.Sprintf("PreReconcileChecks on Organization %s", instance.GetName()))

	if !instance.HasAdmin() {
		return errors.New("organization admin is empty")
	}

	return nil
}

// Initialize on Organization upon Update
func (organization *BaseOrganization) Initialize(instance *current.Organization, update Update) error {
	log.Info(fmt.Sprintf("Checking if organization '%s' needs initialization", instance.GetName()))
	return nil
}

// ReconcileManagers on Organization upon Update
func (organization *BaseOrganization) ReconcileManagers(instance *current.Organization, update Update) error {
	var err error

	err = organization.CreateNamespace(instance)
	if err != nil {
		return err
	}

	// AdminRole
	err = organization.AdminRoleManager.Reconcile(instance, false)
	if err != nil {
		return err
	}

	// ClientRole
	err = organization.ClientRoleManager.Reconcile(instance, false)
	if err != nil {
		return err
	}

	// AdminRoleBinding
	if update.AdminUpdated() {
		err = organization.AdminRoleBindingManager.Reconcile(instance, true)
		if err != nil {
			return err
		}
		// AdminClusterRoleBinding
		err = organization.AdminClusterRoleBindingManager.Reconcile(instance, true)
		if err != nil {
			return err
		}
	}

	if update.AdminUpdated() {
		err = organization.PatchAnnotations(instance)
		if err != nil {
			return err
		}
	}
	// TODO: Deploy CA

	return nil
}

// CheckStates on Organization
func (organization *BaseOrganization) CheckStates(instance *current.Organization) (common.Result, error) {
	return common.Result{
		Status: &current.CRStatus{
			Type:    current.Created,
			Version: version.Operator,
		},
	}, nil
}

// GetLabels from instance.GetLabels
func (organization *BaseOrganization) GetLabels(instance v1.Object) map[string]string {
	return instance.GetLabels()
}

func (organization *BaseOrganization) CreateNamespace(instance *current.Organization) error {
	ns := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:   instance.GetUserNamespace(),
			Labels: instance.GetLabels(),
		},
	}
	ns.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Organization",
			APIVersion: "ibp.com/v1beta1",
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
		},
	}
	return organization.Client.CreateOrUpdate(context.TODO(), &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:   instance.GetUserNamespace(),
			Labels: instance.GetLabels(),
		},
	})
}

// Patch to annotations
func (organization *BaseOrganization) PatchAnnotations(instance *current.Organization) error {
	var err error

	iamuser := &iam.User{}
	err = organization.Client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Admin}, iamuser)
	if err != nil {
		return err
	}

	err = SetAdminAnnotations(iamuser, instance)
	if err != nil {
		return err
	}

	err = organization.Client.Patch(context.TODO(), iamuser, nil, controllerclient.PatchOption{
		Resilient: &controllerclient.ResilientPatch{
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

func SetAdminAnnotations(iamuser *iam.User, instance *current.Organization) error {
	var err error

	annotationList := &current.BlockchainAnnotationList{
		List: make(map[string]current.BlockchainAnnotation),
	}
	err = annotationList.Unmarshal([]byte(iamuser.Annotations[current.BlockchainAnnotationKey]))
	if err != nil {
		return err
	}

	adminAnnotation := instance.GetAdminAnnotations()

	_, err = annotationList.SetOrUpdateAnnotation(instance.GetName(), adminAnnotation)
	if err != nil {
		return err
	}

	raw, err := annotationList.Marshal()
	if err != nil {
		return err
	}

	iamuser.Annotations[current.BlockchainAnnotationKey] = string(raw)

	return nil
}
