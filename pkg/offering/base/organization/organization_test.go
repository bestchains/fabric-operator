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

package organization_test

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	cmocks "github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	orginit "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/organization"
	baseorg "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("BaseOrganization Reconcile Logic", func() {
	var (
		mockKubeClient *cmocks.Client
		organization   *baseorg.BaseOrganization
		instance       *current.Organization
		initializer    *mocks.InitializerOrganization
		update         *mocks.Update
	)
	BeforeEach(func() {
		mockKubeClient = &cmocks.Client{}
		update = &mocks.Update{}

		config := &config.Config{
			OrganizationInitConfig: &orginit.Config{
				StoragePath: "/tmp/orginit",
			},
		}

		organization = &baseorg.BaseOrganization{
			Config: config,
			Client: mockKubeClient,
			Scheme: &runtime.Scheme{},
		}

		initializer = &mocks.InitializerOrganization{}
		organization.Initializer = initializer

		instance = &current.Organization{
			TypeMeta: v1.TypeMeta{
				Kind: "Organization",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "org1",
				Namespace: "test-network",
			},
			Spec: current.OrganizationSpec{
				License:     current.License{Accept: true},
				DisplayName: "org1",
				Admin:       "admin",
				AdminSecret: "",
				CAReference: current.CAReference{
					Name: "org1-ca",
					CA:   "ca",
				},
			},
		}

	})

	Context("Reconcile", func() {
		It("failed due to PreReconcileChecks", func() {
			instance.Spec.Admin = ""
			_, err := organization.Reconcile(instance, update)
			Expect(err.Error()).To(ContainSubstring("organization admin is empty"))
		})

		It("failed due to PreReconcileChecks", func() {
			instance.Spec.CAReference.Name = ""
			_, err := organization.Reconcile(instance, update)
			Expect(err.Error()).To(ContainSubstring("organization caRef is empty"))
		})

		It("Initialize succ with all false in Update", func() {
			result, err := organization.Reconcile(instance, update)
			Expect(err).To(BeNil())
			Expect(result.Status.Type).To(Equal(current.Created))
		})

		It("Initialize succ with adminOrCAUpdated true in Update", func() {
			update.AdminOrCAUpdatedReturns(true)
			initializer.CreateOrUpdateOrgMSPSecretReturns(nil)
			_, err := organization.Reconcile(instance, update)
			Expect(err).To(BeNil())
		})

	})
})
