// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization"
)

type Update struct {
	AdminOrCAUpdatedStub        func() bool
	adminOrCAUpdatedMutex       sync.RWMutex
	adminOrCAUpdatedArgsForCall []struct {
	}
	adminOrCAUpdatedReturns struct {
		result1 bool
	}
	adminOrCAUpdatedReturnsOnCall map[int]struct {
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

func (fake *Update) AdminOrCAUpdated() bool {
	fake.adminOrCAUpdatedMutex.Lock()
	ret, specificReturn := fake.adminOrCAUpdatedReturnsOnCall[len(fake.adminOrCAUpdatedArgsForCall)]
	fake.adminOrCAUpdatedArgsForCall = append(fake.adminOrCAUpdatedArgsForCall, struct {
	}{})
	fake.recordInvocation("AdminOrCAUpdated", []interface{}{})
	fake.adminOrCAUpdatedMutex.Unlock()
	if fake.AdminOrCAUpdatedStub != nil {
		return fake.AdminOrCAUpdatedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.adminOrCAUpdatedReturns
	return fakeReturns.result1
}

func (fake *Update) AdminOrCAUpdatedCallCount() int {
	fake.adminOrCAUpdatedMutex.RLock()
	defer fake.adminOrCAUpdatedMutex.RUnlock()
	return len(fake.adminOrCAUpdatedArgsForCall)
}

func (fake *Update) AdminOrCAUpdatedCalls(stub func() bool) {
	fake.adminOrCAUpdatedMutex.Lock()
	defer fake.adminOrCAUpdatedMutex.Unlock()
	fake.AdminOrCAUpdatedStub = stub
}

func (fake *Update) AdminOrCAUpdatedReturns(result1 bool) {
	fake.adminOrCAUpdatedMutex.Lock()
	defer fake.adminOrCAUpdatedMutex.Unlock()
	fake.AdminOrCAUpdatedStub = nil
	fake.adminOrCAUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) AdminOrCAUpdatedReturnsOnCall(i int, result1 bool) {
	fake.adminOrCAUpdatedMutex.Lock()
	defer fake.adminOrCAUpdatedMutex.Unlock()
	fake.AdminOrCAUpdatedStub = nil
	if fake.adminOrCAUpdatedReturnsOnCall == nil {
		fake.adminOrCAUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.adminOrCAUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) SpecUpdated() bool {
	fake.specUpdatedMutex.Lock()
	ret, specificReturn := fake.specUpdatedReturnsOnCall[len(fake.specUpdatedArgsForCall)]
	fake.specUpdatedArgsForCall = append(fake.specUpdatedArgsForCall, struct {
	}{})
	fake.recordInvocation("SpecUpdated", []interface{}{})
	fake.specUpdatedMutex.Unlock()
	if fake.SpecUpdatedStub != nil {
		return fake.SpecUpdatedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.specUpdatedReturns
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
	fake.adminOrCAUpdatedMutex.RLock()
	defer fake.adminOrCAUpdatedMutex.RUnlock()
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

var _ organization.Update = new(Update)
