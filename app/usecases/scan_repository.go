package usecases

import "github.com/scanner/app/domain"

// A ScanRepository belong to the usecases layer.
type ScanRepository interface {
	Scan(repoUrl string) (*domain.ScanData, error)
	FindAll() (*domain.ScanResults, error)
	FindByID(int64) (*domain.ScanResult, error)
	Store(*domain.ScanData) (*domain.ScanResult, error)
	Update(*domain.ScanData) (*domain.ScanResult, error)
}
