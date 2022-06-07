// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	baseconsole "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/console"
)

type Update struct {
	ConsoleCMUpdatedStub        func() bool
	consoleCMUpdatedMutex       sync.RWMutex
	consoleCMUpdatedArgsForCall []struct {
	}
	consoleCMUpdatedReturns struct {
		result1 bool
	}
	consoleCMUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	DeployerCMUpdatedStub        func() bool
	deployerCMUpdatedMutex       sync.RWMutex
	deployerCMUpdatedArgsForCall []struct {
	}
	deployerCMUpdatedReturns struct {
		result1 bool
	}
	deployerCMUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	EnvCMUpdatedStub        func() bool
	envCMUpdatedMutex       sync.RWMutex
	envCMUpdatedArgsForCall []struct {
	}
	envCMUpdatedReturns struct {
		result1 bool
	}
	envCMUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	RestartNeededStub        func() bool
	restartNeededMutex       sync.RWMutex
	restartNeededArgsForCall []struct {
	}
	restartNeededReturns struct {
		result1 bool
	}
	restartNeededReturnsOnCall map[int]struct {
		result1 bool
	}
	SpecUpdatedStub        func() bool
	specUpdatedMutex       sync.RWMutex
	specUpdatedArgsForCall []struct {
	}
	specUpdatedReturns struct {
		result1 bool
	}
	specUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Update) ConsoleCMUpdated() bool {
	fake.consoleCMUpdatedMutex.Lock()
	ret, specificReturn := fake.consoleCMUpdatedReturnsOnCall[len(fake.consoleCMUpdatedArgsForCall)]
	fake.consoleCMUpdatedArgsForCall = append(fake.consoleCMUpdatedArgsForCall, struct {
	}{})
	stub := fake.ConsoleCMUpdatedStub
	fakeReturns := fake.consoleCMUpdatedReturns
	fake.recordInvocation("ConsoleCMUpdated", []interface{}{})
	fake.consoleCMUpdatedMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Update) ConsoleCMUpdatedCallCount() int {
	fake.consoleCMUpdatedMutex.RLock()
	defer fake.consoleCMUpdatedMutex.RUnlock()
	return len(fake.consoleCMUpdatedArgsForCall)
}

func (fake *Update) ConsoleCMUpdatedCalls(stub func() bool) {
	fake.consoleCMUpdatedMutex.Lock()
	defer fake.consoleCMUpdatedMutex.Unlock()
	fake.ConsoleCMUpdatedStub = stub
}

