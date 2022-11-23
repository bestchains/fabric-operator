// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common/config"
)

type Crypto struct {
	GetCryptoStub        func() (*config.Response, error)
	getCryptoMutex       sync.RWMutex
	getCryptoArgsForCall []struct {
	}
	getCryptoReturns struct {
		result1 *config.Response
		result2 error
	}
	getCryptoReturnsOnCall map[int]struct {
		result1 *config.Response
		result2 error
	}
	PingCAStub        func() error
	pingCAMutex       sync.RWMutex
	pingCAArgsForCall []struct {
	}
	pingCAReturns struct {
		result1 error
	}
	pingCAReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateStub        func() error
	validateMutex       sync.RWMutex
	validateArgsForCall []struct {
	}
	validateReturns struct {
		result1 error
	}
	validateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Crypto) GetCrypto() (*config.Response, error) {
	fake.getCryptoMutex.Lock()
	ret, specificReturn := fake.getCryptoReturnsOnCall[len(fake.getCryptoArgsForCall)]
	fake.getCryptoArgsForCall = append(fake.getCryptoArgsForCall, struct {
	}{})
	fake.recordInvocation("GetCrypto", []interface{}{})
	fake.getCryptoMutex.Unlock()
	if fake.GetCryptoStub != nil {
		return fake.GetCryptoStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getCryptoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Crypto) GetCryptoCallCount() int {
	fake.getCryptoMutex.RLock()
	defer fake.getCryptoMutex.RUnlock()
	return len(fake.getCryptoArgsForCall)
}

func (fake *Crypto) GetCryptoCalls(stub func() (*config.Response, error)) {
	fake.getCryptoMutex.Lock()
	defer fake.getCryptoMutex.Unlock()
	fake.GetCryptoStub = stub
}

func (fake *Crypto) GetCryptoReturns(result1 *config.Response, result2 error) {
	fake.getCryptoMutex.Lock()
	defer fake.getCryptoMutex.Unlock()
	fake.GetCryptoStub = nil
	fake.getCryptoReturns = struct {
		result1 *config.Response
		result2 error
	}{result1, result2}
}

func (fake *Crypto) GetCryptoReturnsOnCall(i int, result1 *config.Response, result2 error) {
	fake.getCryptoMutex.Lock()
	defer fake.getCryptoMutex.Unlock()
	fake.GetCryptoStub = nil
	if fake.getCryptoReturnsOnCall == nil {
		fake.getCryptoReturnsOnCall = make(map[int]struct {
			result1 *config.Response
			result2 error
		})
	}
	fake.getCryptoReturnsOnCall[i] = struct {
		result1 *config.Response
		result2 error
	}{result1, result2}
}

func (fake *Crypto) PingCA() error {
	fake.pingCAMutex.Lock()
	ret, specificReturn := fake.pingCAReturnsOnCall[len(fake.pingCAArgsForCall)]
	fake.pingCAArgsForCall = append(fake.pingCAArgsForCall, struct {
	}{})
	fake.recordInvocation("PingCA", []interface{}{})
	fake.pingCAMutex.Unlock()
	if fake.PingCAStub != nil {
		return fake.PingCAStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.pingCAReturns
	return fakeReturns.result1
}

func (fake *Crypto) PingCACallCount() int {
	fake.pingCAMutex.RLock()
	defer fake.pingCAMutex.RUnlock()
	return len(fake.pingCAArgsForCall)
}

func (fake *Crypto) PingCACalls(stub func() error) {
	fake.pingCAMutex.Lock()
	defer fake.pingCAMutex.Unlock()
	fake.PingCAStub = stub
}

func (fake *Crypto) PingCAReturns(result1 error) {
	fake.pingCAMutex.Lock()
	defer fake.pingCAMutex.Unlock()
	fake.PingCAStub = nil
	fake.pingCAReturns = struct {
		result1 error
	}{result1}
}

func (fake *Crypto) PingCAReturnsOnCall(i int, result1 error) {
	fake.pingCAMutex.Lock()
	defer fake.pingCAMutex.Unlock()
	fake.PingCAStub = nil
	if fake.pingCAReturnsOnCall == nil {
		fake.pingCAReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pingCAReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Crypto) Validate() error {
	fake.validateMutex.Lock()
	ret, specificReturn := fake.validateReturnsOnCall[len(fake.validateArgsForCall)]
	fake.validateArgsForCall = append(fake.validateArgsForCall, struct {
	}{})
	fake.recordInvocation("Validate", []interface{}{})
	fake.validateMutex.Unlock()
	if fake.ValidateStub != nil {
		return fake.ValidateStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.validateReturns
	return fakeReturns.result1
}

func (fake *Crypto) ValidateCallCount() int {
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	return len(fake.validateArgsForCall)
}

func (fake *Crypto) ValidateCalls(stub func() error) {
	fake.validateMutex.Lock()
	defer fake.validateMutex.Unlock()
	fake.ValidateStub = stub
}

func (fake *Crypto) ValidateReturns(result1 error) {
	fake.validateMutex.Lock()
	defer fake.validateMutex.Unlock()
	fake.ValidateStub = nil
	fake.validateReturns = struct {
		result1 error
	}{result1}
}

func (fake *Crypto) ValidateReturnsOnCall(i int, result1 error) {
	fake.validateMutex.Lock()
	defer fake.validateMutex.Unlock()
	fake.ValidateStub = nil
	if fake.validateReturnsOnCall == nil {
		fake.validateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Crypto) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getCryptoMutex.RLock()
	defer fake.getCryptoMutex.RUnlock()
	fake.pingCAMutex.RLock()
	defer fake.pingCAMutex.RUnlock()
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Crypto) recordInvocation(key string, args []interface{}) {
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

var _ config.Crypto = new(Crypto)
