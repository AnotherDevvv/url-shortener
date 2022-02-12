// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"sync"
	"urlShortener/internal/db"
)

// Ensure, that RepositoryMock does implement db.Repository.
// If this is not the case, regenerate this file with moq.
var _ db.Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of db.Repository.
//
// 	func TestSomethingThatUsesRepository(t *testing.T) {
//
// 		// make and configure a mocked db.Repository
// 		mockedRepository := &RepositoryMock{
// 			CloseFunc: func()  {
// 				panic("mock out the Close method")
// 			},
// 			GetFunc: func(key string) (string, error) {
// 				panic("mock out the Get method")
// 			},
// 			InsertFunc: func(key string, value string) error {
// 				panic("mock out the Insert method")
// 			},
// 		}
//
// 		// use mockedRepository in code that requires db.Repository
// 		// and then make assertions.
//
// 	}
type RepositoryMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func()

	// GetFunc mocks the Get method.
	GetFunc func(key string) (string, error)

	// InsertFunc mocks the Insert method.
	InsertFunc func(key string, value string) error

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Key is the key argument value.
			Key string
		}
		// Insert holds details about calls to the Insert method.
		Insert []struct {
			// Key is the key argument value.
			Key string
			// Value is the value argument value.
			Value string
		}
	}
	lockClose  sync.RWMutex
	lockGet    sync.RWMutex
	lockInsert sync.RWMutex
}

// Close calls CloseFunc.
func (mock *RepositoryMock) Close() {
	if mock.CloseFunc == nil {
		panic("RepositoryMock.CloseFunc: method is nil but Repository.Close was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedRepository.CloseCalls())
func (mock *RepositoryMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *RepositoryMock) Get(key string) (string, error) {
	if mock.GetFunc == nil {
		panic("RepositoryMock.GetFunc: method is nil but Repository.Get was just called")
	}
	callInfo := struct {
		Key string
	}{
		Key: key,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(key)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedRepository.GetCalls())
func (mock *RepositoryMock) GetCalls() []struct {
	Key string
} {
	var calls []struct {
		Key string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Insert calls InsertFunc.
func (mock *RepositoryMock) Insert(key string, value string) error {
	if mock.InsertFunc == nil {
		panic("RepositoryMock.InsertFunc: method is nil but Repository.Insert was just called")
	}
	callInfo := struct {
		Key   string
		Value string
	}{
		Key:   key,
		Value: value,
	}
	mock.lockInsert.Lock()
	mock.calls.Insert = append(mock.calls.Insert, callInfo)
	mock.lockInsert.Unlock()
	return mock.InsertFunc(key, value)
}

// InsertCalls gets all the calls that were made to Insert.
// Check the length with:
//     len(mockedRepository.InsertCalls())
func (mock *RepositoryMock) InsertCalls() []struct {
	Key   string
	Value string
} {
	var calls []struct {
		Key   string
		Value string
	}
	mock.lockInsert.RLock()
	calls = mock.calls.Insert
	mock.lockInsert.RUnlock()
	return calls
}
