// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/federation"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common"
)

type Federation struct {
	CheckStatesStub        func(*v1beta1.Federation, federation.Update) (common.Result, error)
	checkStatesMutex       sync.RWMutex
	checkStatesArgsForCall []struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}
	checkStatesReturns struct {
		result1 common.Result
		result2 error
	}
	checkStatesReturnsOnCall map[int]struct {
		result1 common.Result
		result2 error
	}
	InitializeStub        func(*v1beta1.Federation, federation.Update) error
	initializeMutex       sync.RWMutex
	initializeArgsForCall []struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}
	initializeReturns struct {
		result1 error
	}
	initializeReturnsOnCall map[int]struct {
		result1 error
	}
	PreReconcileChecksStub        func(*v1beta1.Federation, federation.Update) error
	preReconcileChecksMutex       sync.RWMutex
	preReconcileChecksArgsForCall []struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}
	preReconcileChecksReturns struct {
		result1 error
	}
	preReconcileChecksReturnsOnCall map[int]struct {
		result1 error
	}
	ReconcileStub        func(*v1beta1.Federation, federation.Update) (common.Result, error)
	reconcileMutex       sync.RWMutex
	reconcileArgsForCall []struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}
	reconcileReturns struct {
		result1 common.Result
		result2 error
	}
	reconcileReturnsOnCall map[int]struct {
		result1 common.Result
		result2 error
	}
	ReconcileManagersStub        func(*v1beta1.Federation, federation.Update) error
	reconcileManagersMutex       sync.RWMutex
	reconcileManagersArgsForCall []struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}
	reconcileManagersReturns struct {
		result1 error
	}
	reconcileManagersReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Federation) CheckStates(arg1 *v1beta1.Federation, arg2 federation.Update) (common.Result, error) {
	fake.checkStatesMutex.Lock()
	ret, specificReturn := fake.checkStatesReturnsOnCall[len(fake.checkStatesArgsForCall)]
	fake.checkStatesArgsForCall = append(fake.checkStatesArgsForCall, struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}{arg1, arg2})
	fake.recordInvocation("CheckStates", []interface{}{arg1, arg2})
	fake.checkStatesMutex.Unlock()
	if fake.CheckStatesStub != nil {
		return fake.CheckStatesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.checkStatesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Federation) CheckStatesCallCount() int {
	fake.checkStatesMutex.RLock()
	defer fake.checkStatesMutex.RUnlock()
	return len(fake.checkStatesArgsForCall)
}

func (fake *Federation) CheckStatesCalls(stub func(*v1beta1.Federation, federation.Update) (common.Result, error)) {
	fake.checkStatesMutex.Lock()
	defer fake.checkStatesMutex.Unlock()
	fake.CheckStatesStub = stub
}

func (fake *Federation) CheckStatesArgsForCall(i int) (*v1beta1.Federation, federation.Update) {
	fake.checkStatesMutex.RLock()
	defer fake.checkStatesMutex.RUnlock()
	argsForCall := fake.checkStatesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Federation) CheckStatesReturns(result1 common.Result, result2 error) {
	fake.checkStatesMutex.Lock()
	defer fake.checkStatesMutex.Unlock()
	fake.CheckStatesStub = nil
	fake.checkStatesReturns = struct {
		result1 common.Result
		result2 error
	}{result1, result2}
}

func (fake *Federation) CheckStatesReturnsOnCall(i int, result1 common.Result, result2 error) {
	fake.checkStatesMutex.Lock()
	defer fake.checkStatesMutex.Unlock()
	fake.CheckStatesStub = nil
	if fake.checkStatesReturnsOnCall == nil {
		fake.checkStatesReturnsOnCall = make(map[int]struct {
			result1 common.Result
			result2 error
		})
	}
	fake.checkStatesReturnsOnCall[i] = struct {
		result1 common.Result
		result2 error
	}{result1, result2}
}

