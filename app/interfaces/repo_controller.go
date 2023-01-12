package interfaces

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/scanner/app/domain"
	"github.com/scanner/app/helper"
	"github.com/scanner/app/usecases"
	"go.uber.org/zap"
)

// A RepoController belong to the interface layer.
type RepoController struct {
	RepoInteractor usecases.RepoInteractor
	Logger         *zap.Logger
}

// Struct for validation
type RepoData struct {
	name string
	url  string
}

// NewRepoController create new instance of repo.
func NewRepoController(sqlHandler SQLHandler, logger *zap.Logger) *RepoController {
	return &RepoController{
		RepoInteractor: usecases.RepoInteractor{
			RepoRepository: &RepoRepository{
				SQLHandler: sqlHandler,
			},
		},
		Logger: logger,
	}
}

// Index return response which contain a listing of the resource of repos.
func (rc *RepoController) Index(w http.ResponseWriter, r *http.Request) {
	rc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))

	repos, err := rc.RepoInteractor.Index()
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, err)
	}
	helper.Write(w, http.StatusOK, repos)
}

// Show return response which contain the specified resource of a repo.
func (rc *RepoController) Show(w http.ResponseWriter, r *http.Request) {
	rc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
	repoID, err := strconv.ParseInt(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	repo, err := rc.RepoInteractor.Show(repoID)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, err)
		return
	}

	// If No repo found
	if repo == nil {
		err := errors.New("no repos found")
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	helper.Write(w, http.StatusOK, repo)
}

// Create return response which contain the resource of a newly created repo.
func (rc *RepoController) Create(w http.ResponseWriter, r *http.Request) {
	rc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
	repoName := r.PostFormValue("name")
	repoUrl := r.PostFormValue("url")
	repoData := &RepoData{
		name: repoName,
		url:  repoUrl,
	}
	err := rc.validate(repoData)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	timeNow := time.Now().UTC()
	repo := &domain.Repo{
		Name:        repoName,
		Url:         repoUrl,
		CreatedTime: &timeNow,
		UpdatedTime: &timeNow,
	}
	newRepo, err := rc.RepoInteractor.Store(repo)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.Write(w, http.StatusOK, newRepo)
}

// Update existing resource and return updated resource.
func (rc *RepoController) Update(w http.ResponseWriter, r *http.Request) {
	rc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
	repoID, err := strconv.ParseInt(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	repoName := r.PostFormValue("name")
	repoUrl := r.PostFormValue("url")
	repoData := &RepoData{
		name: repoName,
		url:  repoUrl,
	}

	err = rc.validate(repoData)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// If repo does not exist
	isExist, _, err := rc.exists(repoID)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, err)
		return
	}
	if !isExist {
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": "no repo found"})
		return
	}
	timeNow := time.Now().UTC()
	repo := &domain.Repo{
		ID:          repoID,
		Name:        repoName,
		Url:         repoUrl,
		UpdatedTime: &timeNow,
	}
	newRepo, err := rc.RepoInteractor.Update(repo)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		//render
		helper.Write(w, http.StatusInternalServerError, err)
		return
	}

	helper.Write(w, http.StatusOK, newRepo)
}

// Soft Delete existing resource.
func (rc *RepoController) Delete(w http.ResponseWriter, r *http.Request) {
	rc.Logger.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
	repoID, err := strconv.ParseInt(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// If repo does not exist
	isExist, repoName, err := rc.exists(repoID)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, err)
		return
	}

	if !isExist {
		helper.Write(w, http.StatusBadRequest, map[string]string{"error": "no repo found"})
		return
	}

	repo := &domain.Repo{
		ID:     repoID,
		Status: 0,
	}
	// Status will be udpated as 0 for soft delete
	_, err = rc.RepoInteractor.Delete(repo)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("%s", err))
		helper.Write(w, http.StatusInternalServerError, err)
		return
	}

	helper.Write(w, http.StatusOK, map[string]string{"msg": fmt.Sprintf("Repository %s deleted successfully", repoName)})
}

func (c *RepoController) validate(rd *RepoData) error {
	return validation.ValidateStruct(rd,
		// Name cannot be empty, and the length must between 5 and 50
		validation.Field(&rd.name, validation.Required, validation.Length(5, 50)),
		// Url cannot be empty, and should be valid url
		validation.Field(&rd.url, validation.Required, is.URL),
	)
}

func (rc *RepoController) exists(repoID int64) (isExist bool, repoName string, err error) {
	repo, err := rc.RepoInteractor.Show(repoID)
	if err != nil {
		return
	}
	// If repo is found
	if repo != nil {
		isExist = true
		repoName = repo.Name
	}
	return
}
