package v1beta1

import "os"

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

func (organization *Organization) GetNamespacedName() string {
	return organization.GetName() + "-" + organization.GetNamespace()
}

func (organization *Organization) GetCAConnectinProfile() string {
	return organization.Spec.CAReference.Name + "-connection-profile"
}

func (organization *Organization) GetAdminSecretName() string {
	if organization.Spec.AdminSecret != "" {
		return organization.Spec.AdminSecret
	}
	return organization.GetName() + "-admin-secret"
}

func (organization *Organization) GetAdminCryptoName() string {
	return organization.GetNamespacedName() + "-admin-crypto"
}

func (organization *Organization) GetOrgMSPCryptoName() string {
	return organization.GetNamespacedName() + "-organization-crypto"
}

func (organization *Organization) GetCACryptoName() string {
	return organization.Spec.CAReference.Name + "-ca-crypto"
}

func (organization *Organization) HasDisplayName() bool {
	return organization.Spec.DisplayName != ""
}

func (organization *Organization) HasCARef() bool {
	return organization.Spec.CAReference.Name != ""
}

func (organization *Organization) HasAdmin() bool {
	return organization.Spec.Admin != ""
}

func (organization *Organization) HasType() bool {
	return organization.Status.CRStatus.Type != ""
}
