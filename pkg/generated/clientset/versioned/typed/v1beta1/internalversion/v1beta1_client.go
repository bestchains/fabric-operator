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

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	"github.com/IBM-Blockchain/fabric-operator/pkg/generated/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type IbpInterface interface {
	RESTClient() rest.Interface
	ChaincodesGetter
	ChaincodeBuildsGetter
	ChannelsGetter
	EndorsePoliciesGetter
	FederationsGetter
	IBPCAsGetter
	IBPConsolesGetter
	IBPOrderersGetter
	IBPPeersGetter
	NetworksGetter
	OrganizationsGetter
	ProposalsGetter
	VotesGetter
}

// IbpClient is used to interact with features provided by the ibp.com group.
type IbpClient struct {
	restClient rest.Interface
}

func (c *IbpClient) Chaincodes() ChaincodeInterface {
	return newChaincodes(c)
}

func (c *IbpClient) ChaincodeBuilds() ChaincodeBuildInterface {
	return newChaincodeBuilds(c)
}

func (c *IbpClient) Channels() ChannelInterface {
	return newChannels(c)
}

func (c *IbpClient) EndorsePolicies() EndorsePolicyInterface {
	return newEndorsePolicies(c)
}

func (c *IbpClient) Federations() FederationInterface {
	return newFederations(c)
}

func (c *IbpClient) IBPCAs(namespace string) IBPCAInterface {
	return newIBPCAs(c, namespace)
}

func (c *IbpClient) IBPConsoles(namespace string) IBPConsoleInterface {
	return newIBPConsoles(c, namespace)
}

func (c *IbpClient) IBPOrderers(namespace string) IBPOrdererInterface {
	return newIBPOrderers(c, namespace)
}

func (c *IbpClient) IBPPeers(namespace string) IBPPeerInterface {
	return newIBPPeers(c, namespace)
}

func (c *IbpClient) Networks() NetworkInterface {
	return newNetworks(c)
}

func (c *IbpClient) Organizations() OrganizationInterface {
	return newOrganizations(c)
}

func (c *IbpClient) Proposals() ProposalInterface {
	return newProposals(c)
}

func (c *IbpClient) Votes(namespace string) VoteInterface {
	return newVotes(c, namespace)
}

// NewForConfig creates a new IbpClient for the given config.
func NewForConfig(c *rest.Config) (*IbpClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &IbpClient{client}, nil
}

// NewForConfigOrDie creates a new IbpClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *IbpClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new IbpClient for the given RESTClient.
func New(c rest.Interface) *IbpClient {
	return &IbpClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	config.APIPath = "/apis"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	if config.GroupVersion == nil || config.GroupVersion.Group != scheme.Scheme.PrioritizedVersionsForGroup("ibp.com")[0].Group {
		gv := scheme.Scheme.PrioritizedVersionsForGroup("ibp.com")[0]
		config.GroupVersion = &gv
	}
	config.NegotiatedSerializer = scheme.Codecs

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *IbpClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
