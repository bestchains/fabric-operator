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
	"sync"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	"github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("predicates", func() {
	var (
		reconciler     *ReconcileOrganization
		client         *mocks.Client
		oldOrg, newOrg *current.Organization
	)

	Context("Predict create event on organization", func() {
		var (
			e event.CreateEvent
		)

		BeforeEach(func() {
			oldOrg = &current.Organization{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-org0",
				},
				Spec: current.OrganizationSpec{},
			}

			newOrg = &current.Organization{
				ObjectMeta: metav1.ObjectMeta{
					Name: oldOrg.GetName(),
				},
				Status: current.OrganizationStatus{
					CRStatus: current.CRStatus{
						Type: current.Created,
					},
				},
			}

			e = event.CreateEvent{
				Object: oldOrg,
			}

			client = &mocks.Client{
				GetStub: func(ctx context.Context, types types.NamespacedName, obj k8sclient.Object) error {
					switch obj := obj.(type) {
					case *corev1.ConfigMap:
						bytes, err := yaml.Marshal(oldOrg.Spec)
						Expect(err).NotTo((HaveOccurred()))
						obj.BinaryData = map[string][]byte{
							"spec": bytes,
						}
					}

					return nil
				},
				ListStub: func(ctx context.Context, obj k8sclient.ObjectList, opts ...k8sclient.ListOption) error {
					switch obj := obj.(type) {
					case *current.OrganizationList:
						obj.Items = []current.Organization{
							{ObjectMeta: metav1.ObjectMeta{Name: "test-org0"}},
							{ObjectMeta: metav1.ObjectMeta{Name: "test-org1"}},
							{ObjectMeta: metav1.ObjectMeta{Name: "test-org2"}},
						}
					case *current.IBPCAList:
						obj.Items = []current.IBPCA{
							{ObjectMeta: metav1.ObjectMeta{Name: "test-org0-ca"}},
							{ObjectMeta: metav1.ObjectMeta{Name: "test-org1-ca"}},
							{ObjectMeta: metav1.ObjectMeta{Name: "test-org2-ca"}},
						}
					}
					return nil
				},
			}

			reconciler = &ReconcileOrganization{
				update: map[string][]Update{},
				client: client,
				mutex:  &sync.Mutex{},
			}
		})

		It("sets update flags to true when a create event is detected", func() {
			create := reconciler.CreateFunc(e)
			Expect(create).To(Equal(true))

			Expect(reconciler.GetUpdateStatus(oldOrg)).To(Equal(&Update{
				adminOrCAUpdated: true,
			}))
		})

		It("sets update flags to true if instance has status type and a create event is detected and spec changes detected", func() {
			spec := current.OrganizationSpec{
				DisplayName: "test-org0",
				CAReference: current.CAReference{
					Name: "test-org0-ca",
					CA:   "ca",
				},
				Admin: "admin",
			}
			binaryData, err := yaml.Marshal(spec)
			Expect(err).NotTo(HaveOccurred())

			client.GetStub = func(ctx context.Context, types types.NamespacedName, obj k8sclient.Object) error {
				switch obj := obj.(type) {
				case *corev1.ConfigMap:
					obj.BinaryData = map[string][]byte{
						"spec": binaryData,
					}
				}
				return nil
			}
			create := reconciler.CreateFunc(e)
			Expect(create).To(Equal(true))

			Expect(reconciler.GetUpdateStatus(newOrg)).To(Equal(&Update{
				adminOrCAUpdated: true,
			}))
		})

		It("trigger create if instance does not have status type and a create event is detected", func() {
			newOrg.Status.Type = ""

			create := reconciler.CreateFunc(e)
			Expect(create).To(Equal(true))

			Expect(reconciler.GetUpdateStatus(newOrg)).To(Equal(&Update{
				adminOrCAUpdated: true,
			}))
		})
	})
	Context("Predict update event on organization", func() {
		var (
			e event.UpdateEvent
		)

		BeforeEach(func() {
			oldOrg = &current.Organization{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-org0",
				},
			}

			newOrg = &current.Organization{
				ObjectMeta: metav1.ObjectMeta{
					Name: oldOrg.Name,
				},
			}

			e = event.UpdateEvent{
				ObjectOld: oldOrg,
				ObjectNew: newOrg,
			}

			client = &mocks.Client{}
			reconciler = &ReconcileOrganization{
				update: map[string][]Update{},
				client: client,
				mutex:  &sync.Mutex{},
			}
		})
		It("returns false old and new objects are equal", func() {
			Expect(reconciler.UpdateFunc(e)).To(Equal(false))
		})

		It("returns true if spec updated", func() {
			newOrg.Spec.DisplayName = "test-org0"
			newOrg.Spec.Admin = "admin"
			newOrg.Spec.CAReference = current.CAReference{
				Name: "org0-ca",
				CA:   "ca",
			}
			Expect(reconciler.UpdateFunc(e)).To(Equal(true))
			Expect(reconciler.GetUpdateStatus(newOrg).AdminOrCAUpdated()).To(Equal(true))
		})
	})

	Context("Predict create/update/delete event on Federation", func() {
		BeforeEach(func() {
			client = &mocks.Client{}
			reconciler = &ReconcileOrganization{
				update: map[string][]Update{},
				client: client,
				mutex:  &sync.Mutex{},
			}
		})
		federation := &current.Federation{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "federation-triple",
				Namespace: "org1",
			},
			Spec: current.FederationSpec{
				Members: []current.Member{
					{Name: "org1", Namespace: "org1", Initiator: true},
					{Name: "org2", Namespace: "org2", Initiator: false},
				},
			},
		}

		newFederation := &current.Federation{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "federation-triple",
				Namespace: "org1",
			},
			Spec: current.FederationSpec{
				Members: []current.Member{
					{Name: "org1", Namespace: "org1", Initiator: true},
					{Name: "org3", Namespace: "org3", Initiator: false},
				},
			},
		}

		It("create event: add federation to each member's status", func() {
			e := event.CreateEvent{Object: federation}

			reconcile := reconciler.CreateFunc(e)
			Expect(reconcile).To(BeFalse())
		})

		It("create event: AddFed error when organizaiton not found", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Organization:
					if nn.Name == "org1" {
						return k8serrors.NewNotFound(
							schema.GroupResource{Group: "ibp.com", Resource: "organizations"},
							"org1",
						)
					}
					obj.Name = nn.Name
					obj.Namespace = nn.Namespace
				}
				return nil
			}

			err := reconciler.AddFed(current.Member{Name: "org1"}, federation)
			Expect(k8serrors.IsNotFound(err)).To(BeTrue())
		})

		It("create event: AddFed error when organization already has federation in its status.Federations ", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Organization:
					obj.Name = nn.Name
					obj.Namespace = nn.Namespace

					obj.Status.Federations = append(obj.Status.Federations, current.NamespacedName{
						Name:      "federation-triple",
						Namespace: "org1",
					})
				}
				return nil
			}

			err := reconciler.AddFed(current.Member{Name: "org1"}, federation)
			Expect(err.Error()).To(ContainSubstring("already exist in organization"))
		})

		It("create event: AddFed succ  ", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Organization:
					obj.Name = nn.Name
					obj.Namespace = nn.Namespace
				}
				return nil
			}
			client.PatchStatusStub = func(ctx context.Context, o k8sclient.Object, p k8sclient.Patch, po ...controllerclient.PatchOption) error {
				return nil
			}

			err := reconciler.AddFed(current.Member{Name: "org1"}, federation)
			Expect(err).To(BeNil())
		})

		It("delete federation event", func() {
			e := event.DeleteEvent{Object: federation}

			reconcile := reconciler.DeleteFunc(e)
			Expect(reconcile).To(BeFalse())
		})

		It("delete federation event: DeleteFed fails due to organizaiton not found", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Organization:
					if nn.Name == "org1" {
						return k8serrors.NewNotFound(
							schema.GroupResource{Group: "ibp.com", Resource: "organizations"},
							"org1",
						)
					}
					obj.Name = nn.Name
					obj.Namespace = nn.Namespace
				}
				return nil
			}

			err := reconciler.DeleteFed(current.Member{Name: "org1"}, federation)
			Expect(k8serrors.IsNotFound(err)).To(BeTrue())
		})

		It("delete fedeartion event: DeleteFed fails due to federation not exist in organization.status.federations", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Organization:
					obj.Name = nn.Name
					obj.Namespace = nn.Namespace

					obj.Status.Federations = []current.NamespacedName{}
				}
				return nil
			}

			err := reconciler.DeleteFed(current.Member{Name: "org1"}, federation)
			Expect(err.Error()).To(ContainSubstring("not exist"))
		})
		It("delete fedeartion event: DeleteFed succ", func() {
			client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
				switch obj := o.(type) {
				case *current.Organization:
					obj.Name = nn.Name
					obj.Namespace = nn.Namespace

					obj.Status.Federations = []current.NamespacedName{
						{Name: "federation-triple", Namespace: "org1"},
					}
				}
				return nil
			}

			err := reconciler.DeleteFed(current.Member{Name: "org1"}, federation)
			Expect(err).To(BeNil())
		})

		It("update federation event: reconcile returns false", func() {
			e := event.UpdateEvent{ObjectOld: federation, ObjectNew: newFederation}

			reconcile := reconciler.UpdateFunc(e)
			Expect(reconcile).To(BeFalse())
		})

		It("update fedeartion event: add/remove organiations", func() {
			added, removed := current.DifferMembers(federation.GetMembers(), newFederation.GetMembers())
			Expect(len(added)).To(Equal(1))
			Expect(added[0].Name).To(Equal("org3"))

			Expect(len(removed)).To(Equal(1))
			Expect(removed[0].Name).To(Equal("org2"))

		})

	})

})
