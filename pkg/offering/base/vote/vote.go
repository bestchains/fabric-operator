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

package basevote

import (
	"context"
	"os"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	resourcemanager "github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/manager"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("base_vote")

type Override interface {
	RoleBinding(v1.Object, *rbacv1.RoleBinding, resources.Action) error
}

type Vote interface {
	ReconcileManagers(ctx context.Context, instance *current.Vote) error
	Reconcile(instance *current.Vote) (common.Result, error)
}

var _ Vote = &BaseVote{}

type BaseVote struct {
	Client   k8sclient.Client
	Scheme   *runtime.Scheme
	Config   *config.Config
	Override Override

	RoleManager        resources.Manager
	RoleBindingManager resources.Manager
}

func New(client k8sclient.Client, scheme *runtime.Scheme, config *config.Config, override Override) *BaseVote {
	vote := &BaseVote{
		Client:   client,
		Scheme:   scheme,
		Config:   config,
		Override: override,
	}

	vote.CreateManagers()
	return vote
}

func (c *BaseVote) CreateManagers() {
	config := c.Config.VoteConfig
	override := c.Override
	resourceManager := resourcemanager.New(c.Client, c.Scheme)
	c.RoleManager = resourceManager.CreateRoleManager("", nil, c.GetLabels, config.RoleFile)
	c.RoleBindingManager = resourceManager.CreateRoleBindingManager("", override.RoleBinding, c.GetLabels, config.RoleBindingFile)
}

func (c *BaseVote) Reconcile(instance *current.Vote) (common.Result, error) {
	if err := c.ReconcileManagers(context.TODO(), instance); err != nil {
		return common.Result{}, errors.Wrap(err, "failed to reconcile managers")
	}
	return common.Result{}, nil
}

func (c *BaseVote) ReconcileManagers(ctx context.Context, instance *current.Vote) error {
	return c.ReconcileRBAC(instance)
}

func (c *BaseVote) ReconcileRBAC(instance *current.Vote) error {
	var err error

	err = c.RoleManager.Reconcile(instance, false)
	if err != nil {
		return err
	}

	err = c.RoleBindingManager.Reconcile(instance, false)
	if err != nil {
		return err
	}

	return nil
}

func (c *BaseVote) GetLabels(instance v1.Object) map[string]string {
	label := os.Getenv("OPERATOR_LABEL_PREFIX")
	if label == "" {
		label = "fabric"
	}

	return map[string]string{
		"app":                          instance.GetName(),
		"creator":                      label,
		"release":                      "operator",
		"helm.sh/chart":                "ibm-" + label,
		"app.kubernetes.io/name":       label,
		"app.kubernetes.io/instance":   label + "Vote",
		"app.kubernetes.io/managed-by": label + "-operator",
	}
}
