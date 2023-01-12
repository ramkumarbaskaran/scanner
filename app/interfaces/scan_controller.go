package interfaces

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/scanner/app/domain"
	"github.com/scanner/app/helper"
	"github.com/scanner/app/usecases"
	"go.uber.org/zap"
)

// ScanController belong to the interface layer.
type ScanController struct {
	ScanInteractor usecases.ScanInteractor
	RepoInteractor usecases.RepoInteractor
	Logger         *zap.Logger
}

// NewScanController returns the instance of Scan controller.
func NewScanController(sqlHandler SQLHandler, logger *zap.Logger) *ScanController {
	// Get search patterns from env file
	searchPattern := os.Getenv("SEARCH_PATTERN")
	searchPatterns := strings.Split(searchPattern, ",")
	noOfWorkers, err := strconv.Atoi(os.Getenv("NOOFWORKERS"))

	// If NOOFWORKERS is invalid input, we can have default worker count as 1
	if err != nil {
		noOfWorkers = 1
	}

	return &ScanController{
		ScanInteractor: usecases.ScanInteractor{
			ScanRepository: &ScanRepository{
				SQLHandler:            sqlHandler,
				SearchPattern:         searchPatterns,
				ScanCloneFolder:       os.Getenv("SCANClONEFOLDER"),
				ScanCloneFolderPrefix: os.Getenv("SCANClONEFOLDERPREFIX"),
				NoOfWorkers:           noOfWorkers,
			},
		},
		RepoInteractor: usecases.RepoInteractor{
			RepoRepository: &RepoRepository{
				SQLHandler: sqlHandler,
			},
		},
		Logger: logger,
	}
}

// Scan the repository.
func (sc *ScanController) Scan(w http.ResponseWriter, r *http.Request) {
	sc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
	repoID, err := strconv.ParseInt(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	repo, err := sc.RepoInteractor.Show(repoID)
	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, err)
		return
	}

	// If No repo found
	if repo == nil {
		err := errors.New("no repos found")
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	// Need to write code if repo exist or not
	scanData := &domain.ScanData{
		RepoID:    repoID,
		Status:    2,
		Result:    `{}`,
		QueueTime: time.Now().UTC(),
		StartTime: time.Now().UTC(),
	}
	scanResult, err := sc.ScanInteractor.Store(scanData)
	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	updatedScanData, err := sc.ScanInteractor.Scan(repo.Url)

	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// If result is empty, store empty json in db
	if updatedScanData.Result == "" {
		updatedScanData.Result = "{}"
	}
	updatedScanData.ID = scanResult.ID
	updatedScanData.RepoID = repoID
	scanResult, err = sc.ScanInteractor.Update(updatedScanData)
	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	helper.Write(w, http.StatusOK, scanResult)
}

// Index return response which contain a listing of the resource of scan results.
func (sc *ScanController) Index(w http.ResponseWriter, r *http.Request) {
	sc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
	results, err := sc.ScanInteractor.Index()
	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	helper.Write(w, http.StatusOK, results)
}

// Show return response which contain the specified resource of a scan result.
func (sc *ScanController) Show(w http.ResponseWriter, r *http.Request) {
	sc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
	resultID, err := strconv.ParseInt(chi.URLParam(r, "resultID"), 10, 64)
	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	scanResult, err := sc.ScanInteractor.Show(resultID)
	if err != nil {
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// If No repo found
	if scanResult == nil {
		err = errors.New("no result found")
		sc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	helper.Write(w, http.StatusOK, scanResult)
}
