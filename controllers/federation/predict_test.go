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

package federation

import (
	"context"
	"sync"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var _ = Describe("Predict federation events", func() {
	var (
		reconciler *ReconcileFederation
		client     *mocks.Client
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		reconciler = &ReconcileFederation{
			update: map[string][]Update{},
			client: client,
			mutex:  &sync.Mutex{},
		}
	})

	Context("Create events", func() {
		It("reconcile create when operator restart", func() {
			federation := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Members: []current.Member{
						{Name: "org1", Namespace: "org1", Initiator: true},
						{Name: "org2", Namespace: "org2", Initiator: false},
					},
				},
			}
			updatedFederation := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Members: []current.Member{
						{Name: "org1", Namespace: "org1", Initiator: true},
						{Name: "org3", Namespace: "org2", Initiator: false},
					},
				},
				// status detected
				Status: current.FederationStatus{
					CRStatus: current.CRStatus{
						Type:   current.Created,
						Status: current.True,
					},
				},
			}

			client.GetStub = func(ctx context.Context, types types.NamespacedName, obj k8sclient.Object) error {
				switch obj := obj.(type) {
				case *corev1.ConfigMap:
					bytes, err := yaml.Marshal(federation.Spec)
					Expect(err).NotTo((HaveOccurred()))
					obj.BinaryData = map[string][]byte{
						"spec": bytes,
					}
				}

				return nil
			}

			e := event.CreateEvent{Object: updatedFederation}

			Expect(reconciler.CreateFunc(e)).To(BeTrue())

			update := reconciler.GetUpdateStatus(federation)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated  memberUpdated")
		})

		It("reconcile Create when new federation comes", func() {
			federation := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Members: []current.Member{
						{Name: "org1", Namespace: "org1", Initiator: true},
						{Name: "org2", Namespace: "org2", Initiator: false},
					},
				},
			}
			e := event.CreateEvent{Object: federation}

			Expect(reconciler.CreateFunc(e)).To(BeTrue())

			update := reconciler.GetUpdateStatus(federation)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated  memberUpdated")
		})

	})

	Context("Update events", func() {
		It("reconcile false when spec not changed", func() {
			oldFed := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Members: []current.Member{
						{Name: "org1", Namespace: "org1", Initiator: true},
						{Name: "org2", Namespace: "org2", Initiator: false},
					},
				},
			}

			newFed := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Members: []current.Member{
						{Name: "org1", Namespace: "org1", Initiator: true},
						{Name: "org2", Namespace: "org2", Initiator: false},
					},
				},
			}

			e := event.UpdateEvent{ObjectOld: oldFed, ObjectNew: newFed}

			Expect(reconciler.UpdateFunc(e)).To(BeFalse())
		})

		It("reconcile true when spec changed but members not changed ", func() {
			oldFed := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
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

			newFed := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Description: "Federation for test",
					Members: []current.Member{
						{Name: "org1", Namespace: "org1", Initiator: true},
						{Name: "org2", Namespace: "org2", Initiator: false},
					},
				},
			}

			e := event.UpdateEvent{ObjectOld: oldFed, ObjectNew: newFed}

			Expect(reconciler.UpdateFunc(e)).To(BeTrue())
			update := reconciler.GetUpdateStatus(oldFed)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated ")
		})

		It("reconcile true when spec changed and members  changed ", func() {
			oldFed := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
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

			newFed := &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Description: "Federation for test",
					Members: []current.Member{
						{Name: "org1", Namespace: "org1", Initiator: true},
						{Name: "org3", Namespace: "org3", Initiator: false},
					},
				},
			}

			e := event.UpdateEvent{ObjectOld: oldFed, ObjectNew: newFed}

			Expect(reconciler.UpdateFunc(e)).To(BeTrue())
			update := reconciler.GetUpdateStatus(oldFed)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated memberChanged")
		})
	})
})
