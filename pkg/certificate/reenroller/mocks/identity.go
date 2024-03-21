// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/certificate/reenroller"
	"github.com/hyperledger/fabric-ca/api"
	"github.com/hyperledger/fabric-ca/lib"
	"github.com/hyperledger/fabric-ca/lib/client/credential/x509"
)

type Identity struct {
	GetClientStub        func() *lib.Client
	getClientMutex       sync.RWMutex
	getClientArgsForCall []struct {
	}
	getClientReturns struct {
		result1 *lib.Client
	}
	getClientReturnsOnCall map[int]struct {
		result1 *lib.Client
	}
	GetECertStub        func() *x509.Signer
	getECertMutex       sync.RWMutex
	getECertArgsForCall []struct {
	}
	getECertReturns struct {
		result1 *x509.Signer
	}
	getECertReturnsOnCall map[int]struct {
		result1 *x509.Signer
	}
	ReenrollStub        func(*api.ReenrollmentRequest) (*lib.EnrollmentResponse, error)
	reenrollMutex       sync.RWMutex
	reenrollArgsForCall []struct {
		arg1 *api.ReenrollmentRequest
	}
	reenrollReturns struct {
		result1 *lib.EnrollmentResponse
		result2 error
	}
	reenrollReturnsOnCall map[int]struct {
		result1 *lib.EnrollmentResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Identity) GetClient() *lib.Client {
	fake.getClientMutex.Lock()
	ret, specificReturn := fake.getClientReturnsOnCall[len(fake.getClientArgsForCall)]
	fake.getClientArgsForCall = append(fake.getClientArgsForCall, struct {
	}{})
	fake.recordInvocation("GetClient", []interface{}{})
	fake.getClientMutex.Unlock()
	if fake.GetClientStub != nil {
		return fake.GetClientStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getClientReturns
	return fakeReturns.result1
}

func (fake *Identity) GetClientCallCount() int {
	fake.getClientMutex.RLock()
	defer fake.getClientMutex.RUnlock()
	return len(fake.getClientArgsForCall)
}

func (fake *Identity) GetClientCalls(stub func() *lib.Client) {
	fake.getClientMutex.Lock()
	defer fake.getClientMutex.Unlock()
	fake.GetClientStub = stub
}

func (fake *Identity) GetClientReturns(result1 *lib.Client) {
	fake.getClientMutex.Lock()
	defer fake.getClientMutex.Unlock()
	fake.GetClientStub = nil
	fake.getClientReturns = struct {
		result1 *lib.Client
	}{result1}
}

func (fake *Identity) GetClientReturnsOnCall(i int, result1 *lib.Client) {
	fake.getClientMutex.Lock()
	defer fake.getClientMutex.Unlock()
	fake.GetClientStub = nil
	if fake.getClientReturnsOnCall == nil {
		fake.getClientReturnsOnCall = make(map[int]struct {
			result1 *lib.Client
		})
	}
	fake.getClientReturnsOnCall[i] = struct {
		result1 *lib.Client
	}{result1}
}

func (fake *Identity) GetECert() *x509.Signer {
	fake.getECertMutex.Lock()
	ret, specificReturn := fake.getECertReturnsOnCall[len(fake.getECertArgsForCall)]
	fake.getECertArgsForCall = append(fake.getECertArgsForCall, struct {
	}{})
	fake.recordInvocation("GetECert", []interface{}{})
	fake.getECertMutex.Unlock()
	if fake.GetECertStub != nil {
		return fake.GetECertStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getECertReturns
	return fakeReturns.result1
}

func (fake *Identity) GetECertCallCount() int {
	fake.getECertMutex.RLock()
	defer fake.getECertMutex.RUnlock()
	return len(fake.getECertArgsForCall)
}

func (fake *Identity) GetECertCalls(stub func() *x509.Signer) {
	fake.getECertMutex.Lock()
	defer fake.getECertMutex.Unlock()
	fake.GetECertStub = stub
}

func (fake *Identity) GetECertReturns(result1 *x509.Signer) {
	fake.getECertMutex.Lock()
	defer fake.getECertMutex.Unlock()
	fake.GetECertStub = nil
	fake.getECertReturns = struct {
		result1 *x509.Signer
	}{result1}
}

func (fake *Identity) GetECertReturnsOnCall(i int, result1 *x509.Signer) {
	fake.getECertMutex.Lock()
	defer fake.getECertMutex.Unlock()
	fake.GetECertStub = nil
	if fake.getECertReturnsOnCall == nil {
		fake.getECertReturnsOnCall = make(map[int]struct {
			result1 *x509.Signer
		})
	}
	fake.getECertReturnsOnCall[i] = struct {
		result1 *x509.Signer
	}{result1}
}

func (fake *Identity) Reenroll(arg1 *api.ReenrollmentRequest) (*lib.EnrollmentResponse, error) {
	fake.reenrollMutex.Lock()
	ret, specificReturn := fake.reenrollReturnsOnCall[len(fake.reenrollArgsForCall)]
	fake.reenrollArgsForCall = append(fake.reenrollArgsForCall, struct {
		arg1 *api.ReenrollmentRequest
	}{arg1})
	fake.recordInvocation("Reenroll", []interface{}{arg1})
	fake.reenrollMutex.Unlock()
	if fake.ReenrollStub != nil {
		return fake.ReenrollStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.reenrollReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Identity) ReenrollCallCount() int {
	fake.reenrollMutex.RLock()
	defer fake.reenrollMutex.RUnlock()
	return len(fake.reenrollArgsForCall)
}

func (fake *Identity) ReenrollCalls(stub func(*api.ReenrollmentRequest) (*lib.EnrollmentResponse, error)) {
	fake.reenrollMutex.Lock()
	defer fake.reenrollMutex.Unlock()
	fake.ReenrollStub = stub
}

func (fake *Identity) ReenrollArgsForCall(i int) *api.ReenrollmentRequest {
	fake.reenrollMutex.RLock()
	defer fake.reenrollMutex.RUnlock()
	argsForCall := fake.reenrollArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Identity) ReenrollReturns(result1 *lib.EnrollmentResponse, result2 error) {
	fake.reenrollMutex.Lock()
	defer fake.reenrollMutex.Unlock()
	fake.ReenrollStub = nil
	fake.reenrollReturns = struct {
		result1 *lib.EnrollmentResponse
		result2 error
	}{result1, result2}
}

func (fake *Identity) ReenrollReturnsOnCall(i int, result1 *lib.EnrollmentResponse, result2 error) {
	fake.reenrollMutex.Lock()
	defer fake.reenrollMutex.Unlock()
	fake.ReenrollStub = nil
	if fake.reenrollReturnsOnCall == nil {
		fake.reenrollReturnsOnCall = make(map[int]struct {
			result1 *lib.EnrollmentResponse
			result2 error
		})
	}
	fake.reenrollReturnsOnCall[i] = struct {
		result1 *lib.EnrollmentResponse
		result2 error
	}{result1, result2}
}

func (fake *Identity) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getClientMutex.RLock()
	defer fake.getClientMutex.RUnlock()
	fake.getECertMutex.RLock()
	defer fake.getECertMutex.RUnlock()
	fake.reenrollMutex.RLock()
	defer fake.reenrollMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Identity) recordInvocation(key string, args []interface{}) {
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

var _ reenroller.Identity = new(Identity)
