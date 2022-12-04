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

package k8sfed_test

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/federation/mocks"
	basefedmocks "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/federation/mocks"
	k8sfed "github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/federation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("K8s Federation Reconcile Logic", func() {
	var (
		federation     *k8sfed.Federation
		baseFederation *basefedmocks.Federation
		instance       *current.Federation
		update         *mocks.Update
	)

	BeforeEach(func() {

		update = &mocks.Update{}

		baseFederation = &basefedmocks.Federation{}

		federation = &k8sfed.Federation{
			BaseFederation: baseFederation,
		}

		instance = &current.Federation{
			TypeMeta: v1.TypeMeta{
				Kind: "Federation",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "federation-sample",
				Namespace: "org1",
			},
			Spec: current.FederationSpec{
				License: current.License{Accept: true},
				Policy:  current.ALL,
				Members: []current.Member{
					{NamespacedName: current.NamespacedName{Name: "org1", Namespace: "org1"}, Initiator: true},
					{NamespacedName: current.NamespacedName{Name: "org2", Namespace: "org2"}, Initiator: false},
					{NamespacedName: current.NamespacedName{Name: "org3", Namespace: "org3"}, Initiator: false},
				},
			},
		}
	})

	Context("K8s reconcile logic", func() {
		It("succ", func() {
			update.SpecUpdatedReturns(true)
			update.MemberUpdatedReturns(true)
			_, err := federation.Reconcile(instance, update)
			Expect(err).To(BeNil())
		})
	})
})
