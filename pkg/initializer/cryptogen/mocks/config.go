// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/apis/common"
	"github.com/IBM-Blockchain/fabric-operator/pkg/initializer/cryptogen"
)

type Config struct {
	GetBCCSPSectionStub        func() *common.BCCSP
	getBCCSPSectionMutex       sync.RWMutex
	getBCCSPSectionArgsForCall []struct {
	}
	getBCCSPSectionReturns struct {
		result1 *common.BCCSP
	}
	getBCCSPSectionReturnsOnCall map[int]struct {
		result1 *common.BCCSP
	}
	SetDefaultKeyStoreStub        func()
	setDefaultKeyStoreMutex       sync.RWMutex
	setDefaultKeyStoreArgsForCall []struct {
	}
	SetPKCS11DefaultsStub        func(bool)
	setPKCS11DefaultsMutex       sync.RWMutex
	setPKCS11DefaultsArgsForCall []struct {
		arg1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Config) GetBCCSPSection() *common.BCCSP {
	fake.getBCCSPSectionMutex.Lock()
	ret, specificReturn := fake.getBCCSPSectionReturnsOnCall[len(fake.getBCCSPSectionArgsForCall)]
	fake.getBCCSPSectionArgsForCall = append(fake.getBCCSPSectionArgsForCall, struct {
	}{})
	fake.recordInvocation("GetBCCSPSection", []interface{}{})
	fake.getBCCSPSectionMutex.Unlock()
	if fake.GetBCCSPSectionStub != nil {
		return fake.GetBCCSPSectionStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getBCCSPSectionReturns
	return fakeReturns.result1
}

func (fake *Config) GetBCCSPSectionCallCount() int {
	fake.getBCCSPSectionMutex.RLock()
	defer fake.getBCCSPSectionMutex.RUnlock()
	return len(fake.getBCCSPSectionArgsForCall)
}

func (fake *Config) GetBCCSPSectionCalls(stub func() *common.BCCSP) {
	fake.getBCCSPSectionMutex.Lock()
	defer fake.getBCCSPSectionMutex.Unlock()
	fake.GetBCCSPSectionStub = stub
}

func (fake *Config) GetBCCSPSectionReturns(result1 *common.BCCSP) {
	fake.getBCCSPSectionMutex.Lock()
	defer fake.getBCCSPSectionMutex.Unlock()
	fake.GetBCCSPSectionStub = nil
	fake.getBCCSPSectionReturns = struct {
		result1 *common.BCCSP
	}{result1}
}

func (fake *Config) GetBCCSPSectionReturnsOnCall(i int, result1 *common.BCCSP) {
	fake.getBCCSPSectionMutex.Lock()
	defer fake.getBCCSPSectionMutex.Unlock()
	fake.GetBCCSPSectionStub = nil
	if fake.getBCCSPSectionReturnsOnCall == nil {
		fake.getBCCSPSectionReturnsOnCall = make(map[int]struct {
			result1 *common.BCCSP
		})
	}
	fake.getBCCSPSectionReturnsOnCall[i] = struct {
		result1 *common.BCCSP
	}{result1}
}

func (fake *Config) SetDefaultKeyStore() {
	fake.setDefaultKeyStoreMutex.Lock()
	fake.setDefaultKeyStoreArgsForCall = append(fake.setDefaultKeyStoreArgsForCall, struct {
	}{})
	fake.recordInvocation("SetDefaultKeyStore", []interface{}{})
	fake.setDefaultKeyStoreMutex.Unlock()
	if fake.SetDefaultKeyStoreStub != nil {
		fake.SetDefaultKeyStoreStub()
	}
}

func (fake *Config) SetDefaultKeyStoreCallCount() int {
	fake.setDefaultKeyStoreMutex.RLock()
	defer fake.setDefaultKeyStoreMutex.RUnlock()
	return len(fake.setDefaultKeyStoreArgsForCall)
}

func (fake *Config) SetDefaultKeyStoreCalls(stub func()) {
	fake.setDefaultKeyStoreMutex.Lock()
	defer fake.setDefaultKeyStoreMutex.Unlock()
	fake.SetDefaultKeyStoreStub = stub
}

func (fake *Config) SetPKCS11Defaults(arg1 bool) {
	fake.setPKCS11DefaultsMutex.Lock()
	fake.setPKCS11DefaultsArgsForCall = append(fake.setPKCS11DefaultsArgsForCall, struct {
		arg1 bool
	}{arg1})
	fake.recordInvocation("SetPKCS11Defaults", []interface{}{arg1})
	fake.setPKCS11DefaultsMutex.Unlock()
	if fake.SetPKCS11DefaultsStub != nil {
		fake.SetPKCS11DefaultsStub(arg1)
	}
}

func (fake *Config) SetPKCS11DefaultsCallCount() int {
	fake.setPKCS11DefaultsMutex.RLock()
	defer fake.setPKCS11DefaultsMutex.RUnlock()
	return len(fake.setPKCS11DefaultsArgsForCall)
}

func (fake *Config) SetPKCS11DefaultsCalls(stub func(bool)) {
	fake.setPKCS11DefaultsMutex.Lock()
	defer fake.setPKCS11DefaultsMutex.Unlock()
	fake.SetPKCS11DefaultsStub = stub
}

func (fake *Config) SetPKCS11DefaultsArgsForCall(i int) bool {
	fake.setPKCS11DefaultsMutex.RLock()
	defer fake.setPKCS11DefaultsMutex.RUnlock()
	argsForCall := fake.setPKCS11DefaultsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Config) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getBCCSPSectionMutex.RLock()
	defer fake.getBCCSPSectionMutex.RUnlock()
	fake.setDefaultKeyStoreMutex.RLock()
	defer fake.setDefaultKeyStoreMutex.RUnlock()
	fake.setPKCS11DefaultsMutex.RLock()
	defer fake.setPKCS11DefaultsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Config) recordInvocation(key string, args []interface{}) {
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

var _ cryptogen.Config = new(Config)
