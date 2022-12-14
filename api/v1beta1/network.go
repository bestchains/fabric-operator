package v1beta1

import "os"

func init() {
	SchemeBuilder.Register(&Network{}, &NetworkList{})
}

func (network *Network) GetLabels() map[string]string {
	label := os.Getenv("OPERATOR_LABEL_PREFIX")
	if label == "" {
		label = "fabric"
	}

	return map[string]string{
		"app":                          network.GetName(),
		"creator":                      label,
		"release":                      "operator",
		"helm.sh/chart":                "ibm-" + label,
		"app.kubernetes.io/name":       label,
		"app.kubernetes.io/instance":   label + "network",
		"app.kubernetes.io/managed-by": label + "-operator",
	}
}

func (network *Network) GetMembers() []Member {
	return network.Spec.Members
}

func (network *Network) HasFederation() bool {
	return network.Spec.Federation != ""
}

func (network *Network) HasConsensus() bool {
	return network.Spec.Consensus.Name != ""
}

func (network *Network) HasType() bool {
	return network.Status.CRStatus.Type != ""
}

func (network *Network) HasMembers() bool {
	return len(network.Spec.Members) != 0
}
