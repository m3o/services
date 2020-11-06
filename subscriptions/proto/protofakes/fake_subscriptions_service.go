// Code generated by counterfeiter. DO NOT EDIT.
package protofakes

import (
	"context"
	"sync"

	subscriptions "github.com/m3o/services/subscriptions/proto"
	"github.com/micro/micro/v3/service/client"
)

type FakeSubscriptionsService struct {
	AddUserStub        func(context.Context, *subscriptions.AddUserRequest, ...client.CallOption) (*subscriptions.AddUserResponse, error)
	addUserMutex       sync.RWMutex
	addUserArgsForCall []struct {
		arg1 context.Context
		arg2 *subscriptions.AddUserRequest
		arg3 []client.CallOption
	}
	addUserReturns struct {
		result1 *subscriptions.AddUserResponse
		result2 error
	}
	addUserReturnsOnCall map[int]struct {
		result1 *subscriptions.AddUserResponse
		result2 error
	}
	CancelStub        func(context.Context, *subscriptions.CancelRequest, ...client.CallOption) (*subscriptions.CancelResponse, error)
	cancelMutex       sync.RWMutex
	cancelArgsForCall []struct {
		arg1 context.Context
		arg2 *subscriptions.CancelRequest
		arg3 []client.CallOption
	}
	cancelReturns struct {
		result1 *subscriptions.CancelResponse
		result2 error
	}
	cancelReturnsOnCall map[int]struct {
		result1 *subscriptions.CancelResponse
		result2 error
	}
	CreateStub        func(context.Context, *subscriptions.CreateRequest, ...client.CallOption) (*subscriptions.CreateResponse, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 context.Context
		arg2 *subscriptions.CreateRequest
		arg3 []client.CallOption
	}
	createReturns struct {
		result1 *subscriptions.CreateResponse
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 *subscriptions.CreateResponse
		result2 error
	}
	UpdateStub        func(context.Context, *subscriptions.UpdateRequest, ...client.CallOption) (*subscriptions.UpdateResponse, error)
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 context.Context
		arg2 *subscriptions.UpdateRequest
		arg3 []client.CallOption
	}
	updateReturns struct {
		result1 *subscriptions.UpdateResponse
		result2 error
	}
	updateReturnsOnCall map[int]struct {
		result1 *subscriptions.UpdateResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSubscriptionsService) AddUser(arg1 context.Context, arg2 *subscriptions.AddUserRequest, arg3 ...client.CallOption) (*subscriptions.AddUserResponse, error) {
	fake.addUserMutex.Lock()
	ret, specificReturn := fake.addUserReturnsOnCall[len(fake.addUserArgsForCall)]
	fake.addUserArgsForCall = append(fake.addUserArgsForCall, struct {
		arg1 context.Context
		arg2 *subscriptions.AddUserRequest
		arg3 []client.CallOption
	}{arg1, arg2, arg3})
	stub := fake.AddUserStub
	fakeReturns := fake.addUserReturns
	fake.recordInvocation("AddUser", []interface{}{arg1, arg2, arg3})
	fake.addUserMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSubscriptionsService) AddUserCallCount() int {
	fake.addUserMutex.RLock()
	defer fake.addUserMutex.RUnlock()
	return len(fake.addUserArgsForCall)
}

func (fake *FakeSubscriptionsService) AddUserCalls(stub func(context.Context, *subscriptions.AddUserRequest, ...client.CallOption) (*subscriptions.AddUserResponse, error)) {
	fake.addUserMutex.Lock()
	defer fake.addUserMutex.Unlock()
	fake.AddUserStub = stub
}

