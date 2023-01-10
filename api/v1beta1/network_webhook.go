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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var networklog = logf.Log.WithName("network-resource")

func (r *Network) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-ibp-com-v1beta1-network,mutating=true,failurePolicy=fail,sideEffects=None,groups=ibp.com,resources=networks,verbs=create;update,versions=v1beta1,name=network.mutate.webhook,admissionReviewVersions=v1

var _ webhook.Defaulter = &Network{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Network) Default() {
	networklog.Info("default", "name", r.Name)
	r.Spec.OrderSpec.License.Accept = true
	r.Spec.OrderSpec.OrdererType = "etcdraft"
	if r.Spec.OrderSpec.ClusterSize == 0 {
		r.Spec.OrderSpec.ClusterSize = 1
	}
}

//+kubebuilder:webhook:path=/validate-ibp-com-v1beta1-network,mutating=false,failurePolicy=fail,sideEffects=None,groups=ibp.com,resources=networks,verbs=create;update;delete,versions=v1beta1,name=network.validate.webhook,admissionReviewVersions=v1

var _ webhook.Validator = &Network{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Network) ValidateCreate() error {
	networklog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Network) ValidateUpdate(old runtime.Object) error {
	networklog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Network) ValidateDelete() error {
	networklog.Info("validate delete", "name", r.Name)

	return nil
}