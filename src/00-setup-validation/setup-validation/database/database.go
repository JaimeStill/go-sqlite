package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/errors"
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/utilities"
	_ "github.com/mattn/go-sqlite3"
)

// Instance is the global database instance
var Instance *Database

// Init initializes the global database connection
func Init(dataSourceName string) error {
	db, err := NewDatabase(dataSourceName)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}
	Instance = db
	return nil
}

// Database wraps the SQL database connection with FTS5-specific operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := utilities.ConnectFTS5(dataSourceName)
	if err != nil {
		return nil, errors.Connectionf("failed to connect to database: %w", err)
	}

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// VerifyFTS5Support checks if FTS5 is available and working
func (d *Database) VerifyFTS5Support(ctx context.Context) error {
	// Check compile options
	var fts5Available bool
	err := d.db.QueryRowContext(ctx, `
		SELECT COUNT(*) > 0 
		FROM pragma_compile_options 
		WHERE compile_options = 'ENABLE_FTS5'
	`).Scan(&fts5Available)

	if err != nil {
		return errors.FTS5f("failed to check FTS5 support: %w", err)
	}

	if !fts5Available {
		return errors.FTS5f("SQLite not compiled with FTS5 support")
	}

	return nil
}

// CreateTestTable creates a test FTS5 table for validation
func (d *Database) CreateTestTable(ctx context.Context, tableName string) error {
	err := utilities.CreateBasicDocumentTable(d.db, tableName)
	if err != nil {
		return errors.FTS5f("failed to create test table: %w", err)
	}
	return nil
}

// InsertTestData inserts sample documents for validation
func (d *Database) InsertTestData(ctx context.Context, tableName string) error {
	err := utilities.InsertSampleDocuments(d.db, tableName)
	if err != nil {
		return errors.Databasef("failed to insert test data: %w", err)
	}
	return nil
}

// QueryWithBM25 performs a BM25 query for validation testing
func (d *Database) QueryWithBM25(ctx context.Context, tableName, searchTerm string) (*sql.Rows, error) {
	rows, err := utilities.QueryFTS5WithBM25(d.db, tableName, searchTerm)
	if err != nil {
		return nil, errors.FTS5f("failed to execute BM25 query: %w", err)
	}
	return rows, nil
}

// CountDocuments returns the number of documents in a table
func (d *Database) CountDocuments(ctx context.Context, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, errors.Databasef("failed to count documents: %w", err)
	}
	return count, nil
}

// GetSQLiteVersion returns the SQLite version string
func (d *Database) GetSQLiteVersion(ctx context.Context) (string, error) {
	var version string
	err := d.db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)
	if err != nil {
		return "", errors.Databasef("failed to get SQLite version: %w", err)
	}
	return version, nil
}

// DB returns the underlying sql.DB for direct queries when needed
func (d *Database) DB() *sql.DB {
	return d.db
}