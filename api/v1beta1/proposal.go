package v1beta1

import (
	"context"
	"fmt"
	"os"

	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"k8s.io/apimachinery/pkg/types"
)

func (p *Proposal) GetVoteName(orgName string) string {
	return fmt.Sprintf("vote-%s-%s", orgName, p.GetName())
}

func (p *Proposal) GetVoteLabel() map[string]string {
	label := os.Getenv("OPERATOR_LABEL_PREFIX")
	if label == "" {
		label = "fabric"
	}

	return map[string]string{
		"app":                          p.GetName(),
		"creator":                      label,
		"release":                      "operator",
		"helm.sh/chart":                "ibm-" + label,
		"app.kubernetes.io/name":       label,
		"app.kubernetes.io/instance":   label + "proposal",
		"app.kubernetes.io/managed-by": label + "-operator",
	}
}

func (p *Proposal) GetCandidateOrganizations(ctx context.Context, client k8sclient.Client) ([]NamespacedName, error) {
	federation := &Federation{}
	if err := client.Get(ctx, types.NamespacedName{Name: p.Spec.Federation.Name, Namespace: p.Spec.Federation.Namespace}, federation); err != nil {
		return nil, err
	}
	orgs := make([]NamespacedName, 0)
	if p.Spec.CreateFederation != nil || p.Spec.DissolveFederation != nil {
		for _, o := range federation.Spec.Members {
			orgs = append(orgs, NamespacedName{
				Name:      o.Name,
				Namespace: o.Namespace,
			})
		}
	} else if p.Spec.AddMember != nil {
		for _, o := range federation.Spec.Members {
			orgs = append(orgs, NamespacedName{
				Name:      o.Name,
				Namespace: o.Namespace,
			})
		}
		for _, o := range p.Spec.AddMember.Members {
			orgs = append(orgs, NamespacedName{
				Name:      o.Name,
				Namespace: o.Namespace,
			})
		}
	} else if p.Spec.DeleteMember != nil {
		for _, o := range federation.Spec.Members {
			if o.Name == p.Spec.DeleteMember.Member.Name && o.Namespace == p.Spec.DeleteMember.Member.Namespace {
				continue
			}
			orgs = append(orgs, NamespacedName{
				Name:      o.Name,
				Namespace: o.Namespace,
			})
		}
	}
	return orgs, nil
}
