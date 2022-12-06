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

package network

import (
	"context"
	"sync"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var _ = Describe("Predict on Network", func() {
	var (
		reconciler     *ReconcileNetwork
		client         *mocks.Client
		network        *current.Network
		updatedNetwork *current.Network
		federation     *current.Federation
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		reconciler = &ReconcileNetwork{
			update: map[string][]Update{},
			client: client,
			mutex:  &sync.Mutex{},
		}

		federation = &current.Federation{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "federation-sample",
				Namespace: "org1",
			},
			Spec: current.FederationSpec{
				Description: "federation for two",
				Members: []current.Member{
					{Name: "org1", Namespace: "org1", Initiator: true},
					{Name: "org2", Namespace: "org2", Initiator: false},
				},
			},
		}
		network = &current.Network{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "network-sample",
				Namespace: "org1",
			},
			Spec: current.NetworkSpec{
				Federation: federation.NamespacedName(),
				Members: []current.Member{
					{Name: "org1", Namespace: "org1", Initiator: true},
					{Name: "org2", Namespace: "org2", Initiator: false},
				},
				Consensus: current.NamespacedName{
					Name:      "ibporderer-org1",
					Namespace: "org1",
				},
			},
		}
		updatedNetwork = &current.Network{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "network-sample",
				Namespace: "org1",
			},
			Spec: current.NetworkSpec{
				Federation: federation.NamespacedName(),
				Members: []current.Member{
					{Name: "org1", Namespace: "org1", Initiator: true},
				},
			},
			// status detected
			Status: current.NetworkStatus{
				CRStatus: current.CRStatus{
					Type:   current.NetworkCreated,
					Status: current.True,
				},
			},
		}
	})

	Context("predict on network create event", func() {
		BeforeEach(func() {
			client.GetStub = func(ctx context.Context, types types.NamespacedName, obj k8sclient.Object) error {
				switch obj := obj.(type) {
				case *corev1.ConfigMap:
					bytes, err := yaml.Marshal(network.Spec)
					Expect(err).NotTo((HaveOccurred()))
					obj.BinaryData = map[string][]byte{
						"spec": bytes,
					}
				}

				return nil
			}
		})
		It("reconcile Create when new network comes", func() {
			e := event.CreateEvent{Object: network}
			Expect(reconciler.CreateFunc(e)).To(BeTrue())

			update := reconciler.GetUpdateStatus(network)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated  memberUpdated")
		})

		It("reconcile create when operator restart", func() {
			e := event.CreateEvent{Object: updatedNetwork}

			Expect(reconciler.CreateFunc(e)).To(BeTrue())

			update := reconciler.GetUpdateStatus(network)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated  memberUpdated")
		})
	})

	Context("predict on network update event", func() {
		It("reconcile false when spec not changed", func() {
			e := event.UpdateEvent{ObjectOld: network, ObjectNew: network}
			Expect(reconciler.UpdateFunc(e)).To(BeFalse())
		})

		It("reconcile true when members changed", func() {
			e := event.UpdateEvent{ObjectOld: network, ObjectNew: updatedNetwork}
			Expect(reconciler.UpdateFunc(e)).To(BeTrue())
			update := reconciler.GetUpdateStatus(network)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated memberChanged")
		})

	})

})
