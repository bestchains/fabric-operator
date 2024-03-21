// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	v25 "github.com/IBM-Blockchain/fabric-operator/pkg/migrator/peer/fabric/v25"
	v1a "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeploymentManager struct {
	DeleteStub        func(v1.Object) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 v1.Object
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	DeploymentStatusStub        func(v1.Object) (v1a.DeploymentStatus, error)
	deploymentStatusMutex       sync.RWMutex
	deploymentStatusArgsForCall []struct {
		arg1 v1.Object
	}
	deploymentStatusReturns struct {
		result1 v1a.DeploymentStatus
		result2 error
	}
	deploymentStatusReturnsOnCall map[int]struct {
		result1 v1a.DeploymentStatus
		result2 error
	}
	GetStub        func(v1.Object) (client.Object, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 v1.Object
	}
	getReturns struct {
		result1 client.Object
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 client.Object
		result2 error
	}
	GetSchemeStub        func() *runtime.Scheme
	getSchemeMutex       sync.RWMutex
	getSchemeArgsForCall []struct {
	}
	getSchemeReturns struct {
		result1 *runtime.Scheme
	}
	getSchemeReturnsOnCall map[int]struct {
		result1 *runtime.Scheme
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DeploymentManager) Delete(arg1 v1.Object) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 v1.Object
	}{arg1})
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1
}

func (fake *DeploymentManager) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *DeploymentManager) DeleteCalls(stub func(v1.Object) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *DeploymentManager) DeleteArgsForCall(i int) v1.Object {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeploymentManager) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeploymentManager) DeleteReturnsOnCall(i int, result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DeploymentManager) DeploymentStatus(arg1 v1.Object) (v1a.DeploymentStatus, error) {
	fake.deploymentStatusMutex.Lock()
	ret, specificReturn := fake.deploymentStatusReturnsOnCall[len(fake.deploymentStatusArgsForCall)]
	fake.deploymentStatusArgsForCall = append(fake.deploymentStatusArgsForCall, struct {
		arg1 v1.Object
	}{arg1})
	fake.recordInvocation("DeploymentStatus", []interface{}{arg1})
	fake.deploymentStatusMutex.Unlock()
	if fake.DeploymentStatusStub != nil {
		return fake.DeploymentStatusStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deploymentStatusReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeploymentManager) DeploymentStatusCallCount() int {
	fake.deploymentStatusMutex.RLock()
	defer fake.deploymentStatusMutex.RUnlock()
	return len(fake.deploymentStatusArgsForCall)
}

func (fake *DeploymentManager) DeploymentStatusCalls(stub func(v1.Object) (v1a.DeploymentStatus, error)) {
	fake.deploymentStatusMutex.Lock()
	defer fake.deploymentStatusMutex.Unlock()
	fake.DeploymentStatusStub = stub
}

func (fake *DeploymentManager) DeploymentStatusArgsForCall(i int) v1.Object {
	fake.deploymentStatusMutex.RLock()
	defer fake.deploymentStatusMutex.RUnlock()
	argsForCall := fake.deploymentStatusArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeploymentManager) DeploymentStatusReturns(result1 v1a.DeploymentStatus, result2 error) {
	fake.deploymentStatusMutex.Lock()
	defer fake.deploymentStatusMutex.Unlock()
	fake.DeploymentStatusStub = nil
	fake.deploymentStatusReturns = struct {
		result1 v1a.DeploymentStatus
		result2 error
	}{result1, result2}
}

func (fake *DeploymentManager) DeploymentStatusReturnsOnCall(i int, result1 v1a.DeploymentStatus, result2 error) {
	fake.deploymentStatusMutex.Lock()
	defer fake.deploymentStatusMutex.Unlock()
	fake.DeploymentStatusStub = nil
	if fake.deploymentStatusReturnsOnCall == nil {
		fake.deploymentStatusReturnsOnCall = make(map[int]struct {
			result1 v1a.DeploymentStatus
			result2 error
		})
	}
	fake.deploymentStatusReturnsOnCall[i] = struct {
		result1 v1a.DeploymentStatus
		result2 error
	}{result1, result2}
}

func (fake *DeploymentManager) Get(arg1 v1.Object) (client.Object, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 v1.Object
	}{arg1})
	fake.recordInvocation("Get", []interface{}{arg1})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeploymentManager) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *DeploymentManager) GetCalls(stub func(v1.Object) (client.Object, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *DeploymentManager) GetArgsForCall(i int) v1.Object {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeploymentManager) GetReturns(result1 client.Object, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 client.Object
		result2 error
	}{result1, result2}
}

func (fake *DeploymentManager) GetReturnsOnCall(i int, result1 client.Object, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 client.Object
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 client.Object
		result2 error
	}{result1, result2}
}

func (fake *DeploymentManager) GetScheme() *runtime.Scheme {
	fake.getSchemeMutex.Lock()
	ret, specificReturn := fake.getSchemeReturnsOnCall[len(fake.getSchemeArgsForCall)]
	fake.getSchemeArgsForCall = append(fake.getSchemeArgsForCall, struct {
	}{})
	fake.recordInvocation("GetScheme", []interface{}{})
	fake.getSchemeMutex.Unlock()
	if fake.GetSchemeStub != nil {
		return fake.GetSchemeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getSchemeReturns
	return fakeReturns.result1
}

func (fake *DeploymentManager) GetSchemeCallCount() int {
	fake.getSchemeMutex.RLock()
	defer fake.getSchemeMutex.RUnlock()
	return len(fake.getSchemeArgsForCall)
}

func (fake *DeploymentManager) GetSchemeCalls(stub func() *runtime.Scheme) {
	fake.getSchemeMutex.Lock()
	defer fake.getSchemeMutex.Unlock()
	fake.GetSchemeStub = stub
}

func (fake *DeploymentManager) GetSchemeReturns(result1 *runtime.Scheme) {
	fake.getSchemeMutex.Lock()
	defer fake.getSchemeMutex.Unlock()
	fake.GetSchemeStub = nil
	fake.getSchemeReturns = struct {
		result1 *runtime.Scheme
	}{result1}
}

func (fake *DeploymentManager) GetSchemeReturnsOnCall(i int, result1 *runtime.Scheme) {
	fake.getSchemeMutex.Lock()
	defer fake.getSchemeMutex.Unlock()
	fake.GetSchemeStub = nil
	if fake.getSchemeReturnsOnCall == nil {
		fake.getSchemeReturnsOnCall = make(map[int]struct {
			result1 *runtime.Scheme
		})
	}
	fake.getSchemeReturnsOnCall[i] = struct {
		result1 *runtime.Scheme
	}{result1}
}

func (fake *DeploymentManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.deploymentStatusMutex.RLock()
	defer fake.deploymentStatusMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.getSchemeMutex.RLock()
	defer fake.getSchemeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DeploymentManager) recordInvocation(key string, args []interface{}) {
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

var _ v25.DeploymentManager = new(DeploymentManager)
