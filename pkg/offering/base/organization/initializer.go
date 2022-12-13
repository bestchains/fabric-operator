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
	"os"
	"path/filepath"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common/secretmanager"
	initializer "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/organization"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

//go:generate counterfeiter -o mocks/crypto.go -fake-name Crypto . Crypto

type Initializer struct {
	Config *initializer.Config

	Scheme *runtime.Scheme
	Client k8sclient.Client

	GetLabels func(instance metav1.Object) map[string]string

	SecretManager *secretmanager.SecretManager
}

func NewInitializer(config *initializer.Config, scheme *runtime.Scheme, client k8sclient.Client, labels func(instance metav1.Object) map[string]string) *Initializer {
	secretManager := secretmanager.New(client, scheme, labels)

	return &Initializer{
		Config:        config,
		Client:        client,
		Scheme:        scheme,
		GetLabels:     labels,
		SecretManager: secretManager,
	}
}

func (i *Initializer) GetInitStoragePath(instance *current.Organization) string {
	return filepath.Join("/", i.Config.StoragePath, instance.GetName())
}

func (i *Initializer) RemoveInitStoragePath(instance *current.Organization) error {
	err := os.RemoveAll(i.GetInitStoragePath(instance))
	if err != nil {
		return err
	}
	return nil
}
