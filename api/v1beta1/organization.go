package v1beta1

import (
	"os"
)

func init() {
	SchemeBuilder.Register(&Organization{}, &OrganizationList{})
}

func (organization *Organization) GetLabels() map[string]string {
	label := os.Getenv("OPERATOR_LABEL_PREFIX")
	if label == "" {
		label = "fabric"
	}

	return map[string]string{
		"app":                          organization.GetName(),
		"creator":                      label,
		"release":                      "operator",
		"helm.sh/chart":                "ibm-" + label,
		"app.kubernetes.io/name":       label,
		"app.kubernetes.io/instance":   label + "organization",
		"app.kubernetes.io/managed-by": label + "-operator",
	}
}

func (organization *Organization) GetUserNamespace() string {
	return organization.GetName()
}

func (organization *Organization) GetCA() NamespacedName {
	return NamespacedName{Namespace: organization.GetUserNamespace(), Name: organization.GetName()}
}

func (organization *Organization) GetAdminRole() string {
	return "blockchain:admin-role"
}

func (organization *Organization) GetClientRole() string {
	return "blockchain:client-role"
}

func (organization *Organization) GetAdminClusterRole() string {
	return "blockchain:admin-cluster-role"
}

func (organization *Organization) GetRoleBinding(role string) string {
	return organization.GetUserNamespace() + ":" + role + "-binding"
}

func (organization *Organization) GetAnnotationKey() string {
	return organization.GetName()
}

func (organization *Organization) GetAdminAnnotations() BlockchainAnnotation {
	return BlockchainAnnotation{
		Organization: organization.GetName(),
		Namespace:    organization.GetUserNamespace(),

		// Enrollment parts
		EnrollmentID:           organization.Spec.Admin,
		Type:                   "admin",
		Affiliation:            "",
		RegistrarRoles:         "*",
		RegistrarDelegateRoles: "*",
		Revoker:                "*",
		IntermediateCA:         "true",
		GenCRL:                 "true",
		RegistrarAttributes:    "*",
	}
}

func (organization *Organization) HasDisplayName() bool {
	return organization.Spec.DisplayName != ""
}

func (organization *Organization) HasAdmin() bool {
	return organization.Spec.Admin != ""
}

func (organization *Organization) HasType() bool {
	return organization.Status.CRStatus.Type != ""
}

func (organizationStatus *OrganizationStatus) AddFederation(federation NamespacedName) bool {
	var conflict bool

	for _, f := range organizationStatus.Federations {
		if f.String() == federation.String() {
			conflict = true
			break
		}
	}

	if !conflict {
		organizationStatus.Federations = append(organizationStatus.Federations, federation)
	}

	return conflict
}

func (organizationStatus *OrganizationStatus) DeleteFederation(federation NamespacedName) bool {
	var exist bool
	var index int

	federations := organizationStatus.Federations

	for curr, f := range federations {
		if f.String() == federation.String() {
			exist = true
			index = curr
			break
		}
	}

	if exist {
		organizationStatus.Federations = append(federations[:index], federations[index+1:]...)

	}

	return exist
}
