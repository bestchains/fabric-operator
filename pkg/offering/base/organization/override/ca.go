package override

import (
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/initializer/ca/config"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) CertificateAuthority(object v1.Object, ca *current.IBPCA, action resources.Action) error {
	instance := object.(*current.Organization)
	switch action {
	case resources.Create, resources.Update:
		return o.CreateOrUpdateCA(instance, ca)
	}

	return nil
}

func (o *Override) CreateOrUpdateCA(instance *current.Organization, ca *current.IBPCA) error {
	var err error
	namespaced := instance.GetCA()
	ca.Namespace = namespaced.Namespace
	ca.Name = namespaced.Name

	ca.Spec.Domain = o.IngressDomain

	if o.IAMEnabled {
		// CA Image
		ca.Spec.Images.CATag = "iam"

		// CA Override
		var caOverride *config.Config
		caOverride, err = GetCAConfigOverride(ca)
		if err != nil {
			caOverride = &config.Config{}
		}
		caOverride.ServerConfig.CAConfig.IAM.Enabled = &o.IAMEnabled
		caOverride.ServerConfig.CAConfig.IAM.URL = o.IAMServer

		caOverride.ServerConfig.CAConfig.Organization = instance.GetName()

		raw, err := util.ConvertToJsonMessage(caOverride.ServerConfig)
		if err != nil {
			return err
		}
		ca.Spec.ConfigOverride.CA.Raw = *raw

		// TLSCA Override
		var tlscaOverride *config.Config
		tlscaOverride, err = GetTLSCAConfigOverride(ca)
		if err != nil {
			tlscaOverride = &config.Config{}
		}
		tlscaOverride.ServerConfig.CAConfig.IAM.Enabled = &o.IAMEnabled
		tlscaOverride.ServerConfig.CAConfig.IAM.URL = o.IAMServer

		tlscaOverride.ServerConfig.CAConfig.Organization = instance.GetName()

		raw, err = util.ConvertToJsonMessage(tlscaOverride.ServerConfig)
		if err != nil {
			return err
		}
		ca.Spec.ConfigOverride.TLSCA.Raw = *raw
	}

	ca.OwnerReferences = []v1.OwnerReference{
		{
			Kind:       "Organization",
			APIVersion: "ibp.com/v1beta1",
			Name:       instance.GetName(),
			UID:        instance.GetUID(),
		},
	}

	return nil
}

func GetCAConfigOverride(ca *current.IBPCA) (*config.Config, error) {
	if ca.Spec.ConfigOverride == nil || ca.Spec.ConfigOverride.CA == nil {
		return &config.Config{}, nil
	}

	configOverride, err := config.ReadFrom(&ca.Spec.ConfigOverride.CA.Raw)
	if err != nil {
		return nil, err
	}
	return configOverride, nil
}

func GetTLSCAConfigOverride(ca *current.IBPCA) (*config.Config, error) {
	if ca.Spec.ConfigOverride == nil || ca.Spec.ConfigOverride.TLSCA == nil {
		return &config.Config{}, nil
	}

	configOverride, err := config.ReadFrom(&ca.Spec.ConfigOverride.TLSCA.Raw)
	if err != nil {
		return nil, err
	}
	return configOverride, nil
}
