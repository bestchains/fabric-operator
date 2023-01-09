/*
 * Copyright contributors to the Hyperledger Fabric Operator project
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 * 	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1beta1

import (
	"errors"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	ErrChangeProposalPurpose = errors.New("the purpose of the proposal cannot be changed")
	ErrNullProposalPurpose   = errors.New("the proposal should have a purpose")
)

// log is for logging in this package.
var proposallog = logf.Log.WithName("proposal-resource")

func (r *Proposal) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-ibp-com-v1beta1-proposal,mutating=true,failurePolicy=fail,sideEffects=None,groups=ibp.com,resources=proposals,verbs=create;update,versions=v1beta1,name=proposal.mutate.webhook,admissionReviewVersions=v1

var _ webhook.Defaulter = &Proposal{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Proposal) Default() {
	proposallog.Info("default", "name", r.Name)
	if r.Spec.StartAt.IsZero() {
		r.Spec.StartAt = metav1.Now()
	}
	if r.Spec.EndAt.IsZero() {
		r.Spec.EndAt = metav1.NewTime(time.Now().Add(time.Hour * 24))
	}
	switch strings.ToLower(string(r.Spec.Policy)) {
	case strings.ToLower(string(ALL)):
		r.Spec.Policy = ALL
	case strings.ToLower(string(Majority)):
		r.Spec.Policy = Majority
	case strings.ToLower(string(OneVoteVeto)):
		r.Spec.Policy = OneVoteVeto
	default:
		r.Spec.Policy = ALL
	}
}

//+kubebuilder:webhook:path=/validate-ibp-com-v1beta1-proposal,mutating=false,failurePolicy=fail,sideEffects=None,groups=ibp.com,resources=proposals,verbs=create;update;delete,versions=v1beta1,name=proposal.validate.webhook,admissionReviewVersions=v1

var _ webhook.Validator = &Proposal{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Proposal) ValidateCreate() error {
	proposallog.Info("validate create", "name", r.Name)
	if r.Spec.CreateFederation == nil && r.Spec.DissolveFederation == nil &&
		r.Spec.AddMember == nil && r.Spec.DeleteMember == nil {
		return ErrNullProposalPurpose
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Proposal) ValidateUpdate(old runtime.Object) error {
	proposallog.Info("validate update", "name", r.Name)
	oldProposal := old.(*Proposal)

	if (oldProposal.Spec.CreateFederation != nil && r.Spec.CreateFederation == nil) ||
		(oldProposal.Spec.DissolveFederation != nil && r.Spec.DissolveFederation == nil) ||
		(oldProposal.Spec.AddMember != nil && r.Spec.AddMember == nil) ||
		(oldProposal.Spec.DeleteMember != nil && r.Spec.DeleteMember == nil) {
		return ErrChangeProposalPurpose
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Proposal) ValidateDelete() error {
	proposallog.Info("validate delete", "name", r.Name)

	return nil
}
