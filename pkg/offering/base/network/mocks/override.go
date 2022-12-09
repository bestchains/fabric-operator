// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	"github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/network"
	v1a "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Override struct {
	ClusterRoleStub        func(v1.Object, *v1a.ClusterRole, resources.Action) error
	clusterRoleMutex       sync.RWMutex
	clusterRoleArgsForCall []struct {
		arg1 v1.Object
		arg2 *v1a.ClusterRole
		arg3 resources.Action
	}
	clusterRoleReturns struct {
		result1 error
	}
	clusterRoleReturnsOnCall map[int]struct {
		result1 error
	}
	ClusterRoleBindingStub        func(v1.Object, *v1a.ClusterRoleBinding, resources.Action) error
	clusterRoleBindingMutex       sync.RWMutex
	clusterRoleBindingArgsForCall []struct {
		arg1 v1.Object
		arg2 *v1a.ClusterRoleBinding
		arg3 resources.Action
	}
	clusterRoleBindingReturns struct {
		result1 error
	}
	clusterRoleBindingReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Override) ClusterRole(arg1 v1.Object, arg2 *v1a.ClusterRole, arg3 resources.Action) error {
	fake.clusterRoleMutex.Lock()
	ret, specificReturn := fake.clusterRoleReturnsOnCall[len(fake.clusterRoleArgsForCall)]
	fake.clusterRoleArgsForCall = append(fake.clusterRoleArgsForCall, struct {
		arg1 v1.Object
		arg2 *v1a.ClusterRole
		arg3 resources.Action
	}{arg1, arg2, arg3})
	fake.recordInvocation("ClusterRole", []interface{}{arg1, arg2, arg3})
	fake.clusterRoleMutex.Unlock()
	if fake.ClusterRoleStub != nil {
		return fake.ClusterRoleStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.clusterRoleReturns
	return fakeReturns.result1
}

func (fake *Override) ClusterRoleCallCount() int {
	fake.clusterRoleMutex.RLock()
	defer fake.clusterRoleMutex.RUnlock()
	return len(fake.clusterRoleArgsForCall)
}

func (fake *Override) ClusterRoleCalls(stub func(v1.Object, *v1a.ClusterRole, resources.Action) error) {
	fake.clusterRoleMutex.Lock()
	defer fake.clusterRoleMutex.Unlock()
	fake.ClusterRoleStub = stub
}

func (fake *Override) ClusterRoleArgsForCall(i int) (v1.Object, *v1a.ClusterRole, resources.Action) {
	fake.clusterRoleMutex.RLock()
	defer fake.clusterRoleMutex.RUnlock()
	argsForCall := fake.clusterRoleArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *Override) ClusterRoleReturns(result1 error) {
	fake.clusterRoleMutex.Lock()
	defer fake.clusterRoleMutex.Unlock()
	fake.ClusterRoleStub = nil
	fake.clusterRoleReturns = struct {
		result1 error
	}{result1}
}

func (fake *Override) ClusterRoleReturnsOnCall(i int, result1 error) {
	fake.clusterRoleMutex.Lock()
	defer fake.clusterRoleMutex.Unlock()
	fake.ClusterRoleStub = nil
	if fake.clusterRoleReturnsOnCall == nil {
		fake.clusterRoleReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.clusterRoleReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Override) ClusterRoleBinding(arg1 v1.Object, arg2 *v1a.ClusterRoleBinding, arg3 resources.Action) error {
	fake.clusterRoleBindingMutex.Lock()
	ret, specificReturn := fake.clusterRoleBindingReturnsOnCall[len(fake.clusterRoleBindingArgsForCall)]
	fake.clusterRoleBindingArgsForCall = append(fake.clusterRoleBindingArgsForCall, struct {
		arg1 v1.Object
		arg2 *v1a.ClusterRoleBinding
		arg3 resources.Action
	}{arg1, arg2, arg3})
	fake.recordInvocation("ClusterRoleBinding", []interface{}{arg1, arg2, arg3})
	fake.clusterRoleBindingMutex.Unlock()
	if fake.ClusterRoleBindingStub != nil {
		return fake.ClusterRoleBindingStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.clusterRoleBindingReturns
	return fakeReturns.result1
}

func (fake *Override) ClusterRoleBindingCallCount() int {
	fake.clusterRoleBindingMutex.RLock()
	defer fake.clusterRoleBindingMutex.RUnlock()
	return len(fake.clusterRoleBindingArgsForCall)
}

func (fake *Override) ClusterRoleBindingCalls(stub func(v1.Object, *v1a.ClusterRoleBinding, resources.Action) error) {
	fake.clusterRoleBindingMutex.Lock()
	defer fake.clusterRoleBindingMutex.Unlock()
	fake.ClusterRoleBindingStub = stub
}

func (fake *Override) ClusterRoleBindingArgsForCall(i int) (v1.Object, *v1a.ClusterRoleBinding, resources.Action) {
	fake.clusterRoleBindingMutex.RLock()
	defer fake.clusterRoleBindingMutex.RUnlock()
	argsForCall := fake.clusterRoleBindingArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *Override) ClusterRoleBindingReturns(result1 error) {
	fake.clusterRoleBindingMutex.Lock()
	defer fake.clusterRoleBindingMutex.Unlock()
	fake.ClusterRoleBindingStub = nil
	fake.clusterRoleBindingReturns = struct {
		result1 error
	}{result1}
}

func (fake *Override) ClusterRoleBindingReturnsOnCall(i int, result1 error) {
	fake.clusterRoleBindingMutex.Lock()
	defer fake.clusterRoleBindingMutex.Unlock()
	fake.ClusterRoleBindingStub = nil
	if fake.clusterRoleBindingReturnsOnCall == nil {
		fake.clusterRoleBindingReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.clusterRoleBindingReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Override) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.clusterRoleMutex.RLock()
	defer fake.clusterRoleMutex.RUnlock()
	fake.clusterRoleBindingMutex.RLock()
	defer fake.clusterRoleBindingMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Override) recordInvocation(key string, args []interface{}) {
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

var _ network.Override = new(Override)