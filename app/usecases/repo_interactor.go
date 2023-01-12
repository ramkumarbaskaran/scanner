package usecases

import (
	"github.com/scanner/app/domain"
)

// A RepoInteractor belong to the usecases layer.
type RepoInteractor struct {
	RepoRepository RepoRepository
}

// Index is display a listing of the resource.
func (ui *RepoInteractor) Index() (repos *domain.Repos, err error) {
	return ui.RepoRepository.FindAll()

}

// Show is display the specified resource.
func (ui *RepoInteractor) Show(repoID int64) (repo *domain.Repo, err error) {
	return ui.RepoRepository.FindByID(repoID)
}

// Store is to create new resource.
func (ui *RepoInteractor) Store(repo *domain.Repo) (newRepo *domain.Repo, err error) {
	return ui.RepoRepository.Store(repo)
}

// Update is to update existing resource.
func (ui *RepoInteractor) Update(repo *domain.Repo) (newRepo *domain.Repo, err error) {
	return ui.RepoRepository.Update(repo)
}

// Delete is to soft delete resource.
func (ui *RepoInteractor) Delete(repo *domain.Repo) (newRepo *domain.Repo, err error) {
	return ui.RepoRepository.Delete(repo)
}
