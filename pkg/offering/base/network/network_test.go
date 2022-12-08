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

package network_test

import (
	"context"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	mgrmocks "github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/mocks"
	basenet "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/network"
	basenetmocks "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/network/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("BaseNetwork Reconcile Logic", func() {
	var (
		err    error
		client *mocks.Client
		o      basenet.Override

		clusterRoleManager        *mgrmocks.ResourceManager
		clusterRoleBindingManager *mgrmocks.ResourceManager

		reconciler *basenet.BaseNetwork

		instance *current.Network
		update   *basenetmocks.Update
	)
	BeforeEach(func() {
		update = &basenetmocks.Update{
			SpecUpdatedStub:   func() bool { return true },
			MemberUpdatedStub: func() bool { return true },
		}

		client = &mocks.Client{}
		o = &basenetmocks.Override{}
		clusterRoleManager = &mgrmocks.ResourceManager{}
		clusterRoleBindingManager = &mgrmocks.ResourceManager{}
		reconciler = &basenet.BaseNetwork{
			Client:                    client,
			ClusterRoleManager:        clusterRoleManager,
			ClusterRoleBindingManager: clusterRoleBindingManager,
			Override:                  o,
		}

	})
	Context("Prereconcile checks", func() {
		BeforeEach(func() {
			instance = &current.Network{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-sample",
					Namespace: "org1",
				},
				Spec: current.NetworkSpec{
					Consensus:  current.NamespacedName{Name: "ibp-orderer", Namespace: "org1"},
					Federation: current.NamespacedName{Name: "federation-sample", Namespace: "org1"},
					Members: []current.Member{
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: false},
						{NamespacedName: current.NamespacedName{Name: "org3", Namespace: "org3"}, Initiator: false},
					},
				},
			}
		})
		It("failed due to missing consensus", func() {
			instance.Spec.Consensus = current.NamespacedName{}
			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("consensus is empty"))
		})
		It("failed due to missing federation", func() {
			instance.Spec.Federation = current.NamespacedName{}
			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("federation is empty"))
		})
		It("failed due to missing members", func() {
			instance.Spec.Members = []current.Member{}
			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("members is empty"))
		})
		It("failed due to federation is not activated yet", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Federation:
					obj.Status.CRStatus.Type = current.FederationPending
				}
				return nil
			}
			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("not activated yet"))

		})
		It("failed due to network contains members which are not  in federation", func() {
			fedMembers := []current.Member{
				{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: false},
				{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
			}

			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Federation:
					obj.Status.CRStatus.Type = current.FederationActivated

					obj.Spec.Members = fedMembers
				}
				return nil
			}

			err = reconciler.PreReconcileChecks(instance, update)
			Expect(err.Error()).To(ContainSubstring("not in Federation"))

			added, _ := current.DifferMembers(fedMembers, instance.GetMembers())
			Expect(len(added)).To(Equal(1))
			Expect(added[0].String()).To(Equal("org3-org3"))
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
			instance = &current.Network{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.NetworkSpec{},
			}
		})

		It("instance has type", func() {
			instance.Status.CRStatus = current.CRStatus{
				Type: current.NetworkCreated,
			}
			result, err := reconciler.CheckStates(instance)
			Expect(err).To(BeNil())
			Expect(result.Status.Type).To(Equal(current.NetworkCreated))
		})

		It("instance do not have type", func() {
			result, err := reconciler.CheckStates(instance)
			Expect(err).To(BeNil())
			Expect(result.Status.Type).To(Equal(current.Created))
		})
	})

})