func (fake *Federation) Initialize(arg1 *v1beta1.Federation, arg2 federation.Update) error {
	fake.initializeMutex.Lock()
	ret, specificReturn := fake.initializeReturnsOnCall[len(fake.initializeArgsForCall)]
	fake.initializeArgsForCall = append(fake.initializeArgsForCall, struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}{arg1, arg2})
	fake.recordInvocation("Initialize", []interface{}{arg1, arg2})
	fake.initializeMutex.Unlock()
	if fake.InitializeStub != nil {
		return fake.InitializeStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.initializeReturns
	return fakeReturns.result1
}

func (fake *Federation) InitializeCallCount() int {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return len(fake.initializeArgsForCall)
}

func (fake *Federation) InitializeCalls(stub func(*v1beta1.Federation, federation.Update) error) {
	fake.initializeMutex.Lock()
	defer fake.initializeMutex.Unlock()
	fake.InitializeStub = stub
}

func (fake *Federation) InitializeArgsForCall(i int) (*v1beta1.Federation, federation.Update) {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	argsForCall := fake.initializeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Federation) InitializeReturns(result1 error) {
	fake.initializeMutex.Lock()
	defer fake.initializeMutex.Unlock()
	fake.InitializeStub = nil
	fake.initializeReturns = struct {
		result1 error
	}{result1}
}

