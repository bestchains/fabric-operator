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
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common/secretmanager"
	initializer "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/organization"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	commonconfig "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common/config"
)

//go:generate counterfeiter -o mocks/crypto.go -fake-name Crypto . Crypto

type Crypto interface {
	GetCrypto() (*commonconfig.Response, error)
	PingCA() error
	Validate() error
}

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

func (i *Initializer) CreateOrUpdateOrgMSPSecret(instance *current.Organization) error {
	log.Info(fmt.Sprintf("CreateOrUpdateOrgMSPSecret on Organization %s", instance.GetNamespacedName()))

	var err error

	orgEnroller, err := i.GetAdminEnroller(instance)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Generate admin crypto on Organization %s", instance.GetNamespacedName()))
	adminCrypto, err := i.GenerateAdminCrypto(instance, orgEnroller)
	if err != nil {
		return err
	}

	// TODO: store admin private key to support auto-renewal
	adminData := map[string][]byte{
		"keystore": adminCrypto.Keystore,
		"signcert": adminCrypto.SignCert,
	}
	err = i.SecretManager.CreateOrUpdateSecret(instance, instance.GetAdminCryptoName(), adminData)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Generate msp crypto on Organization %s", instance.GetNamespacedName()))
	caCrypto, err := i.GetCACryptoSecret(instance)
	if err != nil {
		return err
	}
	data := map[string][]byte{
		"admincerts":       adminCrypto.SignCert,
		"ca_root_certs":    caCrypto.Data["cert.pem"],
		"tlsca_root_certs": caCrypto.Data["tls-cert.pem"],
	}
	err = i.SecretManager.CreateOrUpdateSecret(instance, instance.GetOrgMSPCryptoName(), data)
	if err != nil {
		return err
	}

	return nil
}

func (i *Initializer) GetAdminEnroller(instance *current.Organization) (Crypto, error) {
	log.Info(fmt.Sprintf("GetOrganizationCrypto on Organization %s", instance.GetNamespacedName()))

	enrollmentSpec, err := i.GetEnrollmentSpec(instance)
	if err != nil {
		return nil, err
	}

	cryptos := &commonconfig.Cryptos{}
	// Only software encryption supported for now
	if err := common.GetCommonEnrollers(cryptos, enrollmentSpec, i.GetInitStoragePath(instance)); err != nil {
		return nil, err
	}

	if cryptos.ClientAuth == nil {
		return nil, errors.Errorf("enroller not found on Organization %s", instance.GetNamespacedName())
	}

	return cryptos.ClientAuth, nil
}

func (i *Initializer) GenerateAdminCrypto(instance *current.Organization, enroller Crypto) (*commonconfig.Response, error) {
	adminCrypto, err := commonconfig.GenerateCrypto(enroller)
	defer i.RemoveInitStoragePath(instance)
	if err != nil {
		return nil, err
	}
	return adminCrypto, nil
}

func (i *Initializer) GetEnrollmentSpec(instance *current.Organization) (*current.EnrollmentSpec, error) {
	// Parse EnrollmentSecret
	secret, err := i.GetAdminSecret(instance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get admin-enroll-secret")
	}

	// Read CA's connection profile
	connectionProfile, err := i.GetCAConnectinProfile(instance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ca connection profile")
	}

	caURL, err := url.Parse(connectionProfile.Endpoints.API)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ca url")
	}

	return &current.EnrollmentSpec{
		ClientAuth: &current.Enrollment{
			CAName: instance.Spec.CAReference.CA,
			CAHost: caURL.Hostname(),
			CAPort: caURL.Port(),
			CATLS: &current.CATLS{
				CACert: connectionProfile.TLS.Cert,
			},
			EnrollID:     instance.Spec.Admin,
			EnrollSecret: string(secret.Data["enrollSecret"]),
		},
	}, nil
}

func (i *Initializer) GetAdminSecret(instance *current.Organization) (*corev1.Secret, error) {
	secret, err := i.SecretManager.GetSecret(instance.GetAdminSecretName(), instance)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.Errorf("admin secret %s not found ", instance.GetAdminSecretName())
	}
	return secret, nil
}

func (i *Initializer) GetCAConnectinProfile(instance *current.Organization) (*current.CAConnectionProfile, error) {
	cm := &corev1.ConfigMap{}
	if err := i.Client.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      instance.GetCAConnectinProfile(),
			Namespace: instance.Namespace,
		},
		cm,
	); err != nil {
		return nil, err
	}
	connectionProfile := &current.CAConnectionProfile{}
	if err := json.Unmarshal(cm.BinaryData["profile.json"], connectionProfile); err != nil {
		return nil, err
	}
	return connectionProfile, nil
}

func (i *Initializer) GetCACryptoSecret(instance *current.Organization) (*corev1.Secret, error) {
	secret, err := i.SecretManager.GetSecret(instance.GetCACryptoName(), instance)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (i *Initializer) GetInitStoragePath(instance *current.Organization) string {
	return filepath.Join("/", i.Config.StoragePath, instance.GetNamespacedName())
}

func (i *Initializer) RemoveInitStoragePath(instance *current.Organization) error {
	err := os.RemoveAll(i.GetInitStoragePath(instance))
	if err != nil {
		return err
	}
	return nil
}
