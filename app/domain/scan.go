package domain

import "time"

// A ScanResults belong to the domain layer.
type ScanResults []ScanResult

// A ScanData belong to the domain layer.
type ScanData struct {
	ID        int64
	RepoID    int64
	Status    int8
	QueueTime time.Time
	StartTime time.Time
	EndTime   time.Time
	Result    string
}

// A ScanData belong to the domain layer.
type ScanResult struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Url       string `json:"url,omitempty"`
	QueueTime string `json:"queue_time,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Status    string `json:"status,omitempty"`
	Result    string `json:"result,omitempty"`
}
