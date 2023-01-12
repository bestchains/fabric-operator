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

package k8sorg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	config1 "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common/config"
	"github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common/enroller"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	baseorg "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/cloudflare/cfssl/log"
	"github.com/pkg/errors"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	maxRetryCount = 10
	mspConfig     = `
NodeOUs:
  Enable: true
  lientOUIdentifier:
	Certificate: cacerts/ca-signcert.pem
	OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
	Certificate: cacerts/ca-signcert.pem
	OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
	Certificate: cacerts/ca-signcert.pem
	OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
	Certificate: cacerts/ca-signcert.pem
	OrganizationalUnitIdentifier: orderer`
)

var tryEnrollLogSuffix = []string{"st", "nd", "rd"}

var _ baseorg.Organization = &Organization{}

type Organization struct {
	*baseorg.BaseOrganization
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config) *Organization {
	organization := &Organization{
		BaseOrganization: baseorg.New(client, scheme, config),
	}
	return organization
}

func (organization *Organization) Reconcile(instance *current.Organization, update baseorg.Update) (common.Result, error) {
	var err error

	if err = organization.PreReconcileChecks(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed on prereconcile checks")
	}

	if err = organization.Initialize(instance, update); err != nil {
		return common.Result{}, operatorerrors.Wrap(err, operatorerrors.OrganizationInitilizationFailed, "failed to initialize organization")
	}

	if err = organization.ReconcileManagers(instance, update); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}
	if update.TokenUpdated() {
		go organization.DoEnroll(context.TODO(), instance)
	}

	return organization.CheckStates(instance)
}

// TODO: customize for kubernetes

// PreReconcileChecks on Organization
func (organization *Organization) PreReconcileChecks(instance *current.Organization, update baseorg.Update) error {
	return organization.BaseOrganization.PreReconcileChecks(instance, update)
}

// Initialize on Organization after PreReconcileChecks
func (organization *Organization) Initialize(instance *current.Organization, update baseorg.Update) error {
	return organization.BaseOrganization.Initialize(instance, update)
}

// ReconcileManagers on Organization after Initialize
func (organization *Organization) ReconcileManagers(instance *current.Organization, update baseorg.Update) error {
	return organization.BaseOrganization.ReconcileManagers(instance, update)
}

// CheckStates on Organization after ReconcileManagers
func (organization *Organization) CheckStates(instance *current.Organization) (common.Result, error) {
	return organization.BaseOrganization.CheckStates(instance)
}

func (organization *Organization) enrollArg(ctx context.Context, namespace, name, enrollUser, enrollToken string) (*current.Enrollment, *current.CAConnectionProfile, error) {
	caDeploy := v12.Deployment{}
	var err error
	err = wait.Poll(10*time.Second, 10*maxRetryCount*time.Second, func() (bool, error) {
		err := organization.BaseOrganization.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &caDeploy)
		if err != nil {
			if !k8serrors.IsNotFound(err) {
				log.Error(err, "Unknown error, stop waiting.\n")
				return false, err
			}
			log.Error(err, "No deployment found, keep waiting...\n")
			return false, nil
		}

		if caDeploy.Status.AvailableReplicas != *caDeploy.Spec.Replicas {
			log.Error(fmt.Sprintf("deploy %s is not ready\n", name))
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		log.Info(fmt.Sprintf("after %d attempts, still can't get the information of ca running properly.\n", maxRetryCount))
		return nil, nil, fmt.Errorf("exceed the maximum number of retry %d", maxRetryCount)
	}
	caConnProfile := v1.ConfigMap{}
	caConnProfileName := name + "-connection-profile"
	if err = organization.BaseOrganization.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: caConnProfileName}, &caConnProfile); err != nil {
		log.Error(err, "the deployment of the ca is complete, but the configmap with the connection information is not found")
		return nil, nil, err
	}
	if caConnProfile.BinaryData == nil || len(caConnProfile.BinaryData["profile.json"]) == 0 {
		err = fmt.Errorf("can't find profile json")
		log.Error(err, "configmap does not contain any connection information")
		return nil, nil, err
	}

	profileBytes := caConnProfile.BinaryData["profile.json"]
	profile := current.CAConnectionProfile{}
	if err = json.Unmarshal(profileBytes, &profile); err != nil {
		log.Error(err, "unmarshal ca profile error")
		return nil, nil, err
	}

	// trim protocol
	host := strings.TrimPrefix(profile.Endpoints.API, "https://")
	host = strings.TrimPrefix(host, "http://")
	i := 0
	for ; i < len(host) && host[i] != ':'; i++ {
	}
	dstHost := host[:i]
	dstPort := host[i+1:]
	req := &current.Enrollment{
		CAHost: dstHost,
		CAPort: dstPort,
		CAName: name,
		CATLS: &current.CATLS{
			CACert: profile.TLS.Cert,
		},
		EnrollID:    enrollUser,
		EnrollToken: enrollToken,
		EnrollUser:  enrollUser,
	}
	log.Info("generate enrollment done.\n")
	return req, &profile, nil
}

func (organization *Organization) DoEnroll(ctx context.Context, instance *current.Organization) {
	log.Info("starting do enroll for admin and org")

	enrollment, profileStr, err := organization.enrollArg(ctx, instance.GetName(), instance.GetName(), instance.Spec.Admin, instance.Spec.AdminToken)
	if err != nil {
		return
	}
	certBytes, _ := enrollment.GetCATLSBytes()
	storagePath := organization.BaseOrganization.Config.OrganizationInitConfig.StoragePath + fmt.Sprintf("/%s/%s", instance.GetName(), instance.Spec.Admin)
	defer os.RemoveAll(storagePath)

	caClient := enroller.NewFabCAClient(enrollment, storagePath, nil, certBytes)
	certEnroller := enroller.New(enroller.NewSWEnroller(caClient))

	resp, err := config1.GenerateCrypto(certEnroller)
	if err != nil {
		log.Error(err, "GenerateCrypto error")
		return
	}

	s := v1.Secret{}
	s.Name = instance.GetName() + "-msp-crypto"
	s.Namespace = instance.GetName()
	s.Data = make(map[string][]byte)
	s.Data["admin-signcert"] = resp.SignCert
	s.Data["admin-keystore"] = resp.Keystore

	s.Data["org-ca-signcert"] = []byte(profileStr.CA.SignCerts)
	s.Data["org-tlsca-signcert"] = []byte(profileStr.TLSCA.SignCerts)
	s.Data["msg-config"] = []byte(mspConfig)

	if err = organization.BaseOrganization.Client.Create(ctx, &s); err != nil {
		if !k8serrors.IsAlreadyExists(err) {
			log.Error(err, "create secret %s error.", s.GetName())
			return
		}
		log.Info(fmt.Sprintf("secret %s already exists, try to update...", s.GetName()))

		if organization.BaseOrganization.Client.Get(ctx, types.NamespacedName{Namespace: instance.GetName(), Name: s.GetName()}, &s); err != nil {
			log.Error(fmt.Sprintf("get secret %s error %s", s.GetName(), err))
			return
		}

		s.Data["signcert"] = resp.SignCert
		s.Data["keystore"] = resp.Keystore
		s.Data["org-ca-signcert"] = []byte(profileStr.CA.SignCerts)
		s.Data["org-tlsca-signcert"] = []byte(profileStr.TLSCA.SignCerts)
		s.Data["msg-config"] = []byte(mspConfig)

		_ = organization.BaseOrganization.Client.Update(ctx, &s)
	}
	log.Info(fmt.Sprintf("create or update secret %s successfully", s.GetName()))
}
