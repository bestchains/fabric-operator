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

// FederationSpec defines the desired state of Federation
// +k8s:deepcopy-gen=true
// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
type FederationSpec struct {
	// License should be accepted by the user to be able to setup console
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	License License `json:"license"`

	// Description for this Federation
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Description string `json:"description,omitempty"`

	// Members list all organization in this federation
	// True for Initiator; False for normal organizaiton
	// namespace-name
	Members []Member `json:"members,omitempty"`

	// Policy indicates the rules that this Federation make dicisions
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Policy Policy `json:"policy,omitempty"`
}

// Member in a Fedeartion
// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
type Member struct {
	NamespacedName `json:",inline"`
	Initiator      bool `json:"initiator,omitempty"`
}

// FederationStatus defines the observed state of Federation
type FederationStatus struct {
	CRStatus `json:",inline"`

	// TODO: save networks under this federation
	Networks []NamespacedName `json:"networks,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=fed
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
