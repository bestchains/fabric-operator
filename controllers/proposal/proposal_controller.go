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

package proposal

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	"github.com/IBM-Blockchain/fabric-operator/pkg/global"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
	k8sproposal "github.com/IBM-Blockchain/fabric-operator/pkg/offering/k8s/proposal"
	"github.com/IBM-Blockchain/fabric-operator/pkg/operatorerrors"
	"github.com/pkg/errors"
	"k8s.io/utils/pointer"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
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

var log = logf.Log.WithName("controller_proposal")

// Add creates a new Proposal Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, config *config.Config) error {
	r, err := newReconciler(mgr, config)
	if err != nil {
		return err
	}
	return add(mgr, r)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, cfg *config.Config) (*ReconcileProposal, error) {
	client := k8sclient.New(mgr.GetClient(), &global.ConfigSetter{Config: cfg.Operator.Globals})
	scheme := mgr.GetScheme()

	proposal := &ReconcileProposal{
		client:     client,
		scheme:     scheme,
		Config:     cfg,
		voteResult: map[string]map[string]current.VoteResult{},
		mutex:      &sync.Mutex{},
	}

	switch cfg.Offering {
	case offering.K8S:
		proposal.Offering = k8sproposal.New(client, scheme, cfg)
	default: // todo: maybe add OPENSHIFT support?
		return nil, errors.Errorf("offering %s not supported in Organization controller", cfg.Offering)
	}

	return proposal, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileProposal) error {
	proposalPredicateFuncs := predicate.Funcs{
		CreateFunc: r.CreateFunc,
	}

	c, err := controller.New("proposal-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource proposal
	if err = c.Watch(&source.Kind{Type: &current.Proposal{}}, &handler.EnqueueRequestForObject{}, proposalPredicateFuncs); err != nil {
		return err
	}

	votePredicateFuncs := predicate.Funcs{
		CreateFunc: r.VoteCreateFunc,
		UpdateFunc: r.VoteUpdateFunc,
	}
	// Watch for changes to secondary resource Votes and requeue the owner Proposal
	if err = c.Watch(&source.Kind{Type: &current.Vote{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &current.Proposal{},
	}, votePredicateFuncs); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileProposal{}

//go:generate counterfeiter -o mocks/proposalreconcile.go -fake-name ProposalReconcile . proposalReconcile

type consoleReconcile interface {
	Reconcile(*current.Proposal) (common.Result, error)
}

// ReconcileProposal reconciles a proposal object
type ReconcileProposal struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client k8sclient.Client
	scheme *runtime.Scheme

	Offering consoleReconcile
	Config   *config.Config

	voteResult map[string]map[string]current.VoteResult
	mutex      *sync.Mutex
}

// Reconcile reads that state of the cluster for a proposal object and makes changes based on the state read
// and what is in the proposal.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileProposal) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	var err error

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	// Fetch the proposal instance
	instance := &current.Proposal{}
	err = r.client.Get(ctx, request.NamespacedName, instance)
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
	setStatusErr := r.SetStatus(ctx, instance, err)
	if setStatusErr != nil {
		return reconcile.Result{}, operatorerrors.IsBreakingError(setStatusErr, "failed to update status", log)
	}

	reqLogger.Info("proposal reconcile finished.")
	return result.Result, nil
}

func (r *ReconcileProposal) SetStatus(ctx context.Context, instance *current.Proposal, reconcileErr error) (err error) {
	if err = r.client.Get(ctx, types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()}, instance); err != nil {
		return err
	}

	if reconcileErr != nil {
		_ = r.UpdateCondition(&instance.Status, &current.ProposalCondition{
			Type:               current.ProposalError,
			Status:             v1.ConditionTrue,
			LastTransitionTime: v1.Now(),
			Reason:             "errorOccurredDuringReconcile",
			Message:            reconcileErr.Error(),
		})
		instance.Status.Phase = current.ProposalFinished

		log.Info(fmt.Sprintf("Updating status of Proposal custom resource to %s phase", current.ProposalFinished))
		return r.PatchStatus(ctx, instance)
	}
	if !instance.Spec.EndAt.IsZero() && instance.Spec.EndAt.Time.Before(time.Now()) {
		_ = r.UpdateCondition(&instance.Status, &current.ProposalCondition{
			Type:               current.ProposalExpired,
			Status:             v1.ConditionTrue,
			LastTransitionTime: v1.Now(),
			Reason:             "Success",
			Message:            "Success",
		})
		instance.Status.Phase = current.ProposalFinished
		return r.client.PatchStatus(ctx, instance, nil, k8sclient.PatchOption{
			Resilient: &k8sclient.ResilientPatch{
				Retry:    3,
				Into:     &current.Proposal{},
				Strategy: client.MergeFrom,
			},
		})
	}
	if instance.Status.Phase == "" {
		instance.Status.Phase = current.ProposalPending
		log.Info(fmt.Sprintf("Updating status of Proposal custom resource to %s phase", current.ProposalPending))
		return r.PatchStatus(ctx, instance)
	} else if instance.Status.Phase == current.ProposalPending {
		instance.Status.Phase = current.ProposalVoting
		log.Info(fmt.Sprintf("Updating status of Proposal custom resource to %s phase", current.ProposalVoting))
		return r.PatchStatus(ctx, instance)
	} else if instance.Status.Phase == current.ProposalVoting {
		res, err := r.GetVoteStatus(ctx, instance)
		if err != nil {
			return err
		}
		instance.Status.Votes = res
		if err = r.PatchStatus(ctx, instance); err != nil {
			return err
		}
		var proposalSuccess *bool
		switch instance.Spec.Policy {
		case current.OneVoteVeto, current.ALL: // todo 一票否决 和 全部人都同意 的区别是？
			for index, i := range res {
				if i.Decision != nil && !*i.Decision {
					proposalSuccess = pointer.Bool(false)
					break
				} else if i.Decision == nil {
					break
				}
				if index == len(res)-1 {
					proposalSuccess = pointer.Bool(true)
				}
			}
		case current.Majority:
			sum := len(res)
			agree := 0
			for _, i := range res {
				if pointer.BoolDeref(i.Decision, false) {
					agree += 1
				}
			}
			if agree*2 >= sum {
				proposalSuccess = pointer.Bool(true)
			}
		}
		if proposalSuccess != nil && *proposalSuccess {
			_ = r.UpdateCondition(&instance.Status, &current.ProposalCondition{
				Type:               current.ProposalSucceeded,
				Status:             v1.ConditionTrue,
				LastTransitionTime: v1.Now(),
				Reason:             "Success",
				Message:            "Success",
			})
			instance.Status.Phase = current.ProposalFinished
		} else if proposalSuccess != nil && !*proposalSuccess {
			_ = r.UpdateCondition(&instance.Status, &current.ProposalCondition{
				Type:               current.ProposalFailed,
				Status:             v1.ConditionTrue,
				LastTransitionTime: v1.Now(),
				Reason:             "Failed",
				Message:            "Failed",
			})
			instance.Status.Phase = current.ProposalFinished
		}
		return r.PatchStatus(ctx, instance)
	} else if instance.Status.Phase == current.ProposalFinished {
	}
	return
}

