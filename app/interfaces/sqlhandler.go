package interfaces

// A SQLHandler belong to the inteface layer.
type SQLHandler interface {
	Query(string, ...interface{}) (Row, error)
	Exec(string, ...interface{}) (Result, error)
}

// A Result belong to the inteface layer.
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// A Row belong to the inteface layer.
type Row interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
	Err() error
}
