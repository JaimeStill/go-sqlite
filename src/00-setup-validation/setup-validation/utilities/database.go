// Package utilities provides FTS5 testing utilities for setup validation
package utilities

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver with FTS5 support
)

// ConnectFTS5 creates a new SQLite connection with FTS5 enabled
func ConnectFTS5(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify FTS5 is available
	if err := verifyFTS5Support(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("FTS5 not available: %w", err)
	}

	return db, nil
}

// ConnectMemoryFTS5 creates an in-memory SQLite database with FTS5
func ConnectMemoryFTS5() (*sql.DB, error) {
	return ConnectFTS5(":memory:")
}

// verifyFTS5Support checks if FTS5 is available in the SQLite build
func verifyFTS5Support(db *sql.DB) error {
	// Try to create a temporary FTS5 table to verify support
	_, err := db.Exec("CREATE VIRTUAL TABLE temp.fts5_test USING fts5(content)")
	if err != nil {
		return fmt.Errorf("SQLite was not compiled with FTS5 support: %w", err)
	}
	
	// Clean up the test table
	_, _ = db.Exec("DROP TABLE temp.fts5_test")
	
	return nil
}