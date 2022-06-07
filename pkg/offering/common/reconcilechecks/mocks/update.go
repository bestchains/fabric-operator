// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common/reconcilechecks"
)

type Update struct {
	FabricVersionUpdatedStub        func() bool
	fabricVersionUpdatedMutex       sync.RWMutex
	fabricVersionUpdatedArgsForCall []struct {
	}
	fabricVersionUpdatedReturns struct {
		result1 bool
	}
	fabricVersionUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	ImagesUpdatedStub        func() bool
	imagesUpdatedMutex       sync.RWMutex
	imagesUpdatedArgsForCall []struct {
	}
	imagesUpdatedReturns struct {
		result1 bool
	}
	imagesUpdatedReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Update) FabricVersionUpdated() bool {
	fake.fabricVersionUpdatedMutex.Lock()
	ret, specificReturn := fake.fabricVersionUpdatedReturnsOnCall[len(fake.fabricVersionUpdatedArgsForCall)]
	fake.fabricVersionUpdatedArgsForCall = append(fake.fabricVersionUpdatedArgsForCall, struct {
	}{})
	stub := fake.FabricVersionUpdatedStub
	fakeReturns := fake.fabricVersionUpdatedReturns
	fake.recordInvocation("FabricVersionUpdated", []interface{}{})
	fake.fabricVersionUpdatedMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Update) FabricVersionUpdatedCallCount() int {
	fake.fabricVersionUpdatedMutex.RLock()
	defer fake.fabricVersionUpdatedMutex.RUnlock()
	return len(fake.fabricVersionUpdatedArgsForCall)
}

func (fake *Update) FabricVersionUpdatedCalls(stub func() bool) {
	fake.fabricVersionUpdatedMutex.Lock()
	defer fake.fabricVersionUpdatedMutex.Unlock()
	fake.FabricVersionUpdatedStub = stub
}

func (fake *Update) FabricVersionUpdatedReturns(result1 bool) {
	fake.fabricVersionUpdatedMutex.Lock()
	defer fake.fabricVersionUpdatedMutex.Unlock()
	fake.FabricVersionUpdatedStub = nil
	fake.fabricVersionUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) FabricVersionUpdatedReturnsOnCall(i int, result1 bool) {
	fake.fabricVersionUpdatedMutex.Lock()
	defer fake.fabricVersionUpdatedMutex.Unlock()
	fake.FabricVersionUpdatedStub = nil
	if fake.fabricVersionUpdatedReturnsOnCall == nil {
		fake.fabricVersionUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.fabricVersionUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) ImagesUpdated() bool {
	fake.imagesUpdatedMutex.Lock()
	ret, specificReturn := fake.imagesUpdatedReturnsOnCall[len(fake.imagesUpdatedArgsForCall)]
	fake.imagesUpdatedArgsForCall = append(fake.imagesUpdatedArgsForCall, struct {
	}{})
	stub := fake.ImagesUpdatedStub
	fakeReturns := fake.imagesUpdatedReturns
	fake.recordInvocation("ImagesUpdated", []interface{}{})
	fake.imagesUpdatedMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Update) ImagesUpdatedCallCount() int {
	fake.imagesUpdatedMutex.RLock()
	defer fake.imagesUpdatedMutex.RUnlock()
	return len(fake.imagesUpdatedArgsForCall)
}

func (fake *Update) ImagesUpdatedCalls(stub func() bool) {
	fake.imagesUpdatedMutex.Lock()
	defer fake.imagesUpdatedMutex.Unlock()
	fake.ImagesUpdatedStub = stub
}

func (fake *Update) ImagesUpdatedReturns(result1 bool) {
	fake.imagesUpdatedMutex.Lock()
	defer fake.imagesUpdatedMutex.Unlock()
	fake.ImagesUpdatedStub = nil
	fake.imagesUpdatedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Update) ImagesUpdatedReturnsOnCall(i int, result1 bool) {
	fake.imagesUpdatedMutex.Lock()
	defer fake.imagesUpdatedMutex.Unlock()
	fake.ImagesUpdatedStub = nil
	if fake.imagesUpdatedReturnsOnCall == nil {
		fake.imagesUpdatedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.imagesUpdatedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Update) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.fabricVersionUpdatedMutex.RLock()
	defer fake.fabricVersionUpdatedMutex.RUnlock()
	fake.imagesUpdatedMutex.RLock()
	defer fake.imagesUpdatedMutex.RUnlock()
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

var _ reconcilechecks.Update = new(Update)
