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

package ibpca

import (
	"context"
	"encoding/json"
	"fmt"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/IBM-Blockchain/fabric-operator/pkg/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("ibpca_manager")

type Manager struct {
	Client k8sclient.Client
	Scheme *runtime.Scheme
	CAFile string
	Name   string

	LabelsFunc   func(v1.Object) map[string]string
	OverrideFunc func(v1.Object, *current.IBPCA, resources.Action) error
}

func (m *Manager) GetName(instance v1.Object) string {
	return GetName(instance.GetName(), m.Name)
}

func (m *Manager) Reconcile(instance v1.Object, update bool) error {
	var err error
	name := m.GetName(instance)
	err = m.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: instance.GetNamespace()}, &current.IBPCA{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info(fmt.Sprintf("Creating CA '%s'", name))

			ibpca, err := m.GetCABasedOnCRFromFile(instance)
			if err != nil {
				return err
			}

			err = m.Client.Create(context.TODO(), ibpca, k8sclient.CreateOption{Owner: instance, Scheme: m.Scheme})
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if update {
		// update if exists(same like Create)
		clusterRole, err := m.GetCABasedOnCRFromFile(instance)
		if err != nil {
			return err
		}

		if err = m.Client.Update(context.TODO(), clusterRole); err != nil {
			return err
		}

	}
	return nil
}

func (m *Manager) GetCABasedOnCRFromFile(instance v1.Object) (*current.IBPCA, error) {
	ca, err := GetCAFromFile(m.CAFile)
	if err != nil {
		log.Error(err, fmt.Sprintf("Error reading ca configuration file: %s", m.CAFile))
		return nil, err
	}

	ca.Name = m.GetName(instance)
	ca.Namespace = instance.GetNamespace()
	ca.Labels = m.LabelsFunc(instance)

	return m.BasedOnCR(instance, ca)
}

func (m *Manager) BasedOnCR(instance v1.Object, ca *current.IBPCA) (*current.IBPCA, error) {
	if m.OverrideFunc != nil {
		err := m.OverrideFunc(instance, ca, resources.Create)
		if err != nil {
			return nil, operatorerrors.New(operatorerrors.InvalidClusterRoleCreateRequest, err.Error())
		}
	}

	return ca, nil
}

func (m *Manager) Get(instance v1.Object) (client.Object, error) {
	if instance == nil {
		return nil, nil // Instance has not been reconciled yet
	}

	name := m.GetName(instance)
	ca := &current.IBPCA{}
	err := m.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: instance.GetNamespace()}, ca)
	if err != nil {
		return nil, err
	}

	return ca, nil
}

func (m *Manager) Exists(instance v1.Object) bool {
	_, err := m.Get(instance)

	return err == nil
}

func (m *Manager) Delete(instance v1.Object) error {
	ca, err := m.Get(instance)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
	}

	if ca == nil {
		return nil
	}

	err = m.Client.Delete(context.TODO(), ca)
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

func GetCAFromFile(file string) (*current.IBPCA, error) {
	jsonBytes, err := util.ConvertYamlFileToJson(file)
	if err != nil {
		return nil, err
	}

	ca := &current.IBPCA{}
	err = json.Unmarshal(jsonBytes, &ca)
	if err != nil {
		return nil, err
	}

	return ca, nil
}
