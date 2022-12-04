package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// PropalsalPhase is a label for the condition of a proposal at the current time.
type ProposalPhase string

// These are the valid statuses of pods.
const (
	// ProposalPending means the pod has been accepted by the system, but not all vote has been created.
	ProposalPending ProposalPhase = "Pending"
	// ProposalVoting means all votes has been created, waiting vote by administrator.
	ProposalVoting ProposalPhase = "Voting"
	// ProposalFinished means proposal has been finished.
	ProposalFinished ProposalPhase = "Finished"
)

// ProposalCondition contains details for the current condition of this proposal.
type ProposalCondition struct {
	// Type is the type of the condition.
	Type ProposalConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status metav1.ConditionStatus `json:"status"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// ProposalConditionType is a valid value for ProposalCondition.Type
type ProposalConditionType string

// These are valid conditions of proposal.
const (
	// ProposalInitialized means proposal is created and has been accepted by the system.
	ProposalInitialized ProposalConditionType = "Initialized"
	// ProposalDeployed means all objects is Deployed by controller, for example, vote.
	ProposalDeployed ProposalConditionType = "Deployed"
	// ProposalSucceeded means the proposal was voted and adopted.
	ProposalSucceeded ProposalConditionType = "Succeeded"
	// ProposalFailed means the proposal was failed.
	ProposalFailed ProposalConditionType = "Failed"
	// ProposalExpired means the proposal is expired when waiting vote.
	ProposalExpired ProposalConditionType = "Expired"
	// ProposalError means the proposal is in error.
	ProposalError ProposalConditionType = "Error"
)

// Proposal defines all proposals that require a vote in the federation.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:deepcopy-gen=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type Proposal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Spec ProposalSpec `json:"spec,omitempty"`
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Status ProposalStatus `json:"status,omitempty"`
}

// ProposalList contains a list of Proposal.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:deepcopy-gen=true
type ProposalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Proposal `json:"items"`
}

type ProposalSpec struct {
	Federation     NamespacedName `json:"federation"`
	Policy         Policy         `json:"policy"`
	Initiator      NamespacedName `json:"initiator"`
	ProposalSource `json:",inline"`
	// +optional
	StartAt metav1.Time `json:"startAt,omitempty"`
	// +optional
	EndAt metav1.Time `json:"endAt,omitempty"`
	// +kubebuilder:default=false
	Deprecated bool `json:"deprecated,omitempty"`
}

type ProposalSource struct {
	// +optional
	CreateFederation *CreateFederation `json:"createFederation"`
	// +optional
	AddMember *AddMember `json:"addMember"`
	// +optional
	DeleteMember *DeleteMember `json:"deleteMember"`
	// +optional
	DissolveFederation *DissolveFederation `json:"dissolveFederation"`
}

type AddMember struct {
	Members []NamespacedName `json:"members"`
}

type DeleteMember struct {
	Member NamespacedName `json:"member"`
}
type DissolveFederation struct {
}

type CreateFederation struct {
}

type VoteResult struct {
	Organization NamespacedName `json:"organization"`
	// +optional
	Decision    *bool       `json:"decision"`
	Description string      `json:"description"`
	Phase       VotePhase   `json:"phase,omitempty"`
	VoteTime    metav1.Time `json:"startTime,omitempty"`
}

type ProposalStatus struct {
	// todo comment
	// +optional
	Phase ProposalPhase `json:"phase,omitempty"`
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ProposalCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// A human readable message indicating details about why the proposal is in this condition.
	// +optional
	Message string `json:"message,omitempty"`
	// A brief CamelCase message indicating details about why the proposal is in this state.
	// e.g. 'Expired'
	// +optional
	Reason string `json:"reason,omitempty"`
	// The list has one entry per init container in the manifest. The most recent successful
	// init container will have ready = true, the most recently started container will have
	// startTime set.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-and-container-status
	Votes []VoteResult `json:"votes,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Proposal{}, &ProposalList{})
}
