package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// VotePhase is a label for the condition of a vote at the current time.
type VotePhase string

// These are the valid statuses of pods.
const (
	// VoteCreated means the vote has been created by the controller, The organization administrator
	// has not yet participated in the voting.
	VoteCreated VotePhase = "Created"
	// VoteVoted means the organization administrator has vote for the proposal.
	VoteVoted VotePhase = "Voted"
	// VoteFinished means that the proposal has been finished.
	VoteFinished VotePhase = "Finished"
)

// Vote represents an organization's position on a proposal,
// including voting results and optional reasons.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:deepcopy-gen=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type Vote struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Spec VoteSpec `json:"spec,omitempty"`
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Status VoteStatus `json:"status,omitempty"`
}

// VoteList contains a list of Vote.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:deepcopy-gen=true
type VoteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Vote `json:"items"`
}

type VoteSpec struct {
	ProposalName     string `json:"proposalName"`
	OrganizationName string `json:"organizationName"`
	// +optional
	Decision *bool `json:"decision,omitempty"`
	// +optional
	Description string `json:"description"`
}

type VoteStatus struct {
	// +optional
	Phase VotePhase `json:"phase,omitempty"`
	// Timestamp of voting.
	// +optional
	VoteTime metav1.Time `json:"voteTime,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Vote{}, &VoteList{})
}
