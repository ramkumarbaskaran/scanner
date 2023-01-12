package interfaces

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/scanner/app/domain"
)

// A ScanRepository belong to the inteface layer
type ScanRepository struct {
	SQLHandler            SQLHandler
	SearchPattern         []string
	ScanCloneFolder       string
	ScanCloneFolderPrefix string
	NoOfWorkers           int
}

// Struct for scan result
type result struct {
	Findings findings `json:"findings"`
}

type findings []finding

type finding struct {
	ErrorType string   `json:"type"`
	RuleID    string   `json:"ruleId"`
	Location  location `json:"location"`
	Metadata  metadata `json:"metadata"`
}

type metadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type location struct {
	Path      string     `json:"path"`
	Positions []position `json:"positions"`
}

type position struct {
	Begin begin `json:"begin"`
}

type begin struct {
	Line string   `json:"line"`
	Cols []string `json:"cols"`
}

type resultWrapper struct {
	path      string
	positions []position
	err       error
}

type jsonResultWrapper struct {
	output []byte
	err    error
}

// Scan the repository.
func (sr *ScanRepository) Scan(repoUrl string) (scanData *domain.ScanData, err error) {
	// Create temporary directory to clone the repo
	directory, err := ioutil.TempDir(sr.ScanCloneFolder, sr.ScanCloneFolderPrefix)
	if err != nil {
		return
	}
	// Delete the directory after scanning
	defer os.RemoveAll(directory)

	// Cloning in configured folder
	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		URL:               repoUrl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return
	}

	jobs := make(chan string, sr.NoOfWorkers)
	results := make(chan resultWrapper, sr.NoOfWorkers)
	var wg sync.WaitGroup
	// Spawn go routines to scan the security violation
	// each file will be scanned by each worker
	for w := 1; w <= sr.NoOfWorkers; w++ {
		wg.Add(1)
		go sr.worker(w, &wg, jobs, results)
	}
	jsonResult := make(chan jsonResultWrapper)
	go sr.processResults(results, jsonResult)

	// Traverse each file recursively and pass the file to the worker to check violation
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			jobs <- path
		}
		return nil
	})
	close(jobs)

	wg.Wait()
	close(results)
	// Final json output
	jsonResultWrap := <-jsonResult
	output := jsonResultWrap.output
	err = jsonResultWrap.err
	scanData = &domain.ScanData{
		EndTime: time.Now().UTC(),
	}
	// If there is error then we need to update the status as failed
	if err != nil {
		scanData.Status = 4
	} else {
		scanData.Status = 3
		scanData.Result = string(output)
	}
	return
}

// Merge all the results and send it to the caller
func (sr *ScanRepository) processResults(results chan resultWrapper, jsonResult chan jsonResultWrapper) {
	var (
		findingsOutput findings
		resultOutput   result
		err            error
		output         []byte
	)
	// Receiving results from all the channels
	// and merge all of them
	for resultWrapper := range results {
		err := resultWrapper.err
		if err != nil {
			jsonResult <- jsonResultWrapper{err: err}
			return
		}
		positions := resultWrapper.positions
		// If violation is found, add it in the result
		if len(positions) > 0 {
			findingsOutput = append(findingsOutput, finding{
				ErrorType: "sast",
				RuleID:    "G402",
				Location: location{
					Path:      resultWrapper.path,
					Positions: positions,
				},
				Metadata: metadata{
					Description: "Private/Public key is present",
					Severity:    "HIGH",
				},
			})
			resultOutput = result{Findings: findingsOutput}

		}
	}
	output, err = json.MarshalIndent(resultOutput, "", "  ")
	jsonResult <- jsonResultWrapper{output: output, err: err}
	return
}

// Check security violation in the given path
func (sr *ScanRepository) checkViolation(path string) (positions []position, err error) {
	var file *os.File
	file, err = os.Open(path)
	// If any error while opening the file, we need to report
	if err != nil {
		return
	}
	Scanner := bufio.NewScanner(file)
	lineCount := 1
	for Scanner.Scan() {
		// Read line by line
		line := Scanner.Text()
		words := strings.Split(line, " ")
		errCols := sr.CheckViolationInWord(words)
		// If violation is found, then line no and column nos will be added in the result
		if len(errCols) > 0 {
			positions = append(positions, position{Begin: begin{Line: fmt.Sprintf("%d", lineCount), Cols: errCols}})
		}
		lineCount++
	}

	err = Scanner.Err()
	return
}

// Process each word and check if key is found
func (sr *ScanRepository) CheckViolationInWord(words []string) (errCols []string) {
	colCount := 1
	for _, word := range words {
		// Trim double quote and single around the word
		word := strings.Trim(strings.Trim(word, "\""), ",")
		// check if prefix is found in the word
		for _, pattern := range sr.SearchPattern {
			if strings.HasPrefix(word, pattern) {
				errCols = append(errCols, fmt.Sprintf("%d", colCount))
				// If any one pattern is found we can skip
				break
			}
		}
		colCount++
	}
	return
}

