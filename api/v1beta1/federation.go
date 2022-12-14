package v1beta1

import "os"

func init() {
	SchemeBuilder.Register(&Federation{}, &FederationList{})
}

func (federation *Federation) GetLabels() map[string]string {
	label := os.Getenv("OPERATOR_LABEL_PREFIX")
	if label == "" {
		label = "fabric"
	}

	return map[string]string{
		"app":                          federation.GetName(),
		"creator":                      label,
		"release":                      "operator",
		"helm.sh/chart":                "ibm-" + label,
		"app.kubernetes.io/name":       label,
		"app.kubernetes.io/instance":   label + "federation",
		"app.kubernetes.io/managed-by": label + "-operator",
	}
}

func (federation *Federation) NamespacedName() NamespacedName {
	return NamespacedName{Name: federation.Name, Namespace: federation.Namespace}
}

func (federation *Federation) GetInitiator() (int, Member) {
	for index, m := range federation.Spec.Members {
		if m.Initiator {
			return index, m
		}
	}
	return -1, Member{}
}

func (federation *Federation) GetMembers() []Member {
	return federation.Spec.Members
}

func (federation *Federation) HasInitiator() bool {
	for _, m := range federation.Spec.Members {
		if m.Initiator {
			return true
		}
	}
	return false
}

func (federation *Federation) HasMultiInitiator() bool {
	var multi bool
	for _, m := range federation.Spec.Members {
		if m.Initiator && multi {
			return true
		}
		if m.Initiator {
			multi = true
		}
	}
	return false
}

func (federation *Federation) HasPolicy() bool {
	return federation.Spec.Policy != ""
}

func (federation *Federation) HasType() bool {
	return federation.Status.CRStatus.Type != ""
}

func (federation *Federation) HasMembers() bool {
	return len(federation.Spec.Members) != 0
}

func (m *Member) GetNamespacedName() string {
	return m.Namespace + "-" + m.Name
}

func DifferMembers(old []Member, new []Member) (added []Member, removed []Member) {
	// cache in map
	oldMapper := make(map[string]Member, len(old))
	for _, m := range old {
		oldMapper[m.GetNamespacedName()] = m
	}

	// calculate differences
	for _, m := range new {

		// added: in new ,but not in old
		if _, ok := oldMapper[m.GetNamespacedName()]; !ok {
			added = append(added, m)
		}

		// delete the intersection
		delete(oldMapper, m.GetNamespacedName())
	}

	for _, m := range oldMapper {
		removed = append(removed, m)
	}

	return
}

func (federationStatus *FederationStatus) AddNetwork(network string) bool {
	var conflict bool

	for _, f := range federationStatus.Networks {
		if f == network {
			conflict = true
			break
		}
	}

	if !conflict {
		federationStatus.Networks = append(federationStatus.Networks, network)
	}

	return conflict
}

func (federationStatus *FederationStatus) DeleteNetwork(network string) bool {
	var exist bool
	var index int

	networks := federationStatus.Networks

	for curr, f := range networks {
		if f == network {
			exist = true
			index = curr
			break
		}
	}

	if exist {
		federationStatus.Networks = append(networks[:index], networks[index+1:]...)
	}

	return exist
}
