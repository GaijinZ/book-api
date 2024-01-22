// Code generated by counterfeiter. DO NOT EDIT.
package repositoryfakes

import (
	"library/pkg/postgres"
	"library/users/models"
	"library/users/repository"
	"sync"
)

type FakeUsererRepository struct {
	AddUserStub        func(*models.User) (int, error)
	addUserMutex       sync.RWMutex
	addUserArgsForCall []struct {
		arg1 *models.User
	}
	addUserReturns struct {
		result1 int
		result2 error
	}
	addUserReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	DeleteUserStub        func(int) (int, error)
	deleteUserMutex       sync.RWMutex
	deleteUserArgsForCall []struct {
		arg1 int
	}
	deleteUserReturns struct {
		result1 int
		result2 error
	}
	deleteUserReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	GetAllUsersStub        func() ([]models.UserResponse, error)
	getAllUsersMutex       sync.RWMutex
	getAllUsersArgsForCall []struct {
	}
	getAllUsersReturns struct {
		result1 []models.UserResponse
		result2 error
	}
	getAllUsersReturnsOnCall map[int]struct {
		result1 []models.UserResponse
		result2 error
	}
	GetDBPoolStub        func() postgres.DB
	getDBPoolMutex       sync.RWMutex
	getDBPoolArgsForCall []struct {
	}
	getDBPoolReturns struct {
		result1 postgres.DB
	}
	getDBPoolReturnsOnCall map[int]struct {
		result1 postgres.DB
	}
	GetUserStub        func(int) (*models.UserResponse, error)
	getUserMutex       sync.RWMutex
	getUserArgsForCall []struct {
		arg1 int
	}
	getUserReturns struct {
		result1 *models.UserResponse
		result2 error
	}
	getUserReturnsOnCall map[int]struct {
		result1 *models.UserResponse
		result2 error
	}
	UpdateUserStub        func(*models.User) (*models.UserResponse, error)
	updateUserMutex       sync.RWMutex
	updateUserArgsForCall []struct {
		arg1 *models.User
	}
	updateUserReturns struct {
		result1 *models.UserResponse
		result2 error
	}
	updateUserReturnsOnCall map[int]struct {
		result1 *models.UserResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeUsererRepository) AddUser(arg1 *models.User) (int, error) {
	fake.addUserMutex.Lock()
	ret, specificReturn := fake.addUserReturnsOnCall[len(fake.addUserArgsForCall)]
	fake.addUserArgsForCall = append(fake.addUserArgsForCall, struct {
		arg1 *models.User
	}{arg1})
	stub := fake.AddUserStub
	fakeReturns := fake.addUserReturns
	fake.recordInvocation("AddUser", []interface{}{arg1})
	fake.addUserMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUsererRepository) AddUserCallCount() int {
	fake.addUserMutex.RLock()
	defer fake.addUserMutex.RUnlock()
	return len(fake.addUserArgsForCall)
}

func (fake *FakeUsererRepository) AddUserCalls(stub func(*models.User) (int, error)) {
	fake.addUserMutex.Lock()
	defer fake.addUserMutex.Unlock()
	fake.AddUserStub = stub
}

func (fake *FakeUsererRepository) AddUserArgsForCall(i int) *models.User {
	fake.addUserMutex.RLock()
	defer fake.addUserMutex.RUnlock()
	argsForCall := fake.addUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUsererRepository) AddUserReturns(result1 int, result2 error) {
	fake.addUserMutex.Lock()
	defer fake.addUserMutex.Unlock()
	fake.AddUserStub = nil
	fake.addUserReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) AddUserReturnsOnCall(i int, result1 int, result2 error) {
	fake.addUserMutex.Lock()
	defer fake.addUserMutex.Unlock()
	fake.AddUserStub = nil
	if fake.addUserReturnsOnCall == nil {
		fake.addUserReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.addUserReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) DeleteUser(arg1 int) (int, error) {
	fake.deleteUserMutex.Lock()
	ret, specificReturn := fake.deleteUserReturnsOnCall[len(fake.deleteUserArgsForCall)]
	fake.deleteUserArgsForCall = append(fake.deleteUserArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.DeleteUserStub
	fakeReturns := fake.deleteUserReturns
	fake.recordInvocation("DeleteUser", []interface{}{arg1})
	fake.deleteUserMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUsererRepository) DeleteUserCallCount() int {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	return len(fake.deleteUserArgsForCall)
}

func (fake *FakeUsererRepository) DeleteUserCalls(stub func(int) (int, error)) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = stub
}

func (fake *FakeUsererRepository) DeleteUserArgsForCall(i int) int {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	argsForCall := fake.deleteUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUsererRepository) DeleteUserReturns(result1 int, result2 error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = nil
	fake.deleteUserReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) DeleteUserReturnsOnCall(i int, result1 int, result2 error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = nil
	if fake.deleteUserReturnsOnCall == nil {
		fake.deleteUserReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.deleteUserReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) GetAllUsers() ([]models.UserResponse, error) {
	fake.getAllUsersMutex.Lock()
	ret, specificReturn := fake.getAllUsersReturnsOnCall[len(fake.getAllUsersArgsForCall)]
	fake.getAllUsersArgsForCall = append(fake.getAllUsersArgsForCall, struct {
	}{})
	stub := fake.GetAllUsersStub
	fakeReturns := fake.getAllUsersReturns
	fake.recordInvocation("GetAllUsers", []interface{}{})
	fake.getAllUsersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUsererRepository) GetAllUsersCallCount() int {
	fake.getAllUsersMutex.RLock()
	defer fake.getAllUsersMutex.RUnlock()
	return len(fake.getAllUsersArgsForCall)
}

func (fake *FakeUsererRepository) GetAllUsersCalls(stub func() ([]models.UserResponse, error)) {
	fake.getAllUsersMutex.Lock()
	defer fake.getAllUsersMutex.Unlock()
	fake.GetAllUsersStub = stub
}

func (fake *FakeUsererRepository) GetAllUsersReturns(result1 []models.UserResponse, result2 error) {
	fake.getAllUsersMutex.Lock()
	defer fake.getAllUsersMutex.Unlock()
	fake.GetAllUsersStub = nil
	fake.getAllUsersReturns = struct {
		result1 []models.UserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) GetAllUsersReturnsOnCall(i int, result1 []models.UserResponse, result2 error) {
	fake.getAllUsersMutex.Lock()
	defer fake.getAllUsersMutex.Unlock()
	fake.GetAllUsersStub = nil
	if fake.getAllUsersReturnsOnCall == nil {
		fake.getAllUsersReturnsOnCall = make(map[int]struct {
			result1 []models.UserResponse
			result2 error
		})
	}
	fake.getAllUsersReturnsOnCall[i] = struct {
		result1 []models.UserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) GetDBPool() postgres.DB {
	fake.getDBPoolMutex.Lock()
	ret, specificReturn := fake.getDBPoolReturnsOnCall[len(fake.getDBPoolArgsForCall)]
	fake.getDBPoolArgsForCall = append(fake.getDBPoolArgsForCall, struct {
	}{})
	stub := fake.GetDBPoolStub
	fakeReturns := fake.getDBPoolReturns
	fake.recordInvocation("GetDBPool", []interface{}{})
	fake.getDBPoolMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeUsererRepository) GetDBPoolCallCount() int {
	fake.getDBPoolMutex.RLock()
	defer fake.getDBPoolMutex.RUnlock()
	return len(fake.getDBPoolArgsForCall)
}

func (fake *FakeUsererRepository) GetDBPoolCalls(stub func() postgres.DB) {
	fake.getDBPoolMutex.Lock()
	defer fake.getDBPoolMutex.Unlock()
	fake.GetDBPoolStub = stub
}

func (fake *FakeUsererRepository) GetDBPoolReturns(result1 postgres.DB) {
	fake.getDBPoolMutex.Lock()
	defer fake.getDBPoolMutex.Unlock()
	fake.GetDBPoolStub = nil
	fake.getDBPoolReturns = struct {
		result1 postgres.DB
	}{result1}
}

func (fake *FakeUsererRepository) GetDBPoolReturnsOnCall(i int, result1 postgres.DB) {
	fake.getDBPoolMutex.Lock()
	defer fake.getDBPoolMutex.Unlock()
	fake.GetDBPoolStub = nil
	if fake.getDBPoolReturnsOnCall == nil {
		fake.getDBPoolReturnsOnCall = make(map[int]struct {
			result1 postgres.DB
		})
	}
	fake.getDBPoolReturnsOnCall[i] = struct {
		result1 postgres.DB
	}{result1}
}

func (fake *FakeUsererRepository) GetUser(arg1 int) (*models.UserResponse, error) {
	fake.getUserMutex.Lock()
	ret, specificReturn := fake.getUserReturnsOnCall[len(fake.getUserArgsForCall)]
	fake.getUserArgsForCall = append(fake.getUserArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.GetUserStub
	fakeReturns := fake.getUserReturns
	fake.recordInvocation("GetUser", []interface{}{arg1})
	fake.getUserMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUsererRepository) GetUserCallCount() int {
	fake.getUserMutex.RLock()
	defer fake.getUserMutex.RUnlock()
	return len(fake.getUserArgsForCall)
}

func (fake *FakeUsererRepository) GetUserCalls(stub func(int) (*models.UserResponse, error)) {
	fake.getUserMutex.Lock()
	defer fake.getUserMutex.Unlock()
	fake.GetUserStub = stub
}

func (fake *FakeUsererRepository) GetUserArgsForCall(i int) int {
	fake.getUserMutex.RLock()
	defer fake.getUserMutex.RUnlock()
	argsForCall := fake.getUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUsererRepository) GetUserReturns(result1 *models.UserResponse, result2 error) {
	fake.getUserMutex.Lock()
	defer fake.getUserMutex.Unlock()
	fake.GetUserStub = nil
	fake.getUserReturns = struct {
		result1 *models.UserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) GetUserReturnsOnCall(i int, result1 *models.UserResponse, result2 error) {
	fake.getUserMutex.Lock()
	defer fake.getUserMutex.Unlock()
	fake.GetUserStub = nil
	if fake.getUserReturnsOnCall == nil {
		fake.getUserReturnsOnCall = make(map[int]struct {
			result1 *models.UserResponse
			result2 error
		})
	}
	fake.getUserReturnsOnCall[i] = struct {
		result1 *models.UserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) UpdateUser(arg1 *models.User) (*models.UserResponse, error) {
	fake.updateUserMutex.Lock()
	ret, specificReturn := fake.updateUserReturnsOnCall[len(fake.updateUserArgsForCall)]
	fake.updateUserArgsForCall = append(fake.updateUserArgsForCall, struct {
		arg1 *models.User
	}{arg1})
	stub := fake.UpdateUserStub
	fakeReturns := fake.updateUserReturns
	fake.recordInvocation("UpdateUser", []interface{}{arg1})
	fake.updateUserMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUsererRepository) UpdateUserCallCount() int {
	fake.updateUserMutex.RLock()
	defer fake.updateUserMutex.RUnlock()
	return len(fake.updateUserArgsForCall)
}

func (fake *FakeUsererRepository) UpdateUserCalls(stub func(*models.User) (*models.UserResponse, error)) {
	fake.updateUserMutex.Lock()
	defer fake.updateUserMutex.Unlock()
	fake.UpdateUserStub = stub
}

func (fake *FakeUsererRepository) UpdateUserArgsForCall(i int) *models.User {
	fake.updateUserMutex.RLock()
	defer fake.updateUserMutex.RUnlock()
	argsForCall := fake.updateUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUsererRepository) UpdateUserReturns(result1 *models.UserResponse, result2 error) {
	fake.updateUserMutex.Lock()
	defer fake.updateUserMutex.Unlock()
	fake.UpdateUserStub = nil
	fake.updateUserReturns = struct {
		result1 *models.UserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) UpdateUserReturnsOnCall(i int, result1 *models.UserResponse, result2 error) {
	fake.updateUserMutex.Lock()
	defer fake.updateUserMutex.Unlock()
	fake.UpdateUserStub = nil
	if fake.updateUserReturnsOnCall == nil {
		fake.updateUserReturnsOnCall = make(map[int]struct {
			result1 *models.UserResponse
			result2 error
		})
	}
	fake.updateUserReturnsOnCall[i] = struct {
		result1 *models.UserResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsererRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addUserMutex.RLock()
	defer fake.addUserMutex.RUnlock()
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	fake.getAllUsersMutex.RLock()
	defer fake.getAllUsersMutex.RUnlock()
	fake.getDBPoolMutex.RLock()
	defer fake.getDBPoolMutex.RUnlock()
	fake.getUserMutex.RLock()
	defer fake.getUserMutex.RUnlock()
	fake.updateUserMutex.RLock()
	defer fake.updateUserMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeUsererRepository) recordInvocation(key string, args []interface{}) {
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

var _ repository.UsererRepository = new(FakeUsererRepository)