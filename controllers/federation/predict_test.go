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
	"errors"
	"sync"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
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
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org3", Namespace: "org3"}, Initiator: false},
					},
				},
			}

			e := event.UpdateEvent{ObjectOld: oldFed, ObjectNew: newFed}
			Expect(reconciler.UpdateFunc(e)).To(BeTrue())
			update := reconciler.GetUpdateStatus(oldFed)
			Expect(update.GetUpdateStackWithTrues(), "specUpdated memberChanged")
		})
	})

	Context("PredictNetwork Create/Delete", func() {
		var err error
		var federation *current.Federation
		var network *current.Network
		BeforeEach(func() {
			federation = &current.Federation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "federation-sample",
					Namespace: "org1",
				},
				Spec: current.FederationSpec{
					Description: "federation for two",
					Members: []current.Member{
						{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
						{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
					},
				},
			}

			network = &current.Network{
				ObjectMeta: v1.ObjectMeta{
					Name:      "network-sample",
					Namespace: "org1",
				},
				Spec: current.NetworkSpec{
					Federation: federation.GetName(),
					Members:    federation.GetMembers(),
					Consensus:  current.NamespacedName{Name: "ibporderer-sample", Namespace: "org1"},
				},
			}
		})
		It("PredictCreate failed due to federation not found", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				return k8serrors.NewNotFound(
					schema.GroupResource{Group: "ibp.com", Resource: "federations"},
					"org0",
				)
			}
			err = reconciler.AddNetwork(federation.GetName(), network.Name)

			Expect(k8serrors.IsNotFound(err)).To(BeTrue())
		})

		It("PredictCreate failed due to network already in federation", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Federation:
					obj.ObjectMeta = federation.ObjectMeta
					obj.Spec = federation.Spec

					federation.Status.Networks = append(federation.Status.Networks, network.Name)
					obj.Status = federation.Status
				}
				return nil
			}
			err = reconciler.AddNetwork(federation.GetName(), network.Name)

			Expect(err.Error()).To(ContainSubstring("already exist"))
		})

		It("PredictCreate failed when patch status", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Federation:
					obj.ObjectMeta = federation.ObjectMeta
					obj.Spec = federation.Spec
					obj.Status = federation.Status
				}
				return nil
			}
			errMsg := "patch status failed"
			client.PatchStatusReturns(errors.New(errMsg))
			err = reconciler.AddNetwork(federation.GetName(), network.Name)

			Expect(err.Error()).To(Equal(errMsg))
		})

		It("PredictCreate fails due to federation not found", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				return k8serrors.NewNotFound(
					schema.GroupResource{Group: "ibp.com", Resource: "federations"},
					"org0",
				)
			}
			err = reconciler.DeleteNetwork(federation.GetName(), network.Name)

			Expect(k8serrors.IsNotFound(err)).To(BeTrue())
		})

		It("PredictCreate fails due to network not in federation.status.networks", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Federation:
					obj.ObjectMeta = federation.ObjectMeta
					obj.Spec = federation.Spec

					obj.Status.CRStatus = federation.Status.CRStatus
				}
				return nil
			}
			err = reconciler.DeleteNetwork(federation.GetName(), network.Name)

			Expect(err.Error()).To(ContainSubstring("not exist"))
		})

		It("PredictCreate fails when patch status", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Federation:
					obj.ObjectMeta = federation.ObjectMeta
					obj.Spec = federation.Spec
					obj.Status = federation.Status
					obj.Status.Networks = append(obj.Status.Networks, network.Name)
				}
				return nil
			}
			errMsg := "patch status failed"
			client.PatchStatusReturns(errors.New(errMsg))
			err = reconciler.DeleteNetwork(federation.GetName(), network.Name)

			Expect(err.Error()).To(Equal(errMsg))
		})
	})
})
