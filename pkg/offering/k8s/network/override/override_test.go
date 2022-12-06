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

package override_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/network/override"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/network/override/mocks"
	"github.com/IBM-Blockchain/fabric-operator/pkg/util"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("K8S Network Overrides", func() {
	var (
		mockedBaseOverrider *mocks.BaseOverride
		overrider           *override.Override

		instance *current.Network

		cr  *rbacv1.ClusterRole
		crb *rbacv1.ClusterRoleBinding

		err error
	)

	BeforeEach(func() {
		mockedBaseOverrider = &mocks.BaseOverride{}
		overrider = &override.Override{
			BaseOverride: mockedBaseOverrider,
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
				Members: []current.Member{
					{Name: "org1", Namespace: "org1", Initiator: true},
					{Name: "org2", Namespace: "org3", Initiator: false},
					{Name: "org3", Namespace: "org3", Initiator: false},
				},
			},
		}

		mockedBaseOverrider.CreateClusterRoleStub = func(f *current.Network, cr *rbacv1.ClusterRole) error {
			cr.Name = f.GetNamespacedName()
			cr.Namespace = f.GetNamespace()

			cr.Rules = append(cr.Rules, rbacv1.PolicyRule{
				Verbs:         []string{"get"},
				APIGroups:     []string{"ibp.com"},
				Resources:     []string{"networks"},
				ResourceNames: []string{f.GetName()},
			})

			return nil
		}
		mockedBaseOverrider.UpdateClusterRoleStub = mockedBaseOverrider.CreateClusterRoleStub

		mockedBaseOverrider.CreateClusterRoleBindingStub = func(f *current.Network, crb *rbacv1.ClusterRoleBinding) error {
			crb.Name = f.GetNamespacedName()
			crb.Namespace = f.GetNamespace()

			crb.RoleRef = rbacv1.RoleRef{
				Kind: "ClusterRole",
				Name: f.GetNamespacedName(),
			}

			subjects := make([]rbacv1.Subject, 0, len(instance.GetMembers()))

			for _, member := range f.GetMembers() {
				subjects = append(subjects, rbacv1.Subject{
					Kind:      "ServiceAccount",
					Name:      member.Name + "-admin",
					Namespace: member.Namespace,
				})
			}

			crb.Subjects = subjects
			return nil
		}

		mockedBaseOverrider.UpdateClusterRoleBindingStub = mockedBaseOverrider.CreateClusterRoleBindingStub

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
	By("setting roleRef", func() {
		Expect(crb.RoleRef.Kind).To(Equal("ClusterRole"))
		Expect(crb.RoleRef.Name).To(Equal(instance.GetNamespacedName()))
	})

	By("setting subjects", func() {
		Expect(len(crb.Subjects)).To(Equal(len(instance.Spec.Members)))
	})
}
