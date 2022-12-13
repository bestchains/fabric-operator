package override

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) AdminRoleBinding(object v1.Object, rb *rbacv1.RoleBinding, action resources.Action) error {
	instance := object.(*current.Organization)
	switch action {
	case resources.Create, resources.Update:
		return o.CreateAdminRoleBinding(instance, rb)
	}

	return nil
}

func (o *Override) CreateAdminRoleBinding(instance *current.Organization, rb *rbacv1.RoleBinding) error {
	rb.Namespace = instance.GetUserNamespace()
	rb.Name = instance.GetRoleBinding(instance.GetAdminRole())

	rb.Subjects = []rbacv1.Subject{
		common.GetDefaultSubject(instance.Spec.Admin, instance.Namespace, o.SubjectKind),
	}

	rb.RoleRef = rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     instance.GetAdminRole(),
	}

	rb.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Organization",
			APIVersion: "ibp.com/v1beta1",
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
		},
	}

	return nil
}
