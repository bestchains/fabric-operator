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

package clusterrole

import (
	"context"
	"fmt"

	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/IBM-Blockchain/fabric-operator/pkg/util"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("clusterrole_manager")

type Manager struct {
	Client          k8sclient.Client
	Scheme          *runtime.Scheme
	ClusterRoleFile string
	Name            string

	LabelsFunc   func(v1.Object) map[string]string
	OverrideFunc func(v1.Object, *rbacv1.ClusterRole, resources.Action) error
}

func (m *Manager) GetName(instance v1.Object) string {
	return GetName(instance.GetName(), m.Name)
}

func (m *Manager) Reconcile(instance v1.Object, update bool) error {
	var err error
	name := m.GetName(instance)
	err = m.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: instance.GetNamespace()}, &rbacv1.ClusterRole{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info(fmt.Sprintf("Creating cluster role '%s'", name))

			clusterRole, err := m.GetClusterRoleBasedOnCRFromFile(instance)
			if err != nil {
				return err
			}

			err = m.Client.Create(context.TODO(), clusterRole, k8sclient.CreateOption{Owner: instance, Scheme: m.Scheme})
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if update {
		// update if exists(same like Create)
		clusterRole, err := m.GetClusterRoleBasedOnCRFromFile(instance)
		if err != nil {
			return err
		}

		if err = m.Client.Update(context.TODO(), clusterRole); err != nil {
			return err
		}

	}
	return nil
}

func (m *Manager) GetClusterRoleBasedOnCRFromFile(instance v1.Object) (*rbacv1.ClusterRole, error) {
	clusterRole, err := util.GetClusterRoleFromFile(m.ClusterRoleFile)
	if err != nil {
		log.Error(err, fmt.Sprintf("Error reading role configuration file: %s", m.ClusterRoleFile))
		return nil, err
	}

	clusterRole.Name = m.GetName(instance)
	clusterRole.Namespace = instance.GetNamespace()
	clusterRole.Labels = m.LabelsFunc(instance)

	return m.BasedOnCR(instance, clusterRole)
}

func (m *Manager) BasedOnCR(instance v1.Object, clusterRole *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	if m.OverrideFunc != nil {
		err := m.OverrideFunc(instance, clusterRole, resources.Create)
		if err != nil {
			return nil, operatorerrors.New(operatorerrors.InvalidClusterRoleCreateRequest, err.Error())
		}
	}

	return clusterRole, nil
}

func (m *Manager) Get(instance v1.Object) (client.Object, error) {
	if instance == nil {
		return nil, nil // Instance has not been reconciled yet
	}

	name := m.GetName(instance)
	clusterRole := &rbacv1.ClusterRole{}
	err := m.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: instance.GetNamespace()}, clusterRole)
	if err != nil {
		return nil, err
	}

	return clusterRole, nil
}

func (m *Manager) Exists(instance v1.Object) bool {
	_, err := m.Get(instance)

	return err == nil
}

func (m *Manager) Delete(instance v1.Object) error {
	clusterRole, err := m.Get(instance)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
	}

	if clusterRole == nil {
		return nil
	}

	err = m.Client.Delete(context.TODO(), clusterRole)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func (m *Manager) CheckState(instance v1.Object) error {
	// NO-OP
	return nil
}

func (m *Manager) RestoreState(instance v1.Object) error {
	// NO-OP
	return nil
}

func (m *Manager) SetCustomName(name string) {
	// NO-OP
}

func GetName(instanceName string, suffix ...string) string {
	if len(suffix) != 0 {
		if suffix[0] != "" {
			return fmt.Sprintf("%s-%s-clusterrole", instanceName, suffix[0])
		}
	}
	return fmt.Sprintf("%s-clusterrole", instanceName)
}
