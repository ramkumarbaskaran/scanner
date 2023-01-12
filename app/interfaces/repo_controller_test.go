package interfaces_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/scanner/app/domain"
	"github.com/scanner/app/domain/mocks"
	"github.com/scanner/app/interfaces"
	"github.com/scanner/app/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Test Index endpoint
func TestRepoControllerIndex(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/repos", strings.NewReader(""))
	assert.NoError(t, err)
	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockListRepo := &domain.Repos{}
	mockRepoRepository.On("FindAll").Return(mockListRepo, nil).Once()
	repoController.Index(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// Test Show endpoint
func TestRepoControllerShow(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/repo/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("repoID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test",
		Url:  "www.test.com/repo",
	}
	mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockRepo, nil).Once()
	repoController.Show(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepoRepository.AssertExpectations(t)
}

// Test Store success endpoint
func TestRepoControllerStore(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/repo", strings.NewReader("name=Test1&url=www.test.com/repo"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	assert.NoError(t, err)
	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockRepo := &domain.Repo{
		Name: "Test1",
		Url:  "www.test.com/repo",
	}
	mockRepoRepository.On("Store", mock.AnythingOfType("*domain.Repo")).Return(mockRepo, nil).Once()
	repoController.Create(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepoRepository.AssertExpectations(t)
}

// Test Store bad request endpoint
func TestRepoControllerStoreError(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/repo", strings.NewReader("name=Te&url=www.test.com/repo"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	assert.NoError(t, err)
	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockRepo := &domain.Repo{
		Name: "Te",
		Url:  "www.test.com/repo",
	}
	mockRepoRepository.On("Store", mock.AnythingOfType("*domain.Repo")).Return(mockRepo, nil).Once()
	repoController.Create(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test Update success endpoint
func TestRepoControllerUpdate(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/api/repo/1", strings.NewReader("name=Test1&url=www.test.com/repo"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("repoID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test1",
		Url:  "www.test.com/repo",
	}
	mockExistRepo := &domain.Repo{
		ID:   1,
		Name: "Test",
		Url:  "www.test.com/repo",
	}
	mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockExistRepo, nil).Once()
	mockRepoRepository.On("Update", mock.AnythingOfType("*domain.Repo")).Return(mockRepo, nil).Once()
	repoController.Update(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepoRepository.AssertExpectations(t)
}

// Test Update bad request endpoint
func TestRepoControllerUpdateError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/api/repo/1", strings.NewReader("name=Test1&url=www.test.com/repo"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("repoID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test1",
		Url:  "www.test.com/repo",
	}

	var mockExistRepo *domain.Repo
	mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockExistRepo, nil).Once()
	mockRepoRepository.On("Update", mock.AnythingOfType("*domain.Repo")).Return(mockRepo, nil).Once()
	repoController.Update(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test Delete success endpoint
func TestRepoControllerDelete(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/repo/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("repoID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test1",
		Url:  "www.test.com/repo",
	}
	mockExistRepo := &domain.Repo{
		ID:   1,
		Name: "Test",
		Url:  "www.test.com/repo",
	}
	mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockExistRepo, nil).Once()
	mockRepoRepository.On("Delete", mock.AnythingOfType("*domain.Repo")).Return(mockRepo, nil).Once()
	repoController.Delete(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepoRepository.AssertExpectations(t)
}

// Test Delete bad request endpoint
func TestRepoControllerDeleteError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/repo/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("repoID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	mockRepoRepository := new(mocks.RepoRepository)
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	repoController := interfaces.RepoController{
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	mockRepo := &domain.Repo{
		ID:   1,
		Name: "Test1",
		Url:  "www.test.com/repo",
	}
	var mockExistRepo *domain.Repo
	mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockExistRepo, nil).Once()
	mockRepoRepository.On("Delete", mock.AnythingOfType("*domain.Repo")).Return(mockRepo, nil).Once()
	repoController.Delete(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
