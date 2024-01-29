// Code generated by counterfeiter. DO NOT EDIT.
package repositoryfakes

import (
	"library/transactions/models"
	"library/transactions/repository"
	"sync"
)

type FakeTransactionerRepository struct {
	BuyBookStub        func(int, int, int) (int, error)
	buyBookMutex       sync.RWMutex
	buyBookArgsForCall []struct {
		arg1 int
		arg2 int
		arg3 int
	}
	buyBookReturns struct {
		result1 int
		result2 error
	}
	buyBookReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	TransactionHistoryStub        func(int) ([]models.UserTransactionResponse, error)
	transactionHistoryMutex       sync.RWMutex
	transactionHistoryArgsForCall []struct {
		arg1 int
	}
	transactionHistoryReturns struct {
		result1 []models.UserTransactionResponse
		result2 error
	}
	transactionHistoryReturnsOnCall map[int]struct {
		result1 []models.UserTransactionResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTransactionerRepository) BuyBook(arg1 int, arg2 int, arg3 int) (int, error) {
	fake.buyBookMutex.Lock()
	ret, specificReturn := fake.buyBookReturnsOnCall[len(fake.buyBookArgsForCall)]
	fake.buyBookArgsForCall = append(fake.buyBookArgsForCall, struct {
		arg1 int
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.BuyBookStub
	fakeReturns := fake.buyBookReturns
	fake.recordInvocation("BuyBook", []interface{}{arg1, arg2, arg3})
	fake.buyBookMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTransactionerRepository) BuyBookCallCount() int {
	fake.buyBookMutex.RLock()
	defer fake.buyBookMutex.RUnlock()
	return len(fake.buyBookArgsForCall)
}

func (fake *FakeTransactionerRepository) BuyBookCalls(stub func(int, int, int) (int, error)) {
	fake.buyBookMutex.Lock()
	defer fake.buyBookMutex.Unlock()
	fake.BuyBookStub = stub
}

func (fake *FakeTransactionerRepository) BuyBookArgsForCall(i int) (int, int, int) {
	fake.buyBookMutex.RLock()
	defer fake.buyBookMutex.RUnlock()
	argsForCall := fake.buyBookArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeTransactionerRepository) BuyBookReturns(result1 int, result2 error) {
	fake.buyBookMutex.Lock()
	defer fake.buyBookMutex.Unlock()
	fake.BuyBookStub = nil
	fake.buyBookReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeTransactionerRepository) BuyBookReturnsOnCall(i int, result1 int, result2 error) {
	fake.buyBookMutex.Lock()
	defer fake.buyBookMutex.Unlock()
	fake.BuyBookStub = nil
	if fake.buyBookReturnsOnCall == nil {
		fake.buyBookReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.buyBookReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeTransactionerRepository) TransactionHistory(arg1 int) ([]models.UserTransactionResponse, error) {
	fake.transactionHistoryMutex.Lock()
	ret, specificReturn := fake.transactionHistoryReturnsOnCall[len(fake.transactionHistoryArgsForCall)]
	fake.transactionHistoryArgsForCall = append(fake.transactionHistoryArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.TransactionHistoryStub
	fakeReturns := fake.transactionHistoryReturns
	fake.recordInvocation("TransactionHistory", []interface{}{arg1})
	fake.transactionHistoryMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTransactionerRepository) TransactionHistoryCallCount() int {
	fake.transactionHistoryMutex.RLock()
	defer fake.transactionHistoryMutex.RUnlock()
	return len(fake.transactionHistoryArgsForCall)
}

func (fake *FakeTransactionerRepository) TransactionHistoryCalls(stub func(int) ([]models.UserTransactionResponse, error)) {
	fake.transactionHistoryMutex.Lock()
	defer fake.transactionHistoryMutex.Unlock()
	fake.TransactionHistoryStub = stub
}

func (fake *FakeTransactionerRepository) TransactionHistoryArgsForCall(i int) int {
	fake.transactionHistoryMutex.RLock()
	defer fake.transactionHistoryMutex.RUnlock()
	argsForCall := fake.transactionHistoryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTransactionerRepository) TransactionHistoryReturns(result1 []models.UserTransactionResponse, result2 error) {
	fake.transactionHistoryMutex.Lock()
	defer fake.transactionHistoryMutex.Unlock()
	fake.TransactionHistoryStub = nil
	fake.transactionHistoryReturns = struct {
		result1 []models.UserTransactionResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeTransactionerRepository) TransactionHistoryReturnsOnCall(i int, result1 []models.UserTransactionResponse, result2 error) {
	fake.transactionHistoryMutex.Lock()
	defer fake.transactionHistoryMutex.Unlock()
	fake.TransactionHistoryStub = nil
	if fake.transactionHistoryReturnsOnCall == nil {
		fake.transactionHistoryReturnsOnCall = make(map[int]struct {
			result1 []models.UserTransactionResponse
			result2 error
		})
	}
	fake.transactionHistoryReturnsOnCall[i] = struct {
		result1 []models.UserTransactionResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeTransactionerRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buyBookMutex.RLock()
	defer fake.buyBookMutex.RUnlock()
	fake.transactionHistoryMutex.RLock()
	defer fake.transactionHistoryMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTransactionerRepository) recordInvocation(key string, args []interface{}) {
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

var _ repository.TransactionerRepository = new(FakeTransactionerRepository)