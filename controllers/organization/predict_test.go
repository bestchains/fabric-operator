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
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	Context("create func predict", func() {
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
				Object: newOrg,
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

			Expect(reconciler.GetUpdateStatus(newOrg)).To(Equal(&Update{
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
	Context("update func predict", func() {
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
})
