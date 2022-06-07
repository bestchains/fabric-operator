// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	baseconsole "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/console"
	"github.com/IBM-Blockchain/fabric-operator/pkg/restart"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RestartManager struct {
	ForConfigMapUpdateStub        func(v1.Object) error
	forConfigMapUpdateMutex       sync.RWMutex
	forConfigMapUpdateArgsForCall []struct {
		arg1 v1.Object
	}
	forConfigMapUpdateReturns struct {
		result1 error
	}
	forConfigMapUpdateReturnsOnCall map[int]struct {
		result1 error
	}
	ForRestartActionStub        func(v1.Object) error
	forRestartActionMutex       sync.RWMutex
	forRestartActionArgsForCall []struct {
		arg1 v1.Object
	}
	forRestartActionReturns struct {
		result1 error
	}
	forRestartActionReturnsOnCall map[int]struct {
		result1 error
	}
	TriggerIfNeededStub        func(restart.Instance) error
	triggerIfNeededMutex       sync.RWMutex
	triggerIfNeededArgsForCall []struct {
		arg1 restart.Instance
	}
	triggerIfNeededReturns struct {
		result1 error
	}
	triggerIfNeededReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *RestartManager) ForConfigMapUpdate(arg1 v1.Object) error {
	fake.forConfigMapUpdateMutex.Lock()
	ret, specificReturn := fake.forConfigMapUpdateReturnsOnCall[len(fake.forConfigMapUpdateArgsForCall)]
	fake.forConfigMapUpdateArgsForCall = append(fake.forConfigMapUpdateArgsForCall, struct {
		arg1 v1.Object
	}{arg1})
	stub := fake.ForConfigMapUpdateStub
	fakeReturns := fake.forConfigMapUpdateReturns
	fake.recordInvocation("ForConfigMapUpdate", []interface{}{arg1})
	fake.forConfigMapUpdateMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *RestartManager) ForConfigMapUpdateCallCount() int {
	fake.forConfigMapUpdateMutex.RLock()
	defer fake.forConfigMapUpdateMutex.RUnlock()
	return len(fake.forConfigMapUpdateArgsForCall)
}

func (fake *RestartManager) ForConfigMapUpdateCalls(stub func(v1.Object) error) {
	fake.forConfigMapUpdateMutex.Lock()
	defer fake.forConfigMapUpdateMutex.Unlock()
	fake.ForConfigMapUpdateStub = stub
}

func (fake *RestartManager) ForConfigMapUpdateArgsForCall(i int) v1.Object {
	fake.forConfigMapUpdateMutex.RLock()
	defer fake.forConfigMapUpdateMutex.RUnlock()
	argsForCall := fake.forConfigMapUpdateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *RestartManager) ForConfigMapUpdateReturns(result1 error) {
	fake.forConfigMapUpdateMutex.Lock()
	defer fake.forConfigMapUpdateMutex.Unlock()
	fake.ForConfigMapUpdateStub = nil
	fake.forConfigMapUpdateReturns = struct {
		result1 error
	}{result1}
}

