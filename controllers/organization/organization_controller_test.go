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
	"errors"
	"sync"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	orgmocks "github.com/IBM-Blockchain/fabric-operator/controllers/organization/mocks"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("ReconcileOrganization", func() {
	var (
		reconciler       *ReconcileOrganization
		request          reconcile.Request
		mockKubeClient   *mocks.Client
		mockOrgReconcile *orgmocks.OrganizationReconcile
		instance         *current.Organization
	)

	BeforeEach(func() {
		mockKubeClient = &mocks.Client{}
		mockOrgReconcile = &orgmocks.OrganizationReconcile{}
		instance = &current.Organization{
			Spec: current.OrganizationSpec{},
		}
		instance.Name = "test-org0"

		mockKubeClient.GetStub = func(ctx context.Context, types types.NamespacedName, obj client.Object) error {
			switch obj := obj.(type) {
			case *current.Organization:
				obj.Kind = KIND
				obj.Name = instance.Name

				instance.Status = obj.Status
			}
			return nil
		}

		mockKubeClient.UpdateStatusStub = func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
			switch obj := obj.(type) {
			case *current.Organization:
				instance.Status = obj.Status
			}
			return nil
		}

		mockKubeClient.ListStub = func(ctx context.Context, obj client.ObjectList, opts ...client.ListOption) error {
			switch obj := obj.(type) {
			case *current.OrganizationList:
				org0 := current.Organization{}
				org0.Name = "test-org0"
				obj.Items = []current.Organization{org0}
			}
			return nil
		}

		reconciler = &ReconcileOrganization{
			Offering: mockOrgReconcile,
			client:   mockKubeClient,
			scheme:   &runtime.Scheme{},
			update:   map[string][]Update{},
			mutex:    &sync.Mutex{},
		}
		request = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "test-namespace",
				Name:      "test",
			},
		}
	})

	Context("Reconciles", func() {
		It("does not return an error if the custom resource is 'not found'", func() {
			notFoundErr := &k8serror.StatusError{
				ErrStatus: metav1.Status{
					Reason: metav1.StatusReasonNotFound,
				},
			}
			mockKubeClient.GetReturns(notFoundErr)
			_, err := reconciler.Reconcile(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error if the request to get custom resource return any other error besides 'not found'", func() {
			alreadyExistsErr := &k8serror.StatusError{
				ErrStatus: metav1.Status{
					Message: "already exists",
					Reason:  metav1.StatusReasonAlreadyExists,
				},
			}
			mockKubeClient.GetReturns(alreadyExistsErr)
			_, err := reconciler.Reconcile(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("already exists"))
		})

		It("returns an error if offering.Reconcile failed", func() {
			errMsg := "stopping reconcile loop"
			mockOrgReconcile.ReconcileReturns(common.Result{}, errors.New(errMsg))
			_, err := reconciler.Reconcile(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(errMsg))
		})

		It("returns an error if offering.Reconcile succ but setstatus failed", func() {
			mockOrgReconcile.ReconcileReturns(common.Result{Status: &current.CRStatus{Type: current.Created}}, nil)
			errMsg := "patch status failed"
			mockKubeClient.PatchStatusReturns(errors.New(errMsg))
			_, err := reconciler.Reconcile(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(errMsg))
		})

	})

	Context("set status", func() {
		It("sets the status according to reconcileStatus", func() {
			reconciler.SetStatus(instance,
				&current.CRStatus{Type: current.Created, Status: current.True},
			)
			Expect(instance.Status.Type).To(Equal(current.Created))
			reconciler.SetStatus(instance,
				&current.CRStatus{Type: current.Error, Status: current.True},
			)
			Expect(instance.Status.Type).To(Equal(current.Error))
		})

		It("set status returns error if instance not found", func() {
			errMsg := "instance not found"
			mockKubeClient.GetReturns(errors.New(errMsg))
			err := reconciler.SetStatus(instance,
				&current.CRStatus{Type: current.Created, Status: current.True},
			)
			Expect(err.Error()).To(Equal(errMsg))
		})
	})

	Context("set err status", func() {
		It("set the status to Error", func() {
			reconcileErr := errors.New("failed on pre reconcile checks")
			reconciler.SetErrorStatus(instance, reconcileErr)
			Expect(instance.Status.Type).To(Equal(current.Error))
			Expect(instance.Status.Reason).To(Equal("errorOccurredDuringReconcile"))
			Expect(instance.Status.Message).To(Equal(reconcileErr.Error()))
		})
	})

})
