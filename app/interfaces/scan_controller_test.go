package interfaces_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
func TestScanIndex(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/scan/results", strings.NewReader(""))
	assert.NoError(t, err)
	mockScanRepository := new(mocks.ScanRepository)
	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}
	scanController := interfaces.ScanController{
		ScanInteractor: scanInteractor,
		Logger:         zap.NewNop(),
	}
	mockListScanResult := &domain.ScanResults{}
	mockScanRepository.On("FindAll").Return(mockListScanResult, nil).Once()
	scanController.Index(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// Test Show endpoint
func TestScanShow(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/scan/result/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("resultID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	mockScanRepository := new(mocks.ScanRepository)

	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}

	scanController := interfaces.ScanController{
		ScanInteractor: scanInteractor,
		Logger:         zap.NewNop(),
	}
	timeNowStr := time.Now().UTC().String()
	mockScanResult := &domain.ScanResult{
		ID:        1,
		Name:      "Test",
		Url:       "www.test.com/repo",
		QueueTime: timeNowStr,
		StartTime: timeNowStr,
		EndTime:   timeNowStr,
		Result:    `{}`,
	}
	mockScanRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockScanResult, nil).Once()
	scanController.Show(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockScanRepository.AssertExpectations(t)
}

// Test Show not found endpoint
func TestScanShowError(t *testing.T) {

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/scan/result/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("resultID", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockScanRepository := new(mocks.ScanRepository)
	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}
	scanController := interfaces.ScanController{
		ScanInteractor: scanInteractor,
		Logger:         zap.NewNop(),
	}
	var mockScanResult *domain.ScanResult
	mockScanRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockScanResult, nil).Once()
	scanController.Show(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// Test Show not found endpoint
func TestScan(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/repo/1/scan", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("repoID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockScanRepository := new(mocks.ScanRepository)
	mockRepoRepository := new(mocks.RepoRepository)
	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}
	repoInteractor := usecases.RepoInteractor{
		RepoRepository: mockRepoRepository,
	}
	scanController := interfaces.ScanController{
		ScanInteractor: scanInteractor,
		RepoInteractor: repoInteractor,
		Logger:         zap.NewNop(),
	}
	timeNow := time.Now().UTC()
	mockScanData := &domain.ScanData{
		RepoID:  1,
		Status:  3,
		EndTime: timeNow,
		Result:  `{}`,
	}
	mockScanResult := &domain.ScanResult{
		Name:      "test",
		Url:       "www.test.com/repo",
		StartTime: timeNow.String(),
		Status:    "In Progress",
		Result:    `{}`,
	}
	mockScanUpdateResult := &domain.ScanResult{
		Name:    "test",
		Url:     "www.test.com/repo",
		EndTime: timeNow.String(),
		Status:  "Success",
		Result:  `{}`,
	}
	mockExistRepo := &domain.Repo{
		ID:   1,
		Name: "Test",
		Url:  "www.test.com/repo",
	}

	mockScanRepository.On("Store", mock.AnythingOfType("*domain.ScanData")).Return(mockScanResult, nil).Once()
	mockRepoRepository.On("FindByID", mock.AnythingOfType("int64")).Return(mockExistRepo, nil).Once()
	mockScanRepository.On("Scan", mock.AnythingOfType("string")).Return(mockScanData, nil).Once()
	mockScanRepository.On("Update", mock.AnythingOfType("*domain.ScanData")).Return(mockScanUpdateResult, nil).Once()
	scanController.Scan(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