// Each worker will process each file
func (sr *ScanRepository) worker(id int, wg *sync.WaitGroup, jobs <-chan string, results chan<- resultWrapper) {
	defer wg.Done()
	for path := range jobs {
		positions, err := sr.checkViolation(path)
		results <- resultWrapper{path: path, positions: positions, err: err}
	}
}

// FindAll returns the number of entities.
func (sr *ScanRepository) FindAll() (scanResults *domain.ScanResults, err error) {
	const query = `
		SELECT
			sr.id,
			r.name,
			r.url,
			sr.status,
			sr.result,
			sr.queue_time,
			sr.start_time,
			sr.end_time
		FROM
			scan_results sr 
		JOIN 
			repositories r 
		ON sr.repo_id = r.id
	`
	rows, err := sr.SQLHandler.Query(query)

	defer rows.Close()

	if err != nil {
		return
	}
	scanResults = &domain.ScanResults{}
	for rows.Next() {
		var (
			id        int64
			name      string
			url       string
			status    string
			result    string
			queueTime sql.NullString
			startTime sql.NullString
			endTime   sql.NullString
		)
		if err = rows.Scan(&id, &name, &url, &status, &result, &queueTime, &startTime, &endTime); err != nil {
			return
		}

		scanResult := domain.ScanResult{
			ID:        id,
			Name:      name,
			Url:       url,
			Status:    status,
			Result:    result,
			QueueTime: queueTime.String,
			StartTime: startTime.String,
			EndTime:   endTime.String,
		}
		*scanResults = append(*scanResults, scanResult)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

// FindByID returns the entity identified by the given id.
func (sr *ScanRepository) FindByID(resultID int64) (scanResult *domain.ScanResult, err error) {
	const query = `
		SELECT
			sr.id,
			r.name,
			r.url,
			sr.status,
			sr.result,
			sr.queue_time,
			sr.start_time,
			sr.end_time
		FROM
			scan_results sr 
		JOIN 
			repositories r 
		ON sr.repo_id = r.id
		WHERE
			sr.id = ?
	`
	row, err := sr.SQLHandler.Query(query, resultID)

	defer row.Close()

	if err != nil {
		return
	}

	var (
		id        int64
		name      string
		url       string
		status    string
		result    string
		queueTime sql.NullString
		startTime sql.NullString
		endTime   sql.NullString
	)
	if !row.Next() {
		return
	}
	if err = row.Scan(&id, &name, &url, &status, &result, &queueTime, &startTime, &endTime); err != nil {
		return
	}

	scanResult = &domain.ScanResult{
		ID:        id,
		Name:      name,
		Url:       url,
		Status:    status,
		Result:    result,
		QueueTime: queueTime.String,
		StartTime: startTime.String,
		EndTime:   endTime.String,
	}

	return
}

// Store is to create the new entity.
func (sr *ScanRepository) Store(scanData *domain.ScanData) (newScanResult *domain.ScanResult, err error) {
	query := `
		INSERT INTO scan_results (
			repo_id,
			status,
			result,
			queue_time,
			start_time
		)
		VALUES (
			?,
			?,
			?,
			?,
			?
		)
	`
	var row Result
	row, err = sr.SQLHandler.Exec(query, scanData.RepoID, scanData.Status, scanData.Result, scanData.QueueTime, scanData.StartTime)
	if err != nil {
		return
	}
	var id int64
	id, err = row.LastInsertId()
	if err != nil {
		return
	}
	newScanResult = &domain.ScanResult{
		ID:     id,
		Status: sr.getStatus(scanData.Status),
		Result: scanData.Result,
	}

	return
}

// Update is to update the existing entity.
func (sr *ScanRepository) Update(scanData *domain.ScanData) (updScanResult *domain.ScanResult, err error) {
	query := `
		UPDATE scan_results 
		SET
			result = ?,
			status = ?,
			end_time = ?
		WHERE
			id = ?
	`

	_, err = sr.SQLHandler.Exec(query, scanData.Result, scanData.Status, scanData.EndTime, scanData.ID)
	if err != nil {
		return
	}
	updScanResult = &domain.ScanResult{
		ID:     scanData.ID,
		Status: sr.getStatus(scanData.Status),
		Result: scanData.Result,
	}

	return
}

// Get the user readble status of scan result
func (sr *ScanRepository) getStatus(status int8) (statusStr string) {
	switch status {
	case 1:
		statusStr = "Queued"
	case 2:
		statusStr = "In Progress"
	case 3:
		statusStr = "Success"
	case 4:
		statusStr = "Failure"
	}
	return
}