func (r *ReconcileProposal) PatchStatus(ctx context.Context, instance client.Object) error {
	return r.client.PatchStatus(ctx, instance, nil, k8sclient.PatchOption{
		Resilient: &k8sclient.ResilientPatch{
			Retry:    3,
			Into:     &current.Proposal{},
			Strategy: client.MergeFrom,
		},
	})
}

// GetCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func (r *ReconcileProposal) GetCondition(status *current.ProposalStatus, conditionType current.ProposalConditionType) (int, *current.ProposalCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}

// UpdateCondition updates existing proposal condition or creates a new one.
// Sets LastTransitionTime to now if the status has changed.
// Returns true if proposal condition has changed or has been added.
func (r *ReconcileProposal) UpdateCondition(status *current.ProposalStatus, condition *current.ProposalCondition) bool {
	condition.LastTransitionTime = v1.Now()
	// Try to find this condition.
	conditionIndex, oldCondition := r.GetCondition(status, condition.Type)

	if oldCondition == nil {
		// We are adding new pod condition.
		status.Conditions = append(status.Conditions, *condition)
		return true
	}
	// We are updating an existing condition, so we need to check if it has changed.
	if condition.Status == oldCondition.Status {
		condition.LastTransitionTime = oldCondition.LastTransitionTime
	}

	isEqual := condition.Status == oldCondition.Status &&
		condition.Reason == oldCondition.Reason &&
		condition.Message == oldCondition.Message &&
		condition.LastTransitionTime.Equal(&oldCondition.LastTransitionTime)

	status.Conditions[conditionIndex] = *condition
	// Return true if one of the fields have changed.
	return !isEqual
}

