package usecases

import (
	"github.com/scanner/app/domain"
)

// A RepoRepository belong to the usecases layer.
type RepoRepository interface {
	FindAll() (*domain.Repos, error)
	FindByID(int64) (*domain.Repo, error)
	Store(*domain.Repo) (*domain.Repo, error)
	Update(*domain.Repo) (*domain.Repo, error)
	Delete(*domain.Repo) (*domain.Repo, error)
}
