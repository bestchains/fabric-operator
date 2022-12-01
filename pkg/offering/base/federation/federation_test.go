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

package federation_test

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	mgrmocks "github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/mocks"
	basefed "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/federation"
	basefedmocks "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/federation/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("BaseFederation Reconcile Logic", func() {
	var (
		err    error
		client *mocks.Client
		o      basefed.Override

		clusterRoleManager        *mgrmocks.ResourceManager
		clusterRoleBindingManager *mgrmocks.ResourceManager

		reconciler *basefed.BaseFederation

		instance *current.Federation
		update   *basefedmocks.Update
	)
	BeforeEach(func() {
		update = &basefedmocks.Update{
			SpecUpdatedStub:   func() bool { return true },
			MemberUpdatedStub: func() bool { return true },
		}

		client = &mocks.Client{}
		o = &basefedmocks.Override{}
		clusterRoleManager = &mgrmocks.ResourceManager{}
		clusterRoleBindingManager = &mgrmocks.ResourceManager{}
		reconciler = &basefed.BaseFederation{
			Client:                    client,
			ClusterRoleManager:        clusterRoleManager,
			ClusterRoleBindingManager: clusterRoleBindingManager,
			Override:                  o,
		}

	})
	Context("Prereconcile checks", func() {
		BeforeEach(func() {
			instance = &current.Federation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Description: "federation for unit test",
				},
			}
		})
		It("error on missing initiator", func() {
			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("federation initiator is empty"))
		})
		It("error on multiple initiator", func() {
			instance.Spec.Members = []current.Member{
				{Name: "org1", Namespace: "org1", Initiator: true},
				{Name: "org2", Namespace: "org2", Initiator: true},
			}
			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("only allow one initiator"))
		})
		It("missing policy", func() {
			instance.Spec.Members = []current.Member{
				{Name: "org1", Namespace: "org1", Initiator: true},
				{Name: "org2", Namespace: "org2", Initiator: false},
			}
			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("federation policy is empty"))
		})
	})

	Context("Reconcile managers", func() {
		It("succ", func() {
			err = reconciler.ReconcileManagers(instance, update)
			Expect(err).To(BeNil())
		})
		It("fails due to cluster role manger fails", func() {
			clusterRoleManager.ReconcileReturns(errors.New("cluster role manager reconcile fails"))
			err = reconciler.ReconcileManagers(instance, update)
			Expect(err).To(HaveOccurred())

		})
		It("fails due to cluster role binding manger fails", func() {
			clusterRoleBindingManager.ReconcileReturns(errors.New("cluster role binding manager reconcile fails"))
			err = reconciler.ReconcileManagers(instance, update)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Check states", func() {
		BeforeEach(func() {
			instance = &current.Federation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Description: "federation for unit test",
				},
			}
		})
		It("instance has type", func() {
			instance.Status.CRStatus = current.CRStatus{
				Type: current.FederationActivated,
			}
			result, err := reconciler.CheckStates(instance)
			Expect(err).To(BeNil())
			Expect(result.Status.Type).To(Equal(current.FederationActivated))
		})

		It("instance do not have type", func() {
			result, err := reconciler.CheckStates(instance)
			Expect(err).To(BeNil())
			Expect(result.Status.Type).To(Equal(current.FederationPending))
		})
	})

})