func (fake *Update) ConsoleCMUpdatedReturns(result1 bool) {
	fake.consoleCMUpdatedMutex.Lock()
	defer fake.consoleCMUpdatedMutex.Unlock()
	fake.ConsoleCMUpdatedStub = nil
	fake.consoleCMUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) ConsoleCMUpdatedReturnsOnCall(i int, result1 bool) {
	fake.consoleCMUpdatedMutex.Lock()
	defer fake.consoleCMUpdatedMutex.Unlock()
	fake.ConsoleCMUpdatedStub = nil
	if fake.consoleCMUpdatedReturnsOnCall == nil {
		fake.consoleCMUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.consoleCMUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) DeployerCMUpdated() bool {
	fake.deployerCMUpdatedMutex.Lock()
	ret, specificReturn := fake.deployerCMUpdatedReturnsOnCall[len(fake.deployerCMUpdatedArgsForCall)]
	fake.deployerCMUpdatedArgsForCall = append(fake.deployerCMUpdatedArgsForCall, struct {
	}{})
	stub := fake.DeployerCMUpdatedStub
	fakeReturns := fake.deployerCMUpdatedReturns
	fake.recordInvocation("DeployerCMUpdated", []interface{}{})
	fake.deployerCMUpdatedMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Update) DeployerCMUpdatedCallCount() int {
	fake.deployerCMUpdatedMutex.RLock()
	defer fake.deployerCMUpdatedMutex.RUnlock()
	return len(fake.deployerCMUpdatedArgsForCall)
}

func (fake *Update) DeployerCMUpdatedCalls(stub func() bool) {
	fake.deployerCMUpdatedMutex.Lock()
	defer fake.deployerCMUpdatedMutex.Unlock()
	fake.DeployerCMUpdatedStub = stub
}

func (fake *Update) DeployerCMUpdatedReturns(result1 bool) {
	fake.deployerCMUpdatedMutex.Lock()
	defer fake.deployerCMUpdatedMutex.Unlock()
	fake.DeployerCMUpdatedStub = nil
	fake.deployerCMUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) DeployerCMUpdatedReturnsOnCall(i int, result1 bool) {
	fake.deployerCMUpdatedMutex.Lock()
	defer fake.deployerCMUpdatedMutex.Unlock()
	fake.DeployerCMUpdatedStub = nil
	if fake.deployerCMUpdatedReturnsOnCall == nil {
		fake.deployerCMUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.deployerCMUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) EnvCMUpdated() bool {
	fake.envCMUpdatedMutex.Lock()
	ret, specificReturn := fake.envCMUpdatedReturnsOnCall[len(fake.envCMUpdatedArgsForCall)]
	fake.envCMUpdatedArgsForCall = append(fake.envCMUpdatedArgsForCall, struct {
	}{})
	stub := fake.EnvCMUpdatedStub
	fakeReturns := fake.envCMUpdatedReturns
	fake.recordInvocation("EnvCMUpdated", []interface{}{})
	fake.envCMUpdatedMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Update) EnvCMUpdatedCallCount() int {
	fake.envCMUpdatedMutex.RLock()
	defer fake.envCMUpdatedMutex.RUnlock()
	return len(fake.envCMUpdatedArgsForCall)
}

func (fake *Update) EnvCMUpdatedCalls(stub func() bool) {
	fake.envCMUpdatedMutex.Lock()
	defer fake.envCMUpdatedMutex.Unlock()
	fake.EnvCMUpdatedStub = stub
}

func (fake *Update) EnvCMUpdatedReturns(result1 bool) {
	fake.envCMUpdatedMutex.Lock()
	defer fake.envCMUpdatedMutex.Unlock()
	fake.EnvCMUpdatedStub = nil
	fake.envCMUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) EnvCMUpdatedReturnsOnCall(i int, result1 bool) {
	fake.envCMUpdatedMutex.Lock()
	defer fake.envCMUpdatedMutex.Unlock()
	fake.EnvCMUpdatedStub = nil
	if fake.envCMUpdatedReturnsOnCall == nil {
		fake.envCMUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.envCMUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) RestartNeeded() bool {
	fake.restartNeededMutex.Lock()
	ret, specificReturn := fake.restartNeededReturnsOnCall[len(fake.restartNeededArgsForCall)]
	fake.restartNeededArgsForCall = append(fake.restartNeededArgsForCall, struct {
	}{})
	stub := fake.RestartNeededStub
	fakeReturns := fake.restartNeededReturns
	fake.recordInvocation("RestartNeeded", []interface{}{})
	fake.restartNeededMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Update) RestartNeededCallCount() int {
	fake.restartNeededMutex.RLock()
	defer fake.restartNeededMutex.RUnlock()
	return len(fake.restartNeededArgsForCall)
}

func (fake *Update) RestartNeededCalls(stub func() bool) {
	fake.restartNeededMutex.Lock()
	defer fake.restartNeededMutex.Unlock()
	fake.RestartNeededStub = stub
}

func (fake *Update) RestartNeededReturns(result1 bool) {
	fake.restartNeededMutex.Lock()
	defer fake.restartNeededMutex.Unlock()
	fake.RestartNeededStub = nil
	fake.restartNeededReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) RestartNeededReturnsOnCall(i int, result1 bool) {
	fake.restartNeededMutex.Lock()
	defer fake.restartNeededMutex.Unlock()
	fake.RestartNeededStub = nil
	if fake.restartNeededReturnsOnCall == nil {
		fake.restartNeededReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.restartNeededReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) SpecUpdated() bool {
	fake.specUpdatedMutex.Lock()
	ret, specificReturn := fake.specUpdatedReturnsOnCall[len(fake.specUpdatedArgsForCall)]
	fake.specUpdatedArgsForCall = append(fake.specUpdatedArgsForCall, struct {
	}{})
	stub := fake.SpecUpdatedStub
	fakeReturns := fake.specUpdatedReturns
	fake.recordInvocation("SpecUpdated", []interface{}{})
	fake.specUpdatedMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Update) SpecUpdatedCallCount() int {
	fake.specUpdatedMutex.RLock()
	defer fake.specUpdatedMutex.RUnlock()
	return len(fake.specUpdatedArgsForCall)
}

func (fake *Update) SpecUpdatedCalls(stub func() bool) {
	fake.specUpdatedMutex.Lock()
	defer fake.specUpdatedMutex.Unlock()
	fake.SpecUpdatedStub = stub
}

func (fake *Update) SpecUpdatedReturns(result1 bool) {
	fake.specUpdatedMutex.Lock()
	defer fake.specUpdatedMutex.Unlock()
	fake.SpecUpdatedStub = nil
	fake.specUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) SpecUpdatedReturnsOnCall(i int, result1 bool) {
	fake.specUpdatedMutex.Lock()
	defer fake.specUpdatedMutex.Unlock()
	fake.SpecUpdatedStub = nil
	if fake.specUpdatedReturnsOnCall == nil {
		fake.specUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.specUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.consoleCMUpdatedMutex.RLock()
	defer fake.consoleCMUpdatedMutex.RUnlock()
	fake.deployerCMUpdatedMutex.RLock()
	defer fake.deployerCMUpdatedMutex.RUnlock()
	fake.envCMUpdatedMutex.RLock()
	defer fake.envCMUpdatedMutex.RUnlock()
	fake.restartNeededMutex.RLock()
	defer fake.restartNeededMutex.RUnlock()
	fake.specUpdatedMutex.RLock()
	defer fake.specUpdatedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Update) recordInvocation(key string, args []interface{}) {
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

var _ baseconsole.Update = new(Update)
