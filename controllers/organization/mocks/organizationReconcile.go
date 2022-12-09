// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	organizationa "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
)

type OrganizationReconcile struct {
	ReconcileStub        func(*v1beta1.Organization, organizationa.Update) (common.Result, error)
	reconcileMutex       sync.RWMutex
	reconcileArgsForCall []struct {
		arg1 *v1beta1.Organization
		arg2 organizationa.Update
	}
	reconcileReturns struct {
		result1 common.Result
		result2 error
	}
	reconcileReturnsOnCall map[int]struct {
		result1 common.Result
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OrganizationReconcile) Reconcile(arg1 *v1beta1.Organization, arg2 organizationa.Update) (common.Result, error) {
	fake.reconcileMutex.Lock()
	ret, specificReturn := fake.reconcileReturnsOnCall[len(fake.reconcileArgsForCall)]
	fake.reconcileArgsForCall = append(fake.reconcileArgsForCall, struct {
		arg1 *v1beta1.Organization
		arg2 organizationa.Update
	}{arg1, arg2})
	fake.recordInvocation("Reconcile", []interface{}{arg1, arg2})
	fake.reconcileMutex.Unlock()
	if fake.ReconcileStub != nil {
		return fake.ReconcileStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.reconcileReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *OrganizationReconcile) ReconcileCallCount() int {
	fake.reconcileMutex.RLock()
	defer fake.reconcileMutex.RUnlock()
	return len(fake.reconcileArgsForCall)
}

func (fake *OrganizationReconcile) ReconcileCalls(stub func(*v1beta1.Organization, organizationa.Update) (common.Result, error)) {
	fake.reconcileMutex.Lock()
	defer fake.reconcileMutex.Unlock()
	fake.ReconcileStub = stub
}

func (fake *OrganizationReconcile) ReconcileArgsForCall(i int) (*v1beta1.Organization, organizationa.Update) {
	fake.reconcileMutex.RLock()
	defer fake.reconcileMutex.RUnlock()
	argsForCall := fake.reconcileArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *OrganizationReconcile) ReconcileReturns(result1 common.Result, result2 error) {
	fake.reconcileMutex.Lock()
	defer fake.reconcileMutex.Unlock()
	fake.ReconcileStub = nil
	fake.reconcileReturns = struct {
		result1 common.Result
		result2 error
	}{result1, result2}
}

func (fake *OrganizationReconcile) ReconcileReturnsOnCall(i int, result1 common.Result, result2 error) {
	fake.reconcileMutex.Lock()
	defer fake.reconcileMutex.Unlock()
	fake.ReconcileStub = nil
	if fake.reconcileReturnsOnCall == nil {
		fake.reconcileReturnsOnCall = make(map[int]struct {
			result1 common.Result
			result2 error
		})
	}
	fake.reconcileReturnsOnCall[i] = struct {
		result1 common.Result
		result2 error
	}{result1, result2}
}

func (fake *OrganizationReconcile) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.reconcileMutex.RLock()
	defer fake.reconcileMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OrganizationReconcile) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}