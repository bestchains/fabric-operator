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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	internalinterfaces "github.com/IBM-Blockchain/fabric-operator/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Chaincodes returns a ChaincodeInformer.
	Chaincodes() ChaincodeInformer
	// ChaincodeBuilds returns a ChaincodeBuildInformer.
	ChaincodeBuilds() ChaincodeBuildInformer
	// Channels returns a ChannelInformer.
	Channels() ChannelInformer
	// EndorsePolicies returns a EndorsePolicyInformer.
	EndorsePolicies() EndorsePolicyInformer
	// Federations returns a FederationInformer.
	Federations() FederationInformer
	// IBPCAs returns a IBPCAInformer.
	IBPCAs() IBPCAInformer
	// IBPConsoles returns a IBPConsoleInformer.
	IBPConsoles() IBPConsoleInformer
	// IBPOrderers returns a IBPOrdererInformer.
	IBPOrderers() IBPOrdererInformer
	// IBPPeers returns a IBPPeerInformer.
	IBPPeers() IBPPeerInformer
	// Networks returns a NetworkInformer.
	Networks() NetworkInformer
	// Organizations returns a OrganizationInformer.
	Organizations() OrganizationInformer
	// Proposals returns a ProposalInformer.
	Proposals() ProposalInformer
	// Votes returns a VoteInformer.
	Votes() VoteInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Chaincodes returns a ChaincodeInformer.
func (v *version) Chaincodes() ChaincodeInformer {
	return &chaincodeInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// ChaincodeBuilds returns a ChaincodeBuildInformer.
func (v *version) ChaincodeBuilds() ChaincodeBuildInformer {
	return &chaincodeBuildInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Channels returns a ChannelInformer.
func (v *version) Channels() ChannelInformer {
	return &channelInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// EndorsePolicies returns a EndorsePolicyInformer.
func (v *version) EndorsePolicies() EndorsePolicyInformer {
	return &endorsePolicyInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Federations returns a FederationInformer.
func (v *version) Federations() FederationInformer {
	return &federationInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// IBPCAs returns a IBPCAInformer.
func (v *version) IBPCAs() IBPCAInformer {
	return &iBPCAInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// IBPConsoles returns a IBPConsoleInformer.
func (v *version) IBPConsoles() IBPConsoleInformer {
	return &iBPConsoleInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// IBPOrderers returns a IBPOrdererInformer.
func (v *version) IBPOrderers() IBPOrdererInformer {
	return &iBPOrdererInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// IBPPeers returns a IBPPeerInformer.
func (v *version) IBPPeers() IBPPeerInformer {
	return &iBPPeerInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Networks returns a NetworkInformer.
func (v *version) Networks() NetworkInformer {
	return &networkInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Organizations returns a OrganizationInformer.
func (v *version) Organizations() OrganizationInformer {
	return &organizationInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Proposals returns a ProposalInformer.
func (v *version) Proposals() ProposalInformer {
	return &proposalInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Votes returns a VoteInformer.
func (v *version) Votes() VoteInformer {
	return &voteInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}