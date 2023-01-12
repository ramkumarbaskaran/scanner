package interfaces

import (
	"github.com/scanner/app/domain"
)

// A RepoRepository belong to the inteface layer
type RepoRepository struct {
	SQLHandler SQLHandler
}

// FindAll returns the number of entities.
func (rr *RepoRepository) FindAll() (repos *domain.Repos, err error) {
	const query = `
		SELECT
			id,
			name,
			url
		FROM
			repositories
		WHERE status = 1
	`
	rows, err := rr.SQLHandler.Query(query)

	defer rows.Close()

	if err != nil {
		return
	}
	repos = &domain.Repos{}
	for rows.Next() {
		var id int64
		var name string
		var url string
		if err = rows.Scan(&id, &name, &url); err != nil {
			return
		}
		repo := domain.Repo{
			ID:   id,
			Name: name,
			Url:  url,
		}
		*repos = append(*repos, repo)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

// FindByID returns the entity identified by the given id.
func (rr *RepoRepository) FindByID(repoID int64) (repo *domain.Repo, err error) {
	const query = `
		SELECT
			id,
			name,
			url
		FROM
			repositories
		WHERE
			id = ?
		AND status = 1
	`
	row, err := rr.SQLHandler.Query(query, repoID)

	defer row.Close()

	if err != nil {
		return
	}

	var id int64
	var name string
	var url string
	if !row.Next() {
		return
	}
	if err = row.Scan(&id, &name, &url); err != nil {
		return
	}
	repo = &domain.Repo{
		ID:   id,
		Name: name,
		Url:  url,
	}

	return
}

// Store is to create the new entity.
func (rr *RepoRepository) Store(repo *domain.Repo) (newRepo *domain.Repo, err error) {
	query := `
		INSERT INTO repositories (
			name,
			url,
			created_time,
			updated_time
		)
		VALUES (
			?,
			?,
			?,
			?
		)
	`
	var row Result
	row, err = rr.SQLHandler.Exec(query, repo.Name, repo.Url, repo.CreatedTime, repo.UpdatedTime)
	if err != nil {
		return
	}
	var id int64
	id, err = row.LastInsertId()
	if err != nil {
		return
	}
	newRepo = &domain.Repo{
		ID:          id,
		Name:        repo.Name,
		Url:         repo.Url,
		CreatedTime: repo.CreatedTime,
		UpdatedTime: repo.UpdatedTime,
	}

	return
}

// Update is to update the existing entity.
func (rr *RepoRepository) Update(repo *domain.Repo) (updRepo *domain.Repo, err error) {
	query := `
		UPDATE repositories 
		SET
			name = ?,
			url = ?,
			updated_time = ?
		WHERE
			id = ?
	`

	_, err = rr.SQLHandler.Exec(query, repo.Name, repo.Url, repo.UpdatedTime, repo.ID)
	if err != nil {
		return
	}
	updRepo = &domain.Repo{
		ID:          repo.ID,
		Name:        repo.Name,
		Url:         repo.Url,
		UpdatedTime: repo.UpdatedTime,
	}

	return
}

// Delete is to soft delete the existing entity.
func (rr *RepoRepository) Delete(repo *domain.Repo) (updRepo *domain.Repo, err error) {
	query := `
		UPDATE repositories 
		SET
			status = 0
		WHERE
			id = ?
	`

	_, err = rr.SQLHandler.Exec(query, repo.ID)
	if err != nil {
		return
	}
	updRepo = &domain.Repo{
		ID: repo.ID,
	}

	return
}
