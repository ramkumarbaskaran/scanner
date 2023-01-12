/*package interfaces_test

import (
	"errors"
	"testing"

	"github.com/scanner/app/domain"
	"github.com/scanner/app/domain/mocks"
	"github.com/scanner/app/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepoFindAll(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := domain.Repo{
		Name: "Test",
		Url:  "www.test.com/repo",
	}

	var mockListRepo domain.Repos
	mockListRepo = append(mockListRepo, mockRepo)
	t.Run("success", func(t *testing.T) {
		mockRepoRepository.On("FindAll").Return(mockListRepo, nil).Once()
		list, err := mockRepoRepository.Index()
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListRepo))
		mockRepoRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRepoRepository.On("FindAll").Return(nil, errors.New("Unexpexted Error")).Once()
		list, err := mockRepoRepository.FindAll()
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockRepoRepository.AssertExpectations(t)
	})

}

func TestRepoFindByID(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := domain.Repo{
		ID:   1,
		Name: "Test",
		Url:  "www.test.com/repo",
	}

	t.Run("success", func(t *testing.T) {
		mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockRepo, nil).Once()

		a, err := mockRepoRepository.FindByID(mockRepo.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)
		mockRepoRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(domain.Repo{}, errors.New("Unexpected")).Once()

		a, err := mockRepoRepository.FindByID(mockRepo.ID)
		assert.Error(t, err)
		assert.Equal(t, domain.Repo{}, a)
		mockRepoRepository.AssertExpectations(t)

	})

}

func TestRepoStore(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := domain.Repo{
		Name: "Test1",
		Url:  "www.test.com/repo",
	}

	t.Run("success", func(t *testing.T) {
		tempMockRepo := mockRepo
		tempMockRepo.ID = 0
		mockRepoRepository.On("Store", mock.AnythingOfType("domain.Repo")).Return(tempMockRepo, nil).Once()
		r, err := mockRepoRepository.Store(mockRepo)
		assert.NoError(t, err)
		assert.Equal(t, mockRepo.Name, r.Name)
		mockRepoRepository.AssertExpectations(t)
	})

	mockErrRepo := domain.Repo{
		Name: "Te1",
		Url:  "rest",
	}
	t.Run("error-failed", func(t *testing.T) {
		mockRepoRepository.On("Store", mock.AnythingOfType("domain.Repo")).Return(domain.Repo{}, errors.New("Unexpected")).Once()
		a, err := mockRepoRepository.Store(mockErrRepo)
		assert.Error(t, err)
		assert.Equal(t, domain.Repo{}, a)
		mockRepoRepository.AssertExpectations(t)

	})

}

func TestRepoUpdate(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := domain.Repo{
		ID:   1,
		Name: "Test2",
		Url:  "www.test.com/repo1",
	}

	t.Run("success", func(t *testing.T) {
		tempMockRepo := mockRepo
		tempMockRepo.ID = 0
		mockRepoRepository.On("Update", mock.AnythingOfType("domain.Repo")).Return(tempMockRepo, nil).Once()
		r, err := mockRepoRepository.Update(mockRepo)
		assert.NoError(t, err)
		assert.Equal(t, mockRepo.Name, r.Name)
		mockRepoRepository.AssertExpectations(t)
	})

}

func TestRepoDelete(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := domain.Repo{
		ID:   1,
		Name: "Test2",
		Url:  "www.test.com/repo1",
	}

	t.Run("success", func(t *testing.T) {
		tempMockRepo := mockRepo
		tempMockRepo.ID = 0
		mockRepoRepository.On("Delete", mock.AnythingOfType("domain.Repo")).Return(tempMockRepo, nil).Once()
		r, err := mockRepoRepository.Delete(mockRepo)
		assert.NoError(t, err)
		assert.Equal(t, mockRepo.Name, r.Name)
		mockRepoRepository.AssertExpectations(t)
	})

}
*/

package interfaces_test

