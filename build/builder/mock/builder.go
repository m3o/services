// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/m3o/services/build/builder"
)

type MockBuilder struct {
	BuildStub        func(string, string, string) (string, error)
	buildMutex       sync.RWMutex
	buildArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	buildReturns struct {
		result1 string
		result2 error
	}
	buildReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	PushStub        func(string) (string, error)
	pushMutex       sync.RWMutex
	pushArgsForCall []struct {
		arg1 string
	}
	pushReturns struct {
		result1 string
		result2 error
	}
	pushReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MockBuilder) Build(arg1 string, arg2 string, arg3 string) (string, error) {
	fake.buildMutex.Lock()
	ret, specificReturn := fake.buildReturnsOnCall[len(fake.buildArgsForCall)]
	fake.buildArgsForCall = append(fake.buildArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("Build", []interface{}{arg1, arg2, arg3})
	fake.buildMutex.Unlock()
	if fake.BuildStub != nil {
		return fake.BuildStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.buildReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *MockBuilder) BuildCallCount() int {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return len(fake.buildArgsForCall)
}

func (fake *MockBuilder) BuildCalls(stub func(string, string, string) (string, error)) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = stub
}

func (fake *MockBuilder) BuildArgsForCall(i int) (string, string, string) {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	argsForCall := fake.buildArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *MockBuilder) BuildReturns(result1 string, result2 error) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = nil
	fake.buildReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *MockBuilder) BuildReturnsOnCall(i int, result1 string, result2 error) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = nil
	if fake.buildReturnsOnCall == nil {
		fake.buildReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.buildReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *MockBuilder) Push(arg1 string) (string, error) {
	fake.pushMutex.Lock()
	ret, specificReturn := fake.pushReturnsOnCall[len(fake.pushArgsForCall)]
	fake.pushArgsForCall = append(fake.pushArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Push", []interface{}{arg1})
	fake.pushMutex.Unlock()
	if fake.PushStub != nil {
		return fake.PushStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.pushReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *MockBuilder) PushCallCount() int {
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	return len(fake.pushArgsForCall)
}

func (fake *MockBuilder) PushCalls(stub func(string) (string, error)) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = stub
}

func (fake *MockBuilder) PushArgsForCall(i int) string {
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	argsForCall := fake.pushArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MockBuilder) PushReturns(result1 string, result2 error) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = nil
	fake.pushReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *MockBuilder) PushReturnsOnCall(i int, result1 string, result2 error) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = nil
	if fake.pushReturnsOnCall == nil {
		fake.pushReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.pushReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *MockBuilder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MockBuilder) recordInvocation(key string, args []interface{}) {
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

var _ builder.Builder = new(MockBuilder)
