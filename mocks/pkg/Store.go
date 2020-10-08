// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	pkg "github.com/ONSdigital/ras-rm-print-file/pkg"
	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Add provides a mock function with given fields: printFilename, dataFilename
func (_m *Store) Add(printFilename string, dataFilename string) (*pkg.PrintFileRequest, error) {
	ret := _m.Called(printFilename, dataFilename)

	var r0 *pkg.PrintFileRequest
	if rf, ok := ret.Get(0).(func(string, string) *pkg.PrintFileRequest); ok {
		r0 = rf(printFilename, dataFilename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.PrintFileRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(printFilename, dataFilename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindIncomplete provides a mock function with given fields:
func (_m *Store) FindIncomplete() ([]*pkg.PrintFileRequest, error) {
	ret := _m.Called()

	var r0 []*pkg.PrintFileRequest
	if rf, ok := ret.Get(0).(func() []*pkg.PrintFileRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pkg.PrintFileRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields:
func (_m *Store) Init() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: pfr
func (_m *Store) Update(pfr *pkg.PrintFileRequest) error {
	ret := _m.Called(pfr)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pkg.PrintFileRequest) error); ok {
		r0 = rf(pfr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
