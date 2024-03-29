// Code generated by counterfeiter. DO NOT EDIT.
package repositoryfakes

import (
	"library/books/models"
	"library/books/repository"
	"library/pkg/postgres"
	"sync"
)

type FakeBookerRepository struct {
	AddBookStub        func(*models.BookRequest) (int, error)
	addBookMutex       sync.RWMutex
	addBookArgsForCall []struct {
		arg1 *models.BookRequest
	}
	addBookReturns struct {
		result1 int
		result2 error
	}
	addBookReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	DeleteBookStub        func(int, int) (int, error)
	deleteBookMutex       sync.RWMutex
	deleteBookArgsForCall []struct {
		arg1 int
		arg2 int
	}
	deleteBookReturns struct {
		result1 int
		result2 error
	}
	deleteBookReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	GetAllBooksStub        func() ([]models.BookResponse, error)
	getAllBooksMutex       sync.RWMutex
	getAllBooksArgsForCall []struct {
	}
	getAllBooksReturns struct {
		result1 []models.BookResponse
		result2 error
	}
	getAllBooksReturnsOnCall map[int]struct {
		result1 []models.BookResponse
		result2 error
	}
	GetBookStub        func(int) (*models.BookResponse, error)
	getBookMutex       sync.RWMutex
	getBookArgsForCall []struct {
		arg1 int
	}
	getBookReturns struct {
		result1 *models.BookResponse
		result2 error
	}
	getBookReturnsOnCall map[int]struct {
		result1 *models.BookResponse
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
	UpdateBookStub        func(*models.BookRequest) (*models.BookResponse, error)
	updateBookMutex       sync.RWMutex
	updateBookArgsForCall []struct {
		arg1 *models.BookRequest
	}
	updateBookReturns struct {
		result1 *models.BookResponse
		result2 error
	}
	updateBookReturnsOnCall map[int]struct {
		result1 *models.BookResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBookerRepository) AddBook(arg1 *models.BookRequest) (int, error) {
	fake.addBookMutex.Lock()
	ret, specificReturn := fake.addBookReturnsOnCall[len(fake.addBookArgsForCall)]
	fake.addBookArgsForCall = append(fake.addBookArgsForCall, struct {
		arg1 *models.BookRequest
	}{arg1})
	stub := fake.AddBookStub
	fakeReturns := fake.addBookReturns
	fake.recordInvocation("AddBook", []interface{}{arg1})
	fake.addBookMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBookerRepository) AddBookCallCount() int {
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	return len(fake.addBookArgsForCall)
}

func (fake *FakeBookerRepository) AddBookCalls(stub func(*models.BookRequest) (int, error)) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = stub
}

func (fake *FakeBookerRepository) AddBookArgsForCall(i int) *models.BookRequest {
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	argsForCall := fake.addBookArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBookerRepository) AddBookReturns(result1 int, result2 error) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = nil
	fake.addBookReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) AddBookReturnsOnCall(i int, result1 int, result2 error) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = nil
	if fake.addBookReturnsOnCall == nil {
		fake.addBookReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.addBookReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) DeleteBook(arg1 int, arg2 int) (int, error) {
	fake.deleteBookMutex.Lock()
	ret, specificReturn := fake.deleteBookReturnsOnCall[len(fake.deleteBookArgsForCall)]
	fake.deleteBookArgsForCall = append(fake.deleteBookArgsForCall, struct {
		arg1 int
		arg2 int
	}{arg1, arg2})
	stub := fake.DeleteBookStub
	fakeReturns := fake.deleteBookReturns
	fake.recordInvocation("DeleteBook", []interface{}{arg1, arg2})
	fake.deleteBookMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBookerRepository) DeleteBookCallCount() int {
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	return len(fake.deleteBookArgsForCall)
}

func (fake *FakeBookerRepository) DeleteBookCalls(stub func(int, int) (int, error)) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = stub
}

