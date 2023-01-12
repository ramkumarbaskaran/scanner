package domain

import "time"

// A Repos belong to the domain layer.
type Repos []Repo

// A Repo belong to the domain layer.
type Repo struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	CreatedTime *time.Time `json:"created_time,omitempty"`
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
	Status      int8       `json:"status,omitempty"`
}
