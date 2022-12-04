package v1beta1

func (v *Vote) GetOrganization() NamespacedName {
	return NamespacedName{
		Name:      v.Spec.OrganizationName,
		Namespace: v.GetNamespace(),
	}
}
