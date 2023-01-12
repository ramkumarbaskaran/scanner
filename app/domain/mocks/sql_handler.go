package mocks

import (
	"github.com/scanner/app/interfaces"
	mock "github.com/stretchr/testify/mock"
)

// Mocking SQLHandler
type SQLHandler struct {
	mock.Mock
}

func (s *SQLHandler) Query(string, ...interface{}) (row interfaces.Row, err error) {
	return
}
func (s *SQLHandler) Exec(string, ...interface{}) (result interfaces.Result, err error) {
	return
}