import (
	"errors"
	"testing"

	"github.com/scanner/app/domain"
	"github.com/scanner/app/domain/mocks"
	"github.com/scanner/app/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test FindAll in repository
func TestRepoFindAll(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := domain.Repo{
		Name: "Test",
		Url:  "www.test.com/repo",
	}

	var mockListRepo domain.Repos
	mockListRepo = append(mockListRepo, mockRepo)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	t.Run("success", func(t *testing.T) {
		mockRepoRepository.On("FindAll").Return(&mockListRepo, nil).Once()
		list, err := repoInteractor.Index()
		assert.NoError(t, err)
		assert.Len(t, *list, len(mockListRepo))
		mockRepoRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRepoRepository.On("FindAll").Return(nil, errors.New("Unexpexted Error")).Once()
		list, err := repoInteractor.Index()
		assert.Error(t, err)
		assert.Nil(t, list)
		mockRepoRepository.AssertExpectations(t)
	})

}

// Test FindByID in repository
func TestRepoFindByID(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test",
		Url:  "www.test.com/repo",
	}
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	t.Run("success", func(t *testing.T) {
		mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockRepo, nil).Once()
		a, err := repoInteractor.Show(mockRepo.ID)
		assert.NoError(t, err)
		assert.NotNil(t, a)
		mockRepoRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(&domain.Repo{}, errors.New("Unexpected")).Once()
		a, err := repoInteractor.Show(mockRepo.ID)
		assert.Error(t, err)
		assert.Equal(t, &domain.Repo{}, a)
		mockRepoRepository.AssertExpectations(t)

	})

}

// Test Store in repository
func TestRepoStore(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := &domain.Repo{
		Name: "Test1",
		Url:  "www.test.com/repo",
	}
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	t.Run("success", func(t *testing.T) {
		tempMockRepo := mockRepo
		tempMockRepo.ID = 0
		mockRepoRepository.On("Store", mock.AnythingOfType("*domain.Repo")).Return(tempMockRepo, nil).Once()
		r, err := repoInteractor.Store(mockRepo)
		assert.NoError(t, err)
		assert.Equal(t, mockRepo.Name, r.Name)
		mockRepoRepository.AssertExpectations(t)
	})

	mockErrRepo := &domain.Repo{
		Name: "Te1",
		Url:  "rest",
	}
	t.Run("error-failed", func(t *testing.T) {
		mockRepoRepository.On("Store", mock.AnythingOfType("*domain.Repo")).Return(&domain.Repo{}, errors.New("Unexpected")).Once()
		a, err := repoInteractor.Store(mockErrRepo)
		assert.Error(t, err)
		assert.Equal(t, &domain.Repo{}, a)
		mockRepoRepository.AssertExpectations(t)

	})

}

// Test Update in repository
func TestRepoUpdate(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test2",
		Url:  "www.test.com/repo1",
	}
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	t.Run("success", func(t *testing.T) {
		tempMockRepo := mockRepo
		tempMockRepo.ID = 0
		mockRepoRepository.On("Update", mock.AnythingOfType("*domain.Repo")).Return(tempMockRepo, nil).Once()
		r, err := repoInteractor.Update(mockRepo)
		assert.NoError(t, err)
		assert.Equal(t, mockRepo.Name, r.Name)
		mockRepoRepository.AssertExpectations(t)
	})

}

// Test Delete in repository
func TestRepoDelete(t *testing.T) {
	mockRepoRepository := new(mocks.RepoRepository)
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test2",
		Url:  "www.test.com/repo1",
	}
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	t.Run("success", func(t *testing.T) {
		tempMockRepo := mockRepo
		tempMockRepo.ID = 0
		mockRepoRepository.On("Delete", mock.AnythingOfType("*domain.Repo")).Return(tempMockRepo, nil).Once()
		r, err := repoInteractor.Delete(mockRepo)
		assert.NoError(t, err)
		assert.Equal(t, mockRepo.Name, r.Name)
		mockRepoRepository.AssertExpectations(t)
	})

}
