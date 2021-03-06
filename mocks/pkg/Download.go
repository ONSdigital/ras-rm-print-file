// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	pkg "github.com/ONSdigital/ras-rm-print-file/pkg"
	mock "github.com/stretchr/testify/mock"
)

// Download is an autogenerated mock type for the Download type
type Download struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Download) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DownloadFile provides a mock function with given fields: filename
func (_m *Download) DownloadFile(filename string) (*pkg.PrintFile, error) {
	ret := _m.Called(filename)

	var r0 *pkg.PrintFile
	if rf, ok := ret.Get(0).(func(string) *pkg.PrintFile); ok {
		r0 = rf(filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.PrintFile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields:
func (_m *Download) Init() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
