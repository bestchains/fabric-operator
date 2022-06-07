// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/common/reconcilechecks/images"
)

type Instance struct {
	GetArchStub        func() []string
	getArchMutex       sync.RWMutex
	getArchArgsForCall []struct {
	}
	getArchReturns struct {
		result1 []string
	}
	getArchReturnsOnCall map[int]struct {
		result1 []string
	}
	GetFabricVersionStub        func() string
	getFabricVersionMutex       sync.RWMutex
	getFabricVersionArgsForCall []struct {
	}
	getFabricVersionReturns struct {
		result1 string
	}
	getFabricVersionReturnsOnCall map[int]struct {
		result1 string
	}
	GetRegistryURLStub        func() string
	getRegistryURLMutex       sync.RWMutex
	getRegistryURLArgsForCall []struct {
	}
	getRegistryURLReturns struct {
		result1 string
	}
	getRegistryURLReturnsOnCall map[int]struct {
		result1 string
	}
	ImagesSetStub        func() bool
	imagesSetMutex       sync.RWMutex
	imagesSetArgsForCall []struct {
	}
	imagesSetReturns struct {
		result1 bool
	}
	imagesSetReturnsOnCall map[int]struct {
		result1 bool
	}
	SetFabricVersionStub        func(string)
	setFabricVersionMutex       sync.RWMutex
	setFabricVersionArgsForCall []struct {
		arg1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Instance) GetArch() []string {
	fake.getArchMutex.Lock()
	ret, specificReturn := fake.getArchReturnsOnCall[len(fake.getArchArgsForCall)]
	fake.getArchArgsForCall = append(fake.getArchArgsForCall, struct {
	}{})
	stub := fake.GetArchStub
	fakeReturns := fake.getArchReturns
	fake.recordInvocation("GetArch", []interface{}{})
	fake.getArchMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Instance) GetArchCallCount() int {
	fake.getArchMutex.RLock()
	defer fake.getArchMutex.RUnlock()
	return len(fake.getArchArgsForCall)
}

func (fake *Instance) GetArchCalls(stub func() []string) {
	fake.getArchMutex.Lock()
	defer fake.getArchMutex.Unlock()
	fake.GetArchStub = stub
}

func (fake *Instance) GetArchReturns(result1 []string) {
	fake.getArchMutex.Lock()
	defer fake.getArchMutex.Unlock()
	fake.GetArchStub = nil
	fake.getArchReturns = struct {
		result1 []string
	}{result1}
}

func (fake *Instance) GetArchReturnsOnCall(i int, result1 []string) {
	fake.getArchMutex.Lock()
	defer fake.getArchMutex.Unlock()
	fake.GetArchStub = nil
	if fake.getArchReturnsOnCall == nil {
		fake.getArchReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.getArchReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *Instance) GetFabricVersion() string {
	fake.getFabricVersionMutex.Lock()
	ret, specificReturn := fake.getFabricVersionReturnsOnCall[len(fake.getFabricVersionArgsForCall)]
	fake.getFabricVersionArgsForCall = append(fake.getFabricVersionArgsForCall, struct {
	}{})
	stub := fake.GetFabricVersionStub
	fakeReturns := fake.getFabricVersionReturns
	fake.recordInvocation("GetFabricVersion", []interface{}{})
	fake.getFabricVersionMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Instance) GetFabricVersionCallCount() int {
	fake.getFabricVersionMutex.RLock()
	defer fake.getFabricVersionMutex.RUnlock()
	return len(fake.getFabricVersionArgsForCall)
}

func (fake *Instance) GetFabricVersionCalls(stub func() string) {
	fake.getFabricVersionMutex.Lock()
	defer fake.getFabricVersionMutex.Unlock()
	fake.GetFabricVersionStub = stub
}

func (fake *Instance) GetFabricVersionReturns(result1 string) {
	fake.getFabricVersionMutex.Lock()
	defer fake.getFabricVersionMutex.Unlock()
	fake.GetFabricVersionStub = nil
	fake.getFabricVersionReturns = struct {
		result1 string
	}{result1}
}

func (fake *Instance) GetFabricVersionReturnsOnCall(i int, result1 string) {
	fake.getFabricVersionMutex.Lock()
	defer fake.getFabricVersionMutex.Unlock()
	fake.GetFabricVersionStub = nil
	if fake.getFabricVersionReturnsOnCall == nil {
		fake.getFabricVersionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getFabricVersionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *Instance) GetRegistryURL() string {
	fake.getRegistryURLMutex.Lock()
	ret, specificReturn := fake.getRegistryURLReturnsOnCall[len(fake.getRegistryURLArgsForCall)]
	fake.getRegistryURLArgsForCall = append(fake.getRegistryURLArgsForCall, struct {
	}{})
	stub := fake.GetRegistryURLStub
	fakeReturns := fake.getRegistryURLReturns
	fake.recordInvocation("GetRegistryURL", []interface{}{})
	fake.getRegistryURLMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Instance) GetRegistryURLCallCount() int {
	fake.getRegistryURLMutex.RLock()
	defer fake.getRegistryURLMutex.RUnlock()
	return len(fake.getRegistryURLArgsForCall)
}

func (fake *Instance) GetRegistryURLCalls(stub func() string) {
	fake.getRegistryURLMutex.Lock()
	defer fake.getRegistryURLMutex.Unlock()
	fake.GetRegistryURLStub = stub
}

func (fake *Instance) GetRegistryURLReturns(result1 string) {
	fake.getRegistryURLMutex.Lock()
	defer fake.getRegistryURLMutex.Unlock()
	fake.GetRegistryURLStub = nil
	fake.getRegistryURLReturns = struct {
		result1 string
	}{result1}
}

func (fake *Instance) GetRegistryURLReturnsOnCall(i int, result1 string) {
	fake.getRegistryURLMutex.Lock()
	defer fake.getRegistryURLMutex.Unlock()
	fake.GetRegistryURLStub = nil
	if fake.getRegistryURLReturnsOnCall == nil {
		fake.getRegistryURLReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getRegistryURLReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *Instance) ImagesSet() bool {
	fake.imagesSetMutex.Lock()
	ret, specificReturn := fake.imagesSetReturnsOnCall[len(fake.imagesSetArgsForCall)]
	fake.imagesSetArgsForCall = append(fake.imagesSetArgsForCall, struct {
	}{})
	stub := fake.ImagesSetStub
	fakeReturns := fake.imagesSetReturns
	fake.recordInvocation("ImagesSet", []interface{}{})
	fake.imagesSetMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *Instance) ImagesSetCallCount() int {
	fake.imagesSetMutex.RLock()
	defer fake.imagesSetMutex.RUnlock()
	return len(fake.imagesSetArgsForCall)
}

func (fake *Instance) ImagesSetCalls(stub func() bool) {
	fake.imagesSetMutex.Lock()
	defer fake.imagesSetMutex.Unlock()
	fake.ImagesSetStub = stub
}

func (fake *Instance) ImagesSetReturns(result1 bool) {
	fake.imagesSetMutex.Lock()
	defer fake.imagesSetMutex.Unlock()
	fake.ImagesSetStub = nil
	fake.imagesSetReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Instance) ImagesSetReturnsOnCall(i int, result1 bool) {
	fake.imagesSetMutex.Lock()
	defer fake.imagesSetMutex.Unlock()
	fake.ImagesSetStub = nil
	if fake.imagesSetReturnsOnCall == nil {
		fake.imagesSetReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.imagesSetReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Instance) SetFabricVersion(arg1 string) {
	fake.setFabricVersionMutex.Lock()
	fake.setFabricVersionArgsForCall = append(fake.setFabricVersionArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.SetFabricVersionStub
	fake.recordInvocation("SetFabricVersion", []interface{}{arg1})
	fake.setFabricVersionMutex.Unlock()
	if stub != nil {
		fake.SetFabricVersionStub(arg1)
	}
}

func (fake *Instance) SetFabricVersionCallCount() int {
	fake.setFabricVersionMutex.RLock()
	defer fake.setFabricVersionMutex.RUnlock()
	return len(fake.setFabricVersionArgsForCall)
}

func (fake *Instance) SetFabricVersionCalls(stub func(string)) {
	fake.setFabricVersionMutex.Lock()
	defer fake.setFabricVersionMutex.Unlock()
	fake.SetFabricVersionStub = stub
}

func (fake *Instance) SetFabricVersionArgsForCall(i int) string {
	fake.setFabricVersionMutex.RLock()
	defer fake.setFabricVersionMutex.RUnlock()
	argsForCall := fake.setFabricVersionArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Instance) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getArchMutex.RLock()
	defer fake.getArchMutex.RUnlock()
	fake.getFabricVersionMutex.RLock()
	defer fake.getFabricVersionMutex.RUnlock()
	fake.getRegistryURLMutex.RLock()
	defer fake.getRegistryURLMutex.RUnlock()
	fake.imagesSetMutex.RLock()
	defer fake.imagesSetMutex.RUnlock()
	fake.setFabricVersionMutex.RLock()
	defer fake.setFabricVersionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Instance) recordInvocation(key string, args []interface{}) {
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

var _ images.Instance = new(Instance)
