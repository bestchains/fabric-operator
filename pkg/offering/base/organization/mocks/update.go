// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization"
)

type Update struct {
	AdminTransferedStub        func() string
	adminTransferedMutex       sync.RWMutex
	adminTransferedArgsForCall []struct {
	}
	adminTransferedReturns struct {
		result1 string
	}
	adminTransferedReturnsOnCall map[int]struct {
		result1 string
	}
	AdminUpdatedStub        func() bool
	adminUpdatedMutex       sync.RWMutex
	adminUpdatedArgsForCall []struct {
	}
	adminUpdatedReturns struct {
		result1 bool
	}
	adminUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	ClientsRemovedStub        func() string
	clientsRemovedMutex       sync.RWMutex
	clientsRemovedArgsForCall []struct {
	}
	clientsRemovedReturns struct {
		result1 string
	}
	clientsRemovedReturnsOnCall map[int]struct {
		result1 string
	}
	ClientsUpdatedStub        func() bool
	clientsUpdatedMutex       sync.RWMutex
	clientsUpdatedArgsForCall []struct {
	}
	clientsUpdatedReturns struct {
		result1 bool
	}
	clientsUpdatedReturnsOnCall map[int]struct {
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
	TokenUpdatedStub        func() bool
	tokenUpdatedMutex       sync.RWMutex
	tokenUpdatedArgsForCall []struct {
	}
	tokenUpdatedReturns struct {
		result1 bool
	}
	tokenUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Update) AdminTransfered() string {
	fake.adminTransferedMutex.Lock()
	ret, specificReturn := fake.adminTransferedReturnsOnCall[len(fake.adminTransferedArgsForCall)]
	fake.adminTransferedArgsForCall = append(fake.adminTransferedArgsForCall, struct {
	}{})
	fake.recordInvocation("AdminTransfered", []interface{}{})
	fake.adminTransferedMutex.Unlock()
	if fake.AdminTransferedStub != nil {
		return fake.AdminTransferedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.adminTransferedReturns
	return fakeReturns.result1
}

func (fake *Update) AdminTransferedCallCount() int {
	fake.adminTransferedMutex.RLock()
	defer fake.adminTransferedMutex.RUnlock()
	return len(fake.adminTransferedArgsForCall)
}

func (fake *Update) AdminTransferedCalls(stub func() string) {
	fake.adminTransferedMutex.Lock()
	defer fake.adminTransferedMutex.Unlock()
	fake.AdminTransferedStub = stub
}

func (fake *Update) AdminTransferedReturns(result1 string) {
	fake.adminTransferedMutex.Lock()
	defer fake.adminTransferedMutex.Unlock()
	fake.AdminTransferedStub = nil
	fake.adminTransferedReturns = struct {
		result1 string
	}{result1}
}

func (fake *Update) AdminTransferedReturnsOnCall(i int, result1 string) {
	fake.adminTransferedMutex.Lock()
	defer fake.adminTransferedMutex.Unlock()
	fake.AdminTransferedStub = nil
	if fake.adminTransferedReturnsOnCall == nil {
		fake.adminTransferedReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.adminTransferedReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *Update) AdminUpdated() bool {
	fake.adminUpdatedMutex.Lock()
	ret, specificReturn := fake.adminUpdatedReturnsOnCall[len(fake.adminUpdatedArgsForCall)]
	fake.adminUpdatedArgsForCall = append(fake.adminUpdatedArgsForCall, struct {
	}{})
	fake.recordInvocation("AdminUpdated", []interface{}{})
	fake.adminUpdatedMutex.Unlock()
	if fake.AdminUpdatedStub != nil {
		return fake.AdminUpdatedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.adminUpdatedReturns
	return fakeReturns.result1
}

func (fake *Update) AdminUpdatedCallCount() int {
	fake.adminUpdatedMutex.RLock()
	defer fake.adminUpdatedMutex.RUnlock()
	return len(fake.adminUpdatedArgsForCall)
}

func (fake *Update) AdminUpdatedCalls(stub func() bool) {
	fake.adminUpdatedMutex.Lock()
	defer fake.adminUpdatedMutex.Unlock()
	fake.AdminUpdatedStub = stub
}

func (fake *Update) AdminUpdatedReturns(result1 bool) {
	fake.adminUpdatedMutex.Lock()
	defer fake.adminUpdatedMutex.Unlock()
	fake.AdminUpdatedStub = nil
	fake.adminUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) AdminUpdatedReturnsOnCall(i int, result1 bool) {
	fake.adminUpdatedMutex.Lock()
	defer fake.adminUpdatedMutex.Unlock()
	fake.AdminUpdatedStub = nil
	if fake.adminUpdatedReturnsOnCall == nil {
		fake.adminUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.adminUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) ClientsRemoved() string {
	fake.clientsRemovedMutex.Lock()
	ret, specificReturn := fake.clientsRemovedReturnsOnCall[len(fake.clientsRemovedArgsForCall)]
	fake.clientsRemovedArgsForCall = append(fake.clientsRemovedArgsForCall, struct {
	}{})
	fake.recordInvocation("ClientsRemoved", []interface{}{})
	fake.clientsRemovedMutex.Unlock()
	if fake.ClientsRemovedStub != nil {
		return fake.ClientsRemovedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.clientsRemovedReturns
	return fakeReturns.result1
}

func (fake *Update) ClientsRemovedCallCount() int {
	fake.clientsRemovedMutex.RLock()
	defer fake.clientsRemovedMutex.RUnlock()
	return len(fake.clientsRemovedArgsForCall)
}

func (fake *Update) ClientsRemovedCalls(stub func() string) {
	fake.clientsRemovedMutex.Lock()
	defer fake.clientsRemovedMutex.Unlock()
	fake.ClientsRemovedStub = stub
}

func (fake *Update) ClientsRemovedReturns(result1 string) {
	fake.clientsRemovedMutex.Lock()
	defer fake.clientsRemovedMutex.Unlock()
	fake.ClientsRemovedStub = nil
	fake.clientsRemovedReturns = struct {
		result1 string
	}{result1}
}

func (fake *Update) ClientsRemovedReturnsOnCall(i int, result1 string) {
	fake.clientsRemovedMutex.Lock()
	defer fake.clientsRemovedMutex.Unlock()
	fake.ClientsRemovedStub = nil
	if fake.clientsRemovedReturnsOnCall == nil {
		fake.clientsRemovedReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.clientsRemovedReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *Update) ClientsUpdated() bool {
	fake.clientsUpdatedMutex.Lock()
	ret, specificReturn := fake.clientsUpdatedReturnsOnCall[len(fake.clientsUpdatedArgsForCall)]
	fake.clientsUpdatedArgsForCall = append(fake.clientsUpdatedArgsForCall, struct {
	}{})
	fake.recordInvocation("ClientsUpdated", []interface{}{})
	fake.clientsUpdatedMutex.Unlock()
	if fake.ClientsUpdatedStub != nil {
		return fake.ClientsUpdatedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.clientsUpdatedReturns
	return fakeReturns.result1
}

func (fake *Update) ClientsUpdatedCallCount() int {
	fake.clientsUpdatedMutex.RLock()
	defer fake.clientsUpdatedMutex.RUnlock()
	return len(fake.clientsUpdatedArgsForCall)
}

func (fake *Update) ClientsUpdatedCalls(stub func() bool) {
	fake.clientsUpdatedMutex.Lock()
	defer fake.clientsUpdatedMutex.Unlock()
	fake.ClientsUpdatedStub = stub
}

func (fake *Update) ClientsUpdatedReturns(result1 bool) {
	fake.clientsUpdatedMutex.Lock()
	defer fake.clientsUpdatedMutex.Unlock()
	fake.ClientsUpdatedStub = nil
	fake.clientsUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) ClientsUpdatedReturnsOnCall(i int, result1 bool) {
	fake.clientsUpdatedMutex.Lock()
	defer fake.clientsUpdatedMutex.Unlock()
	fake.ClientsUpdatedStub = nil
	if fake.clientsUpdatedReturnsOnCall == nil {
		fake.clientsUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.clientsUpdatedReturnsOnCall[i] = struct {
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

func (fake *Update) TokenUpdated() bool {
	fake.tokenUpdatedMutex.Lock()
	ret, specificReturn := fake.tokenUpdatedReturnsOnCall[len(fake.tokenUpdatedArgsForCall)]
	fake.tokenUpdatedArgsForCall = append(fake.tokenUpdatedArgsForCall, struct {
	}{})
	fake.recordInvocation("TokenUpdated", []interface{}{})
	fake.tokenUpdatedMutex.Unlock()
	if fake.TokenUpdatedStub != nil {
		return fake.TokenUpdatedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.tokenUpdatedReturns
	return fakeReturns.result1
}

func (fake *Update) TokenUpdatedCallCount() int {
	fake.tokenUpdatedMutex.RLock()
	defer fake.tokenUpdatedMutex.RUnlock()
	return len(fake.tokenUpdatedArgsForCall)
}

func (fake *Update) TokenUpdatedCalls(stub func() bool) {
	fake.tokenUpdatedMutex.Lock()
	defer fake.tokenUpdatedMutex.Unlock()
	fake.TokenUpdatedStub = stub
}

func (fake *Update) TokenUpdatedReturns(result1 bool) {
	fake.tokenUpdatedMutex.Lock()
	defer fake.tokenUpdatedMutex.Unlock()
	fake.TokenUpdatedStub = nil
	fake.tokenUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) TokenUpdatedReturnsOnCall(i int, result1 bool) {
	fake.tokenUpdatedMutex.Lock()
	defer fake.tokenUpdatedMutex.Unlock()
	fake.TokenUpdatedStub = nil
	if fake.tokenUpdatedReturnsOnCall == nil {
		fake.tokenUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.tokenUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.adminTransferedMutex.RLock()
	defer fake.adminTransferedMutex.RUnlock()
	fake.adminUpdatedMutex.RLock()
	defer fake.adminUpdatedMutex.RUnlock()
	fake.clientsRemovedMutex.RLock()
	defer fake.clientsRemovedMutex.RUnlock()
	fake.clientsUpdatedMutex.RLock()
	defer fake.clientsUpdatedMutex.RUnlock()
	fake.specUpdatedMutex.RLock()
	defer fake.specUpdatedMutex.RUnlock()
	fake.tokenUpdatedMutex.RLock()
	defer fake.tokenUpdatedMutex.RUnlock()
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