// todo labels
func (r *ReconcileProposal) getLabels(instance v1.Object) map[string]string {
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
		"app.kubernetes.io/instance":   label + "Proposal",
		"app.kubernetes.io/managed-by": label + "-operator",
	}
}

func (r *ReconcileProposal) getSelectorLabels(instance v1.Object) map[string]string {
	return map[string]string{
		"app": instance.GetName(),
	}
}

func (r *ReconcileProposal) CreateFunc(e event.CreateEvent) bool {
	proposal := e.Object.(*current.Proposal)
	// todo more validate in spec
	if proposal.Status.Phase != "" && proposal.Status.Phase != current.ProposalPending {
		log.Error(nil, "create a wrong phase proposal")
		return false
	}
	// todo phase and conditions union check
	return true
}

func (r *ReconcileProposal) UpdateFunc(e event.UpdateEvent) bool {
	oldProposal := e.ObjectOld.(*current.Proposal)
	newProposal := e.ObjectNew.(*current.Proposal)

	// todo add more valid check
	if reflect.DeepEqual(oldProposal.Spec, newProposal.Spec) {
		log.Error(errors.New("Spec update is not allowed"), "invalid proposal.spec update")
		return false
	}

	log.Info(fmt.Sprintf("Spec update detected on proposal custom resource: %s", oldProposal.Name))
	return true
}

func (r *ReconcileProposal) VoteCreateFunc(e event.CreateEvent) bool {
	vote := e.Object.(*current.Vote)
	// todo add more valid check
	if vote.Status.Phase != "" && vote.Status.Phase != current.VoteCreated {
		return false
	}

	if len(vote.OwnerReferences) != 1 || vote.OwnerReferences[0].Kind != "Proposal" {
		return false
	}
	// todo check other phase should not occur.

	voteResult := current.VoteResult{
		Organization: vote.GetOrganization(),
		Phase:        vote.Status.Phase,
	}
	r.UpdateVoteResult(vote.OwnerReferences[0].Name, vote.Spec.OrganizationName, voteResult)
	return true
}

func (r *ReconcileProposal) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&current.Proposal{}).
		Complete(r)
}

func (r *ReconcileProposal) GetVoteStatus(ctx context.Context, instance v1.Object) (res []current.VoteResult, err error) {
	res = r.GetAllVoteResult(instance.GetName())
	return res, nil
}

func (r *ReconcileProposal) VoteUpdateFunc(e event.UpdateEvent) bool {
	oldVote := e.ObjectOld.(*current.Vote)
	newVote := e.ObjectNew.(*current.Vote)

	if reflect.DeepEqual(oldVote.Spec, newVote.Spec) && reflect.DeepEqual(oldVote.Status, newVote.Status) {
		return false
	}

	if len(newVote.OwnerReferences) != 1 || newVote.OwnerReferences[0].Kind != "Proposal" {
		return false
	}
	// todo add more valid check
	if newVote.Status.Phase != current.VoteVoted {
		return false
	}

	voteResult := current.VoteResult{
		Organization: newVote.GetOrganization(), // todo vote -> org name
		Decision:     newVote.Spec.Decision,
		Description:  newVote.Spec.Description,
		Phase:        newVote.Status.Phase,
		VoteTime:     newVote.Status.VoteTime,
	}
	r.UpdateVoteResult(newVote.OwnerReferences[0].Name, newVote.Spec.OrganizationName, voteResult)
	log.Info("voted will update proposal")
	return true
}

func (r *ReconcileProposal) GetVoteResult(proposalName, organizaionName string) (res current.VoteResult) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	votes, exist := r.voteResult[proposalName]
	if !exist {
		return
	}
	res, exist = votes[organizaionName]
	return
}

func (r *ReconcileProposal) GetAllVoteResult(proposalName string) (res []current.VoteResult) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	votes, exist := r.voteResult[proposalName]
	if !exist {
		return
	}
	res = make([]current.VoteResult, 0, len(votes))
	for _, v := range votes {
		res = append(res, v)
	}
	return
}

func (r *ReconcileProposal) UpdateVoteResult(proposalName, organizaionName string, res current.VoteResult) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, exist := r.voteResult[proposalName]
	if !exist {
		r.voteResult[proposalName] = make(map[string]current.VoteResult)
	}
	r.voteResult[proposalName][organizaionName] = res
}
