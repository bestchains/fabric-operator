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

package k8snet_test

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/network/mocks"
	basenetmocks "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/network/mocks"
	k8snet "github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("K8s Network Reconcile Logic", func() {
	var (
		network     *k8snet.Network
		baseNetwork *basenetmocks.Network
		instance    *current.Network
		update      *mocks.Update
	)

	BeforeEach(func() {

		update = &mocks.Update{}

		baseNetwork = &basenetmocks.Network{}

		network = &k8snet.Network{
			BaseNetwork: baseNetwork,
		}

		instance = &current.Network{
			TypeMeta: v1.TypeMeta{
				Kind: "Network",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "network-sample",
				Namespace: "org1",
			},
			Spec: current.NetworkSpec{
				Consensus:  current.NamespacedName{Name: "ibp-orderer", Namespace: "org1"},
				Federation: current.NamespacedName{Name: "federation-sample", Namespace: "org1"},
				Members: []current.Member{
					{Name: "org1", Namespace: "org1"},
					{Name: "org3", Namespace: "org3"},
				},
			},
		}
	})

	Context("K8s reconcile logic", func() {
		It("succ", func() {
			update.SpecUpdatedReturns(true)
			update.MemberUpdatedReturns(true)
			_, err := network.Reconcile(instance, update)
			Expect(err).To(BeNil())
		})
	})
})