func (fake *FakeBookerRepository) DeleteBookArgsForCall(i int) (int, int) {
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	argsForCall := fake.deleteBookArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBookerRepository) DeleteBookReturns(result1 int, result2 error) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = nil
	fake.deleteBookReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) DeleteBookReturnsOnCall(i int, result1 int, result2 error) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = nil
	if fake.deleteBookReturnsOnCall == nil {
		fake.deleteBookReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.deleteBookReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) GetAllBooks() ([]models.BookResponse, error) {
	fake.getAllBooksMutex.Lock()
	ret, specificReturn := fake.getAllBooksReturnsOnCall[len(fake.getAllBooksArgsForCall)]
	fake.getAllBooksArgsForCall = append(fake.getAllBooksArgsForCall, struct {
	}{})
	stub := fake.GetAllBooksStub
	fakeReturns := fake.getAllBooksReturns
	fake.recordInvocation("GetAllBooks", []interface{}{})
	fake.getAllBooksMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBookerRepository) GetAllBooksCallCount() int {
	fake.getAllBooksMutex.RLock()
	defer fake.getAllBooksMutex.RUnlock()
	return len(fake.getAllBooksArgsForCall)
}

func (fake *FakeBookerRepository) GetAllBooksCalls(stub func() ([]models.BookResponse, error)) {
	fake.getAllBooksMutex.Lock()
	defer fake.getAllBooksMutex.Unlock()
	fake.GetAllBooksStub = stub
}

