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

package vote

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	"github.com/IBM-Blockchain/fabric-operator/pkg/global"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	k8svote "github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/vote"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_vote")

// Add creates a new Vote Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, config *config.Config) error {
	r, err := newReconciler(mgr, config)
	if err != nil {
		return err
	}
	return add(mgr, r)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, cfg *config.Config) (*ReconcileVote, error) {
	client := k8sclient.New(mgr.GetClient(), &global.ConfigSetter{Config: cfg.Operator.Globals})
	scheme := mgr.GetScheme()

	vote := &ReconcileVote{
		client:   client,
		scheme:   scheme,
		Config:   cfg,
		voted:    make(map[string]bool),
		finished: make(map[string]bool),
		mutex:    &sync.RWMutex{},
	}

	switch cfg.Offering {
	case offering.K8S:
		vote.Offering = k8svote.New(client, scheme, cfg)
	default: // todo: maybe add OPENSHIFT support?
		return nil, errors.Errorf("offering %s not supported in Organization controller", cfg.Offering)
	}

	return vote, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileVote) error {
	predicateFuncs := predicate.Funcs{
		UpdateFunc: r.UpdateFunc,
	}

	c, err := controller.New("vote-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	if err = c.Watch(&source.Kind{Type: &current.Vote{}}, &handler.EnqueueRequestForObject{}, predicateFuncs); err != nil {
		return err
	}

	proposalPredicateFuncs := predicate.Funcs{
		UpdateFunc: r.ProposalUpdateFunc,
	}
	// Watch for changes to secondary resource proposal
	if err = c.Watch(&source.Kind{Type: &current.Proposal{}}, &handler.EnqueueRequestForObject{}, proposalPredicateFuncs); err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileVote{}

//go:generate counterfeiter -o mocks/votereconcile.go -fake-name VoteReconcile . voteReconcile

type voteReconcile interface {
	Reconcile(*current.Vote) (common.Result, error)
}

// ReconcileVote reconciles a Vote object
type ReconcileVote struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client k8sclient.Client
	scheme *runtime.Scheme

	Offering voteReconcile
	Config   *config.Config

	voted    map[string]bool // key:voteName
	finished map[string]bool // key:proposalName
	mutex    *sync.RWMutex
}

// Reconcile reads that state of the cluster for a Vote object and makes changes based on the state read
// and what is in the Vote.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileVote) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	var err error

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Vote")

	instance := &current.Vote{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	result, err := r.Offering.Reconcile(instance)
	setStatusErr := r.SetStatus(context.TODO(), instance)
	if setStatusErr != nil {
		return reconcile.Result{}, operatorerrors.IsBreakingError(setStatusErr, "failed to update status", log)
	}

	if err != nil {
		return reconcile.Result{}, operatorerrors.IsBreakingError(errors.Wrapf(err, "Vote instance '%s' encountered error", instance.GetName()), "stopping reconcile loop", log)
	}

	reqLogger.Info(fmt.Sprintf("Finished reconciling Vote '%s'", instance.GetName()))
	return result.Result, nil
}

func (r *ReconcileVote) SetStatus(ctx context.Context, instance *current.Vote) (err error) {
	if err = r.client.Get(ctx, types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()}, instance); err != nil {
		return err
	}

	if instance.Status.Phase == "" {
		instance.Status.Phase = current.VoteCreated
	}

	if r.GetVoted(instance.GetName()) {
		instance.Status.Phase = current.VoteVoted
		instance.Status.VoteTime = v1.NewTime(time.Now())
	}

	if r.GetFinished(instance.Spec.ProposalName) {
		instance.Status.Phase = current.VoteFinished
	}
	return r.PatchStatus(ctx, instance)
}

func (r *ReconcileVote) PatchStatus(ctx context.Context, instance client.Object) error {
	return r.client.PatchStatus(ctx, instance, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    3,
			Into:     &current.Vote{},
			Strategy: client.MergeFrom,
		},
	})
}

func (r *ReconcileVote) GetFinished(proposalName string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.finished[proposalName]
}

func (r *ReconcileVote) GetVoted(voteName string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.voted[voteName]
}

func (r *ReconcileVote) UpdateFinished(proposalName string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.finished[proposalName] = true
}

func (r *ReconcileVote) UpdateVoted(voteName string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.voted[voteName] = true
}

func (r *ReconcileVote) UpdateFunc(e event.UpdateEvent) bool {
	oldVote := e.ObjectOld.(*current.Vote)
	newVote := e.ObjectNew.(*current.Vote)

	if reflect.DeepEqual(oldVote.Spec, newVote.Spec) && reflect.DeepEqual(oldVote.Status, newVote.Status) {
		return false
	}

	if oldVote.Spec.Decision == nil && newVote.Spec.Decision != nil {
		r.UpdateVoted(newVote.GetName())
		log.Info(fmt.Sprintf("vote:%s voted\n", newVote.GetName()))
		return true
	}
	return false
}

func (r *ReconcileVote) ProposalUpdateFunc(e event.UpdateEvent) bool {
	oldProposal := e.ObjectOld.(*current.Proposal)
	newProposal := e.ObjectNew.(*current.Proposal)

	if oldProposal.Status.Phase != current.ProposalFinished && newProposal.Status.Phase == current.ProposalFinished {
		r.UpdateFinished(newProposal.GetName())
		log.Info(fmt.Sprintf("proposal:%s status update to finished\n", newProposal.GetName()))
		return true
	}
	return false
}
