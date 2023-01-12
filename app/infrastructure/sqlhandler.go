package infrastructure

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/scanner/app/interfaces"
)

// A SQLHandler belong to the infrastructure layer.
type SQLHandler struct {
	Conn *sql.DB
}

// A Result belong to the infrastructure layer.
type Result struct {
	Result sql.Result
}

// A Row belong to the infrastructure layer.
type Row struct {
	Rows *sql.Rows
}

// NewSQLHandler returns connection and methods which is related to database handling.
func NewSQLHandler() (interfaces.SQLHandler, error) {
	sqlHandler := &SQLHandler{}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	conn, err := sql.Open(os.Getenv("DB_DRIVER"), dataSourceName)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	sqlHandler.Conn = conn

	return sqlHandler, nil
}

// Query returns results of a Query method.
func (s *SQLHandler) Query(query string, args ...interface{}) (interfaces.Row, error) {
	rows, err := s.Conn.Query(query, args...)

	if err != nil {
		return nil, err
	}

	row := &Row{}
	row.Rows = rows

	return row, nil
}

// Exec is execute statement
func (s *SQLHandler) Exec(query string, args ...interface{}) (interfaces.Result, error) {
	result, err := s.Conn.Exec(query, args...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// LastInsertId returns results of a LastInsertId method.
func (r Result) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

// RowsAffected returns results of a RowsAffected method.
func (r Result) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

// Scan returns results of a Scan method.
func (r Row) Scan(value ...interface{}) error {
	return r.Rows.Scan(value...)
}

// Next returns results of a Next method.
func (r Row) Next() bool {
	return r.Rows.Next()
}

// Close returns results of a Close method.
func (r Row) Close() error {
	return r.Rows.Close()
}

// Err returns results of a Err method.
func (r Row) Err() error {
	return r.Rows.Err()
}