func (fake *RestartManager) ForConfigMapUpdateReturnsOnCall(i int, result1 error) {
	fake.forConfigMapUpdateMutex.Lock()
	defer fake.forConfigMapUpdateMutex.Unlock()
	fake.ForConfigMapUpdateStub = nil
	if fake.forConfigMapUpdateReturnsOnCall == nil {
		fake.forConfigMapUpdateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.forConfigMapUpdateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *RestartManager) ForRestartAction(arg1 v1.Object) error {
	fake.forRestartActionMutex.Lock()
	ret, specificReturn := fake.forRestartActionReturnsOnCall[len(fake.forRestartActionArgsForCall)]
	fake.forRestartActionArgsForCall = append(fake.forRestartActionArgsForCall, struct {
		arg1 v1.Object
	}{arg1})
	stub := fake.ForRestartActionStub
	fakeReturns := fake.forRestartActionReturns
	fake.recordInvocation("ForRestartAction", []interface{}{arg1})
	fake.forRestartActionMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *RestartManager) ForRestartActionCallCount() int {
	fake.forRestartActionMutex.RLock()
	defer fake.forRestartActionMutex.RUnlock()
	return len(fake.forRestartActionArgsForCall)
}

func (fake *RestartManager) ForRestartActionCalls(stub func(v1.Object) error) {
	fake.forRestartActionMutex.Lock()
	defer fake.forRestartActionMutex.Unlock()
	fake.ForRestartActionStub = stub
}

func (fake *RestartManager) ForRestartActionArgsForCall(i int) v1.Object {
	fake.forRestartActionMutex.RLock()
	defer fake.forRestartActionMutex.RUnlock()
	argsForCall := fake.forRestartActionArgsForCall[i]
	return argsForCall.arg1
}

func (fake *RestartManager) ForRestartActionReturns(result1 error) {
	fake.forRestartActionMutex.Lock()
	defer fake.forRestartActionMutex.Unlock()
	fake.ForRestartActionStub = nil
	fake.forRestartActionReturns = struct {
		result1 error
	}{result1}
}

func (fake *RestartManager) ForRestartActionReturnsOnCall(i int, result1 error) {
	fake.forRestartActionMutex.Lock()
	defer fake.forRestartActionMutex.Unlock()
	fake.ForRestartActionStub = nil
	if fake.forRestartActionReturnsOnCall == nil {
		fake.forRestartActionReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.forRestartActionReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *RestartManager) TriggerIfNeeded(arg1 restart.Instance) error {
	fake.triggerIfNeededMutex.Lock()
	ret, specificReturn := fake.triggerIfNeededReturnsOnCall[len(fake.triggerIfNeededArgsForCall)]
	fake.triggerIfNeededArgsForCall = append(fake.triggerIfNeededArgsForCall, struct {
		arg1 restart.Instance
	}{arg1})
	stub := fake.TriggerIfNeededStub
	fakeReturns := fake.triggerIfNeededReturns
	fake.recordInvocation("TriggerIfNeeded", []interface{}{arg1})
	fake.triggerIfNeededMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *RestartManager) TriggerIfNeededCallCount() int {
	fake.triggerIfNeededMutex.RLock()
	defer fake.triggerIfNeededMutex.RUnlock()
	return len(fake.triggerIfNeededArgsForCall)
}

func (fake *RestartManager) TriggerIfNeededCalls(stub func(restart.Instance) error) {
	fake.triggerIfNeededMutex.Lock()
	defer fake.triggerIfNeededMutex.Unlock()
	fake.TriggerIfNeededStub = stub
}

func (fake *RestartManager) TriggerIfNeededArgsForCall(i int) restart.Instance {
	fake.triggerIfNeededMutex.RLock()
	defer fake.triggerIfNeededMutex.RUnlock()
	argsForCall := fake.triggerIfNeededArgsForCall[i]
	return argsForCall.arg1
}

func (fake *RestartManager) TriggerIfNeededReturns(result1 error) {
	fake.triggerIfNeededMutex.Lock()
	defer fake.triggerIfNeededMutex.Unlock()
	fake.TriggerIfNeededStub = nil
	fake.triggerIfNeededReturns = struct {
		result1 error
	}{result1}
}

func (fake *RestartManager) TriggerIfNeededReturnsOnCall(i int, result1 error) {
	fake.triggerIfNeededMutex.Lock()
	defer fake.triggerIfNeededMutex.Unlock()
	fake.TriggerIfNeededStub = nil
	if fake.triggerIfNeededReturnsOnCall == nil {
		fake.triggerIfNeededReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.triggerIfNeededReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *RestartManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.forConfigMapUpdateMutex.RLock()
	defer fake.forConfigMapUpdateMutex.RUnlock()
	fake.forRestartActionMutex.RLock()
	defer fake.forRestartActionMutex.RUnlock()
	fake.triggerIfNeededMutex.RLock()
	defer fake.triggerIfNeededMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *RestartManager) recordInvocation(key string, args []interface{}) {
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

var _ baseconsole.RestartManager = new(RestartManager)
