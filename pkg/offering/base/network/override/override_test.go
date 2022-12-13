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

package override

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/clusterrolebinding"
	"github.com/IBM-Blockchain/fabric-operator/pkg/util"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Base network Overrides", func() {
	var (
		client *mocks.Client

		overrider *Override

		instance *current.Network

		cr  *rbacv1.ClusterRole
		crb *rbacv1.ClusterRoleBinding

		err error
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		overrider = &Override{
			Client:      client,
			SubjectKind: clusterrolebinding.ServiceAccount,
		}

		cr, err = util.GetClusterRoleFromFile("../../../../../definitions/network/clusterrole.yaml")
		Expect(err).NotTo(HaveOccurred())

		crb, err = util.GetClusterRoleBindingFromFile("../../../../../definitions/network/clusterrolebinding.yaml")
		Expect(err).NotTo(HaveOccurred())

		instance = &current.Network{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "network-sample",
				Namespace: "org1",
			},
			Spec: current.NetworkSpec{
				Federation: "federation-sample",
				Members: []current.Member{
					{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
					{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
					{NamespacedName: current.NamespacedName{Name: "org3", Namespace: "org3"}, Initiator: false},
				},
			},
		}

		client.GetStub = func(ctx context.Context, nn types.NamespacedName, o k8sclient.Object) error {
			switch obj := o.(type) {
			case *current.Organization:
				obj.Name = nn.Name
				obj.Namespace = nn.Namespace
				obj.Spec.Admin = nn.Name + "-admin"
			}
			return nil
		}

	})

	Context("ClusterRole", func() {
		It("creating cluster role", func() {
			err := overrider.ClusterRole(instance, cr, resources.Create)
			Expect(err).To(BeNil())
			ValidateClusterRole(instance, cr)
		})

		It("updating cluster role", func() {
			err := overrider.ClusterRole(instance, cr, resources.Update)
			Expect(err).To(BeNil())
			ValidateClusterRole(instance, cr)
		})

	})

	Context("ClusterRoleBinding", func() {
		It("creating cluster role binding", func() {
			err := overrider.ClusterRoleBinding(instance, crb, resources.Create)
			Expect(err).To(BeNil())
			ValidateClusterRoleBinding(instance, crb)
		})

		It("updating cluster role binding", func() {
			err := overrider.ClusterRoleBinding(instance, crb, resources.Create)
			Expect(err).To(BeNil())
			ValidateClusterRoleBinding(instance, crb)
		})

	})
})

func ValidateClusterRole(instance *current.Network, cr *rbacv1.ClusterRole) {
	By("setting resource name", func() {
		Expect(cr.Rules[0].ResourceNames).Should(ContainElements(instance.GetName()))
	})
}

func ValidateClusterRoleBinding(instance *current.Network, crb *rbacv1.ClusterRoleBinding) {
	By("setting subjects", func() {
		Expect(len(crb.Subjects)).To(Equal(len(instance.Spec.Members)))
	})
}
