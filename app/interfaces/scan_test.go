package interfaces_test

import (
	"errors"
	"testing"
	"time"

	"github.com/scanner/app/domain"
	"github.com/scanner/app/domain/mocks"
	"github.com/scanner/app/interfaces"
	"github.com/scanner/app/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test FindAll in repository
func TestScanResultFindAll(t *testing.T) {
	mockScanRepository := new(mocks.ScanRepository)
	timeNowStr := time.Now().UTC().String()
	mockScanResult := domain.ScanResult{
		ID:        1,
		Name:      "Test",
		Url:       "www.test.com/repo",
		QueueTime: timeNowStr,
		StartTime: timeNowStr,
		EndTime:   timeNowStr,
		Result:    `{"findings": [{"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/test.json", "positions": [{"begin": {"cols": ["4"], "line": "2"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}, {"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/README.md", "positions": [{"begin": {"cols": ["4"], "line": "33"}}, {"begin": {"cols": ["10"], "line": "109"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}]}`,
	}
	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}
	mockListScanResult := &domain.ScanResults{}
	*mockListScanResult = append(*mockListScanResult, mockScanResult)

	t.Run("success", func(t *testing.T) {
		mockScanRepository.On("FindAll").Return(mockListScanResult, nil).Once()
		list, err := scanInteractor.Index()
		assert.NoError(t, err)
		assert.Len(t, *list, len(*mockListScanResult))
		mockScanRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockScanRepository.On("FindAll").Return(nil, errors.New("Unexpexted Error")).Once()
		list, err := scanInteractor.Index()
		assert.Error(t, err)
		assert.Nil(t, list)
		mockScanRepository.AssertExpectations(t)
	})

}

// Test FindByID in repository
func TestScanResultFindByID(t *testing.T) {
	mockScanRepository := new(mocks.ScanRepository)
	timeNowStr := time.Now().UTC().String()
	mockScanResult := domain.ScanResult{
		ID:        1,
		Name:      "Test",
		Url:       "www.test.com/repo",
		QueueTime: timeNowStr,
		StartTime: timeNowStr,
		EndTime:   timeNowStr,
		Result:    `{"findings": [{"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/test.json", "positions": [{"begin": {"cols": ["4"], "line": "2"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}, {"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/README.md", "positions": [{"begin": {"cols": ["4"], "line": "33"}}, {"begin": {"cols": ["10"], "line": "109"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}]}`,
	}
	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}
	t.Run("success", func(t *testing.T) {
		mockScanRepository.On("FindByID", mock.AnythingOfType("int64")).Return(&mockScanResult, nil).Once()

		a, err := scanInteractor.Show(mockScanResult.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)
		mockScanRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockScanRepository.On("FindByID", mock.AnythingOfType("int64")).Return(&domain.ScanResult{}, errors.New("Unexpected")).Once()
		a, err := mockScanRepository.FindByID(mockScanResult.ID)
		assert.Error(t, err)
		assert.Equal(t, domain.ScanResult{}, *a)
		mockScanRepository.AssertExpectations(t)

	})

}

// Test Store in repository
func TestScanResultStore(t *testing.T) {
	mockScanRepository := new(mocks.ScanRepository)
	timeNow := time.Now().UTC()
	mockScanData := &domain.ScanData{
		RepoID:    1,
		Status:    2,
		StartTime: timeNow,
		Result:    `{}`,
	}
	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}
	t.Run("success", func(t *testing.T) {
		tempMockScanResult := &domain.ScanResult{
			Name:      "test",
			Url:       "www.test.com/repo",
			StartTime: timeNow.String(),
			Status:    "In Progress",
			Result:    `{}`,
		}

		mockScanRepository.On("Store", mock.AnythingOfType("*domain.ScanData")).Return(tempMockScanResult, nil).Once()
		r, err := scanInteractor.Store(mockScanData)
		assert.NoError(t, err)
		assert.Equal(t, mockScanData.StartTime.String(), r.StartTime)
		assert.Equal(t, mockScanData.Result, r.Result)
		mockScanRepository.AssertExpectations(t)
	})

}

// Test Update in repository
func TestScanResultUpdate(t *testing.T) {
	mockScanRepository := new(mocks.ScanRepository)
	timeNow := time.Now().UTC()
	mockScanData := &domain.ScanData{
		ID:      1,
		RepoID:  1,
		Status:  3,
		EndTime: timeNow,
		Result:  `{"findings": [{"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/test.json", "positions": [{"begin": {"cols": ["4"], "line": "2"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}, {"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/README.md", "positions": [{"begin": {"cols": ["4"], "line": "33"}}, {"begin": {"cols": ["10"], "line": "109"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}]}`,
	}
	scanInteractor := usecases.ScanInteractor{
		ScanRepository: mockScanRepository,
	}
	t.Run("success", func(t *testing.T) {
		tempMockScanResult := &domain.ScanResult{
			Name:      "test",
			Url:       "www.test.com/repo",
			QueueTime: timeNow.String(),
			StartTime: timeNow.String(),
			EndTime:   timeNow.String(),
			Status:    "In Progress",
			Result:    `{"findings": [{"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/test.json", "positions": [{"begin": {"cols": ["4"], "line": "2"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}, {"type": "sast", "ruleId": "G402", "location": {"path": "/tmp/repo3236763213/README.md", "positions": [{"begin": {"cols": ["4"], "line": "33"}}, {"begin": {"cols": ["10"], "line": "109"}}]}, "metadata": {"severity": "HIGH", "description": "Private/Public key is present"}}]}`,
		}
		mockScanRepository.On("Update", mock.AnythingOfType("*domain.ScanData")).Return(tempMockScanResult, nil).Once()
		r, err := scanInteractor.Update(mockScanData)
		assert.NoError(t, err)
		assert.Equal(t, mockScanData.EndTime.String(), r.StartTime)
		assert.Equal(t, mockScanData.Result, r.Result)
		mockScanRepository.AssertExpectations(t)
	})

}

// Test CheckViolationInWord in repository
func TestScanCheckViolationInWord(t *testing.T) {
	scanRepository := &interfaces.ScanRepository{
		SQLHandler:    &mocks.SQLHandler{},
		SearchPattern: []string{"public_key", "prefix_key"},
	}

	// Success case with double quote
	t.Run("successwithquote", func(t *testing.T) {
		words := []string{"test", `"public_key"`}
		results := scanRepository.CheckViolationInWord(words)
		assert.Equal(t, 1, len(results))
	})

	// Success case without double quote
	t.Run("successwithoutquote", func(t *testing.T) {
		words := []string{"test", `public_key`}
		results := scanRepository.CheckViolationInWord(words)
		assert.Equal(t, 1, len(results))
	})

	// Success case without double quote
	t.Run("notfound", func(t *testing.T) {
		words := []string{"test", `testing`}
		results := scanRepository.CheckViolationInWord(words)
		assert.Equal(t, 0, len(results))
	})

}
