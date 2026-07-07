package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Ensure SQLiteStore implements Store.
// var _ Store = (*SQLiteStore)(nil)

// SQLiteStore implements Store using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens a SQLite database.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}