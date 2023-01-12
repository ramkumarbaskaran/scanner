package mocks

import (
	domain "github.com/scanner/app/domain"
	mock "github.com/stretchr/testify/mock"
)

// Mocking Repo Repository
type RepoRepository struct {
	mock.Mock
}

// FindAll provides a mock function
func (_m *RepoRepository) FindAll() (*domain.Repos, error) {
	ret := _m.Called()

	var r0 *domain.Repos
	if rf, ok := ret.Get(0).(func() *domain.Repos); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Repos)
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
func (_m *RepoRepository) FindByID(id int64) (*domain.Repo, error) {
	ret := _m.Called(id)

	var r0 *domain.Repo
	if rf, ok := ret.Get(0).(func(int64) *domain.Repo); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(*domain.Repo)
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
func (_m *RepoRepository) Store(_a0 *domain.Repo) (newRepo *domain.Repo, err error) {
	ret := _m.Called(_a0)

	var r0 *domain.Repo
	if rf, ok := ret.Get(0).(func(*domain.Repo) *domain.Repo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(*domain.Repo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.Repo) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *RepoRepository) Update(_a0 *domain.Repo) (newRepo *domain.Repo, err error) {
	ret := _m.Called(_a0)

	var r0 *domain.Repo
	if rf, ok := ret.Get(0).(func(*domain.Repo) *domain.Repo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(*domain.Repo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.Repo) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0
func (_m *RepoRepository) Delete(_a0 *domain.Repo) (newRepo *domain.Repo, err error) {
	ret := _m.Called(_a0)

	var r0 *domain.Repo
	if rf, ok := ret.Get(0).(func(*domain.Repo) *domain.Repo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(*domain.Repo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.Repo) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
