package override

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) AdminRole(object v1.Object, role *rbacv1.Role, action resources.Action) error {
	instance := object.(*current.Organization)
	switch action {
	case resources.Create, resources.Update:
		return o.CreateAdminRole(instance, role)
	}

	return nil
}

func (o *Override) CreateAdminRole(instance *current.Organization, role *rbacv1.Role) error {
	role.Namespace = instance.GetUserNamespace()
	role.Name = instance.GetAdminRole()
	role.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Organization",
			APIVersion: "ibp.com/v1beta1",
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
		},
	}
	return nil
}