func (fake *Federation) InitializeReturnsOnCall(i int, result1 error) {
	fake.initializeMutex.Lock()
	defer fake.initializeMutex.Unlock()
	fake.InitializeStub = nil
	if fake.initializeReturnsOnCall == nil {
		fake.initializeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.initializeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Federation) PreReconcileChecks(arg1 *v1beta1.Federation, arg2 federation.Update) error {
	fake.preReconcileChecksMutex.Lock()
	ret, specificReturn := fake.preReconcileChecksReturnsOnCall[len(fake.preReconcileChecksArgsForCall)]
	fake.preReconcileChecksArgsForCall = append(fake.preReconcileChecksArgsForCall, struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}{arg1, arg2})
	fake.recordInvocation("PreReconcileChecks", []interface{}{arg1, arg2})
	fake.preReconcileChecksMutex.Unlock()
	if fake.PreReconcileChecksStub != nil {
		return fake.PreReconcileChecksStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.preReconcileChecksReturns
	return fakeReturns.result1
}

func (fake *Federation) PreReconcileChecksCallCount() int {
	fake.preReconcileChecksMutex.RLock()
	defer fake.preReconcileChecksMutex.RUnlock()
	return len(fake.preReconcileChecksArgsForCall)
}

func (fake *Federation) PreReconcileChecksCalls(stub func(*v1beta1.Federation, federation.Update) error) {
	fake.preReconcileChecksMutex.Lock()
	defer fake.preReconcileChecksMutex.Unlock()
	fake.PreReconcileChecksStub = stub
}

func (fake *Federation) PreReconcileChecksArgsForCall(i int) (*v1beta1.Federation, federation.Update) {
	fake.preReconcileChecksMutex.RLock()
	defer fake.preReconcileChecksMutex.RUnlock()
	argsForCall := fake.preReconcileChecksArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Federation) PreReconcileChecksReturns(result1 error) {
	fake.preReconcileChecksMutex.Lock()
	defer fake.preReconcileChecksMutex.Unlock()
	fake.PreReconcileChecksStub = nil
	fake.preReconcileChecksReturns = struct {
		result1 error
	}{result1}
}

func (fake *Federation) PreReconcileChecksReturnsOnCall(i int, result1 error) {
	fake.preReconcileChecksMutex.Lock()
	defer fake.preReconcileChecksMutex.Unlock()
	fake.PreReconcileChecksStub = nil
	if fake.preReconcileChecksReturnsOnCall == nil {
		fake.preReconcileChecksReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.preReconcileChecksReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Federation) Reconcile(arg1 *v1beta1.Federation, arg2 federation.Update) (common.Result, error) {
	fake.reconcileMutex.Lock()
	ret, specificReturn := fake.reconcileReturnsOnCall[len(fake.reconcileArgsForCall)]
	fake.reconcileArgsForCall = append(fake.reconcileArgsForCall, struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
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

func (fake *Federation) ReconcileCallCount() int {
	fake.reconcileMutex.RLock()
	defer fake.reconcileMutex.RUnlock()
	return len(fake.reconcileArgsForCall)
}

func (fake *Federation) ReconcileCalls(stub func(*v1beta1.Federation, federation.Update) (common.Result, error)) {
	fake.reconcileMutex.Lock()
	defer fake.reconcileMutex.Unlock()
	fake.ReconcileStub = stub
}

func (fake *Federation) ReconcileArgsForCall(i int) (*v1beta1.Federation, federation.Update) {
	fake.reconcileMutex.RLock()
	defer fake.reconcileMutex.RUnlock()
	argsForCall := fake.reconcileArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Federation) ReconcileReturns(result1 common.Result, result2 error) {
	fake.reconcileMutex.Lock()
	defer fake.reconcileMutex.Unlock()
	fake.ReconcileStub = nil
	fake.reconcileReturns = struct {
		result1 common.Result
		result2 error
	}{result1, result2}
}

func (fake *Federation) ReconcileReturnsOnCall(i int, result1 common.Result, result2 error) {
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

func (fake *Federation) ReconcileManagers(arg1 *v1beta1.Federation, arg2 federation.Update) error {
	fake.reconcileManagersMutex.Lock()
	ret, specificReturn := fake.reconcileManagersReturnsOnCall[len(fake.reconcileManagersArgsForCall)]
	fake.reconcileManagersArgsForCall = append(fake.reconcileManagersArgsForCall, struct {
		arg1 *v1beta1.Federation
		arg2 federation.Update
	}{arg1, arg2})
	fake.recordInvocation("ReconcileManagers", []interface{}{arg1, arg2})
	fake.reconcileManagersMutex.Unlock()
	if fake.ReconcileManagersStub != nil {
		return fake.ReconcileManagersStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.reconcileManagersReturns
	return fakeReturns.result1
}

func (fake *Federation) ReconcileManagersCallCount() int {
	fake.reconcileManagersMutex.RLock()
	defer fake.reconcileManagersMutex.RUnlock()
	return len(fake.reconcileManagersArgsForCall)
}

func (fake *Federation) ReconcileManagersCalls(stub func(*v1beta1.Federation, federation.Update) error) {
	fake.reconcileManagersMutex.Lock()
	defer fake.reconcileManagersMutex.Unlock()
	fake.ReconcileManagersStub = stub
}

func (fake *Federation) ReconcileManagersArgsForCall(i int) (*v1beta1.Federation, federation.Update) {
	fake.reconcileManagersMutex.RLock()
	defer fake.reconcileManagersMutex.RUnlock()
	argsForCall := fake.reconcileManagersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Federation) ReconcileManagersReturns(result1 error) {
	fake.reconcileManagersMutex.Lock()
	defer fake.reconcileManagersMutex.Unlock()
	fake.ReconcileManagersStub = nil
	fake.reconcileManagersReturns = struct {
		result1 error
	}{result1}
}

func (fake *Federation) ReconcileManagersReturnsOnCall(i int, result1 error) {
	fake.reconcileManagersMutex.Lock()
	defer fake.reconcileManagersMutex.Unlock()
	fake.ReconcileManagersStub = nil
	if fake.reconcileManagersReturnsOnCall == nil {
		fake.reconcileManagersReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.reconcileManagersReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Federation) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkStatesMutex.RLock()
	defer fake.checkStatesMutex.RUnlock()
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	fake.preReconcileChecksMutex.RLock()
	defer fake.preReconcileChecksMutex.RUnlock()
	fake.reconcileMutex.RLock()
	defer fake.reconcileMutex.RUnlock()
	fake.reconcileManagersMutex.RLock()
	defer fake.reconcileManagersMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Federation) recordInvocation(key string, args []interface{}) {
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

var _ federation.Federation = new(Federation)