func (fake *FakeBookerRepository) GetAllBooksReturns(result1 []models.BookResponse, result2 error) {
	fake.getAllBooksMutex.Lock()
	defer fake.getAllBooksMutex.Unlock()
	fake.GetAllBooksStub = nil
	fake.getAllBooksReturns = struct {
		result1 []models.BookResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) GetAllBooksReturnsOnCall(i int, result1 []models.BookResponse, result2 error) {
	fake.getAllBooksMutex.Lock()
	defer fake.getAllBooksMutex.Unlock()
	fake.GetAllBooksStub = nil
	if fake.getAllBooksReturnsOnCall == nil {
		fake.getAllBooksReturnsOnCall = make(map[int]struct {
			result1 []models.BookResponse
			result2 error
		})
	}
	fake.getAllBooksReturnsOnCall[i] = struct {
		result1 []models.BookResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) GetBook(arg1 int) (*models.BookResponse, error) {
	fake.getBookMutex.Lock()
	ret, specificReturn := fake.getBookReturnsOnCall[len(fake.getBookArgsForCall)]
	fake.getBookArgsForCall = append(fake.getBookArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.GetBookStub
	fakeReturns := fake.getBookReturns
	fake.recordInvocation("GetBook", []interface{}{arg1})
	fake.getBookMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBookerRepository) GetBookCallCount() int {
	fake.getBookMutex.RLock()
	defer fake.getBookMutex.RUnlock()
	return len(fake.getBookArgsForCall)
}

func (fake *FakeBookerRepository) GetBookCalls(stub func(int) (*models.BookResponse, error)) {
	fake.getBookMutex.Lock()
	defer fake.getBookMutex.Unlock()
	fake.GetBookStub = stub
}

func (fake *FakeBookerRepository) GetBookArgsForCall(i int) int {
	fake.getBookMutex.RLock()
	defer fake.getBookMutex.RUnlock()
	argsForCall := fake.getBookArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBookerRepository) GetBookReturns(result1 *models.BookResponse, result2 error) {
	fake.getBookMutex.Lock()
	defer fake.getBookMutex.Unlock()
	fake.GetBookStub = nil
	fake.getBookReturns = struct {
		result1 *models.BookResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) GetBookReturnsOnCall(i int, result1 *models.BookResponse, result2 error) {
	fake.getBookMutex.Lock()
	defer fake.getBookMutex.Unlock()
	fake.GetBookStub = nil
	if fake.getBookReturnsOnCall == nil {
		fake.getBookReturnsOnCall = make(map[int]struct {
			result1 *models.BookResponse
			result2 error
		})
	}
	fake.getBookReturnsOnCall[i] = struct {
		result1 *models.BookResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) GetDBPool() postgres.DB {
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

func (fake *FakeBookerRepository) GetDBPoolCallCount() int {
	fake.getDBPoolMutex.RLock()
	defer fake.getDBPoolMutex.RUnlock()
	return len(fake.getDBPoolArgsForCall)
}

func (fake *FakeBookerRepository) GetDBPoolCalls(stub func() postgres.DB) {
	fake.getDBPoolMutex.Lock()
	defer fake.getDBPoolMutex.Unlock()
	fake.GetDBPoolStub = stub
}

func (fake *FakeBookerRepository) GetDBPoolReturns(result1 postgres.DB) {
	fake.getDBPoolMutex.Lock()
	defer fake.getDBPoolMutex.Unlock()
	fake.GetDBPoolStub = nil
	fake.getDBPoolReturns = struct {
		result1 postgres.DB
	}{result1}
}

func (fake *FakeBookerRepository) GetDBPoolReturnsOnCall(i int, result1 postgres.DB) {
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

func (fake *FakeBookerRepository) UpdateBook(arg1 *models.BookRequest) (*models.BookResponse, error) {
	fake.updateBookMutex.Lock()
	ret, specificReturn := fake.updateBookReturnsOnCall[len(fake.updateBookArgsForCall)]
	fake.updateBookArgsForCall = append(fake.updateBookArgsForCall, struct {
		arg1 *models.BookRequest
	}{arg1})
	stub := fake.UpdateBookStub
	fakeReturns := fake.updateBookReturns
	fake.recordInvocation("UpdateBook", []interface{}{arg1})
	fake.updateBookMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBookerRepository) UpdateBookCallCount() int {
	fake.updateBookMutex.RLock()
	defer fake.updateBookMutex.RUnlock()
	return len(fake.updateBookArgsForCall)
}

func (fake *FakeBookerRepository) UpdateBookCalls(stub func(*models.BookRequest) (*models.BookResponse, error)) {
	fake.updateBookMutex.Lock()
	defer fake.updateBookMutex.Unlock()
	fake.UpdateBookStub = stub
}

func (fake *FakeBookerRepository) UpdateBookArgsForCall(i int) *models.BookRequest {
	fake.updateBookMutex.RLock()
	defer fake.updateBookMutex.RUnlock()
	argsForCall := fake.updateBookArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBookerRepository) UpdateBookReturns(result1 *models.BookResponse, result2 error) {
	fake.updateBookMutex.Lock()
	defer fake.updateBookMutex.Unlock()
	fake.UpdateBookStub = nil
	fake.updateBookReturns = struct {
		result1 *models.BookResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) UpdateBookReturnsOnCall(i int, result1 *models.BookResponse, result2 error) {
	fake.updateBookMutex.Lock()
	defer fake.updateBookMutex.Unlock()
	fake.UpdateBookStub = nil
	if fake.updateBookReturnsOnCall == nil {
		fake.updateBookReturnsOnCall = make(map[int]struct {
			result1 *models.BookResponse
			result2 error
		})
	}
	fake.updateBookReturnsOnCall[i] = struct {
		result1 *models.BookResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeBookerRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	fake.getAllBooksMutex.RLock()
	defer fake.getAllBooksMutex.RUnlock()
	fake.getBookMutex.RLock()
	defer fake.getBookMutex.RUnlock()
	fake.getDBPoolMutex.RLock()
	defer fake.getDBPoolMutex.RUnlock()
	fake.updateBookMutex.RLock()
	defer fake.updateBookMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBookerRepository) recordInvocation(key string, args []interface{}) {
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

var _ repository.BookerRepository = new(FakeBookerRepository)
