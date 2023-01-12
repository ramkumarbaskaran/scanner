package usecases

import "github.com/scanner/app/domain"

// A ScanInteractor belong to the usecases layer.
type ScanInteractor struct {
	ScanRepository ScanRepository
}

// Scan the repository.
func (si *ScanInteractor) Scan(repoUrl string) (scanData *domain.ScanData, err error) {
	return si.ScanRepository.Scan(repoUrl)
}

// Index is display a listing of the resource.
func (si *ScanInteractor) Index() (scanResults *domain.ScanResults, err error) {
	return si.ScanRepository.FindAll()
}

// Show is display the specified resource.
func (si *ScanInteractor) Show(resultID int64) (scanResult *domain.ScanResult, err error) {
	return si.ScanRepository.FindByID(resultID)
}

// Store is to create new resource.
func (si *ScanInteractor) Store(scanData *domain.ScanData) (newScanResult *domain.ScanResult, err error) {
	return si.ScanRepository.Store(scanData)
}

// Update is to update existing resource.
func (si *ScanInteractor) Update(scanData *domain.ScanData) (newScanResult *domain.ScanResult, err error) {
	return si.ScanRepository.Update(scanData)
}
