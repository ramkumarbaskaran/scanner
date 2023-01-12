package mocks

import (
	domain "github.com/scanner/app/domain"
	mock "github.com/stretchr/testify/mock"
)

// Mocking Scan Repository
type ScanRepository struct {
	mock.Mock
}

func (_m *ScanRepository) Scan(repoUrl string) (*domain.ScanData, error) {
	ret := _m.Called(repoUrl)

	var r0 *domain.ScanData
	if rf, ok := ret.Get(0).(func(string) *domain.ScanData); ok {
		r0 = rf(repoUrl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ScanData)
		}
	}

	var r2 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r2 = rf(repoUrl)
	} else {
		r2 = ret.Error(1)
	}

	return r0, r2
}

// FindAll provides a mock function
func (_m *ScanRepository) FindAll() (*domain.ScanResults, error) {
	ret := _m.Called()

	var r0 *domain.ScanResults
	if rf, ok := ret.Get(0).(func() *domain.ScanResults); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ScanResults)
		}
	}

	var r2 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(1)
	}

	return r0, r2
}

// FindByID provides a mock function with given fields: id
func (_m *ScanRepository) FindByID(id int64) (*domain.ScanResult, error) {
	ret := _m.Called(id)

	var r0 *domain.ScanResult
	if rf, ok := ret.Get(0).(func(int64) *domain.ScanResult); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(*domain.ScanResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: _a0
func (_m *ScanRepository) Store(_a0 *domain.ScanData) (newScanResult *domain.ScanResult, err error) {
	ret := _m.Called(_a0)

	var r0 *domain.ScanResult
	if rf, ok := ret.Get(0).(func(*domain.ScanData) *domain.ScanResult); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(*domain.ScanResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.ScanData) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *ScanRepository) Update(_a0 *domain.ScanData) (updScanResult *domain.ScanResult, err error) {
	ret := _m.Called(_a0)

	var r0 *domain.ScanResult
	if rf, ok := ret.Get(0).(func(*domain.ScanData) *domain.ScanResult); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(*domain.ScanResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.ScanData) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