func (fake *FakeSubscriptionsService) AddUserArgsForCall(i int) (context.Context, *subscriptions.AddUserRequest, []client.CallOption) {
	fake.addUserMutex.RLock()
	defer fake.addUserMutex.RUnlock()
	argsForCall := fake.addUserArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeSubscriptionsService) AddUserReturns(result1 *subscriptions.AddUserResponse, result2 error) {
	fake.addUserMutex.Lock()
	defer fake.addUserMutex.Unlock()
	fake.AddUserStub = nil
	fake.addUserReturns = struct {
		result1 *subscriptions.AddUserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) AddUserReturnsOnCall(i int, result1 *subscriptions.AddUserResponse, result2 error) {
	fake.addUserMutex.Lock()
	defer fake.addUserMutex.Unlock()
	fake.AddUserStub = nil
	if fake.addUserReturnsOnCall == nil {
		fake.addUserReturnsOnCall = make(map[int]struct {
			result1 *subscriptions.AddUserResponse
			result2 error
		})
	}
	fake.addUserReturnsOnCall[i] = struct {
		result1 *subscriptions.AddUserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) Cancel(arg1 context.Context, arg2 *subscriptions.CancelRequest, arg3 ...client.CallOption) (*subscriptions.CancelResponse, error) {
	fake.cancelMutex.Lock()
	ret, specificReturn := fake.cancelReturnsOnCall[len(fake.cancelArgsForCall)]
	fake.cancelArgsForCall = append(fake.cancelArgsForCall, struct {
		arg1 context.Context
		arg2 *subscriptions.CancelRequest
		arg3 []client.CallOption
	}{arg1, arg2, arg3})
	stub := fake.CancelStub
	fakeReturns := fake.cancelReturns
	fake.recordInvocation("Cancel", []interface{}{arg1, arg2, arg3})
	fake.cancelMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSubscriptionsService) CancelCallCount() int {
	fake.cancelMutex.RLock()
	defer fake.cancelMutex.RUnlock()
	return len(fake.cancelArgsForCall)
}

func (fake *FakeSubscriptionsService) CancelCalls(stub func(context.Context, *subscriptions.CancelRequest, ...client.CallOption) (*subscriptions.CancelResponse, error)) {
	fake.cancelMutex.Lock()
	defer fake.cancelMutex.Unlock()
	fake.CancelStub = stub
}

func (fake *FakeSubscriptionsService) CancelArgsForCall(i int) (context.Context, *subscriptions.CancelRequest, []client.CallOption) {
	fake.cancelMutex.RLock()
	defer fake.cancelMutex.RUnlock()
	argsForCall := fake.cancelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeSubscriptionsService) CancelReturns(result1 *subscriptions.CancelResponse, result2 error) {
	fake.cancelMutex.Lock()
	defer fake.cancelMutex.Unlock()
	fake.CancelStub = nil
	fake.cancelReturns = struct {
		result1 *subscriptions.CancelResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) CancelReturnsOnCall(i int, result1 *subscriptions.CancelResponse, result2 error) {
	fake.cancelMutex.Lock()
	defer fake.cancelMutex.Unlock()
	fake.CancelStub = nil
	if fake.cancelReturnsOnCall == nil {
		fake.cancelReturnsOnCall = make(map[int]struct {
			result1 *subscriptions.CancelResponse
			result2 error
		})
	}
	fake.cancelReturnsOnCall[i] = struct {
		result1 *subscriptions.CancelResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) Create(arg1 context.Context, arg2 *subscriptions.CreateRequest, arg3 ...client.CallOption) (*subscriptions.CreateResponse, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 context.Context
		arg2 *subscriptions.CreateRequest
		arg3 []client.CallOption
	}{arg1, arg2, arg3})
	stub := fake.CreateStub
	fakeReturns := fake.createReturns
	fake.recordInvocation("Create", []interface{}{arg1, arg2, arg3})
	fake.createMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSubscriptionsService) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeSubscriptionsService) CreateCalls(stub func(context.Context, *subscriptions.CreateRequest, ...client.CallOption) (*subscriptions.CreateResponse, error)) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeSubscriptionsService) CreateArgsForCall(i int) (context.Context, *subscriptions.CreateRequest, []client.CallOption) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeSubscriptionsService) CreateReturns(result1 *subscriptions.CreateResponse, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 *subscriptions.CreateResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) CreateReturnsOnCall(i int, result1 *subscriptions.CreateResponse, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 *subscriptions.CreateResponse
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 *subscriptions.CreateResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) Update(arg1 context.Context, arg2 *subscriptions.UpdateRequest, arg3 ...client.CallOption) (*subscriptions.UpdateResponse, error) {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 context.Context
		arg2 *subscriptions.UpdateRequest
		arg3 []client.CallOption
	}{arg1, arg2, arg3})
	stub := fake.UpdateStub
	fakeReturns := fake.updateReturns
	fake.recordInvocation("Update", []interface{}{arg1, arg2, arg3})
	fake.updateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSubscriptionsService) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeSubscriptionsService) UpdateCalls(stub func(context.Context, *subscriptions.UpdateRequest, ...client.CallOption) (*subscriptions.UpdateResponse, error)) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *FakeSubscriptionsService) UpdateArgsForCall(i int) (context.Context, *subscriptions.UpdateRequest, []client.CallOption) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeSubscriptionsService) UpdateReturns(result1 *subscriptions.UpdateResponse, result2 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 *subscriptions.UpdateResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) UpdateReturnsOnCall(i int, result1 *subscriptions.UpdateResponse, result2 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 *subscriptions.UpdateResponse
			result2 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 *subscriptions.UpdateResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSubscriptionsService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addUserMutex.RLock()
	defer fake.addUserMutex.RUnlock()
	fake.cancelMutex.RLock()
	defer fake.cancelMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSubscriptionsService) recordInvocation(key string, args []interface{}) {
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

var _ subscriptions.SubscriptionsService = new(FakeSubscriptionsService)
