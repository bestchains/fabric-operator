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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FederationSpec defines the desired state of Federation
// +k8s:deepcopy-gen=true
// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
type FederationSpec struct {
	// License should be accepted by the user to be able to setup CA
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	License License `json:"license"`

	/* Federation basic info*/
	// Name to declare federation uniqueness
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Name string `json:"name,omitempty"`

	// Description states what is this federation for and how it works
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Description string `json:"description,omitempty"`

	// Initiator who initiates this federation (Initiator is reponsible for building ordering service)
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Initiator FederationMember `json:"initiator,omitempty"`

	/* Federation Members */
	// Official members in this federation
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Official map[string]FederationMember `json:"official,omitempty"`

	// Rounds for federation proposals.Every member can initiate a new round
	// - NewMember
	// - ExpelMember
	// - Dissolve this federation
	// ......
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Rounds map[string]Round `json:"rounds,omitempty"`
}

type FederationMember struct {
	// Namespace (Optional) this FederationMember owns
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Namespace string `json:"namespace,omitempty"`

	// Agent is the service account to represents this organization in this federation(Only one agent account in each federation member)
	// This account will be granted a ClusterRole(`get/list/update/patch`) to this `Federation Resource`
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Agent string `json:"agent,omitempty"`

	Initiator bool `json:"initiator,omitempty"`

	// Federation should be granted a ClusterRole(once federation member aggreed to join this federation) to retrieve `MSP`(a secret generated once CA deployed) based on `{Name_MSPID}`
	// Organization name
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Name string `json:"name,omitempty"`

	// MSPID to make organization unique in this federation
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	MSPID string `json:"mspID,omitempty"`

	// Decision to indicate whether or not this member aggrees to build this federation(Only used during federation built)
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Decision bool `json:"decision,omitempty"`
}

// Round is the mechanism for federation to make proposals and decisions
// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
type Round struct {
	// Proposal details of this VOTE Round
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Proposal `json:"proposal,omitempty"`

	// Policy to check whether this round succeeds
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	PolicyKind PolicyKind `json:"policyKind,omitempty"`

	// Votes Indexed by FederationMember.MSPID
	// Auto-filterd by official member with Decision `false ` by default during `Round` create
	Votes map[string]Vote `json:"votes,omitempty"`
}

type PolicyKind string

const (
	Any      PolicyKind = "Any"
	Majority PolicyKind = "Majority"
)

type Proposal struct {
	VoteStart metav1.Time `json:"voteStart,omitempty"`
	VoteEnd   metav1.Time `json:"voteEnd,omitempty"`

	Initiator string `json:"initiator,omitempty"`

	Kind    ProposalKind `json:"kind,omitempty"`
	Payload []byte       `json:"payload,omitempty"`
}

type ProposalKind string

const (
	NewMemberProposal ProposalKind = "NewMemberProposal"
	ExpelMember       ProposalKind = "ExpelMember"
	DissolveMember    ProposalKind = "DissolveMember"
)

// Vote describes the decision made by Federation member on a specific proposal
// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
type Vote struct {
	// Decision shows member's opinion on this proposal
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Decision bool `json:"decision,omitempty"`
	// Reason detailed member's decision
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Reason string `json:"reason,omitempty"`
	// VotedAt detailed the first time member voted at this proposal
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	VotedAt metav1.Time `json:"votedAt,omitempty"`
	// UpdatedAt detailed the time member re-vote(update its vote)d at this proposal
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}

// FederationStatus defines the observed state of Federation
// TODO: need to be detailed
type FederationStatus struct {
	// Condition - Federation
	Condition metav1.Condition `json:"condition,omitempty"`

	// RoundConditions - Round
	RoundConditions []metav1.Condition `json:"roundConditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Federation is the Schema for the federations API
type Federation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FederationSpec   `json:"spec,omitempty"`
	Status FederationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FederationList contains a list of Federation
type FederationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Federation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Federation{}, &FederationList{})
}
