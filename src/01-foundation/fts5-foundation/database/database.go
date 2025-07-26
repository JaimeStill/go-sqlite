package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/errors"
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
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, errors.Databasef("failed to open database: %w", err)
	}

	// Configure SQLite for FTS5 operations
	if err := configureSQLite(db); err != nil {
		db.Close()
		return nil, err
	}

	return &Database{db: db}, nil
}

// configureSQLite applies optimal settings for FTS5 operations
func configureSQLite(db *sql.DB) error {
	settings := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=10000",
		"PRAGMA foreign_keys=ON",
		"PRAGMA temp_store=MEMORY",
	}

	for _, setting := range settings {
		if _, err := db.Exec(setting); err != nil {
			return errors.Databasef("failed to configure SQLite: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// VerifyFTS5Support checks if FTS5 is available
func (d *Database) VerifyFTS5Support(ctx context.Context) error {
	var available bool
	err := d.db.QueryRowContext(ctx, `
		SELECT COUNT(*) > 0 
		FROM pragma_compile_options 
		WHERE compile_options = 'ENABLE_FTS5'
	`).Scan(&available)

	if err != nil {
		return errors.FTS5f("failed to check FTS5 support: %w", err)
	}

	if !available {
		return errors.FTS5f("SQLite not compiled with FTS5 support")
	}

	return nil
}

// InitSchema creates the FTS5 tables and indexes
func (d *Database) InitSchema(ctx context.Context) error {
	// Check FTS5 support first
	if err := d.VerifyFTS5Support(ctx); err != nil {
		return err
	}

	// Create the FTS5 virtual table
	createTableSQL := `
		CREATE VIRTUAL TABLE IF NOT EXISTS documents USING fts5(
			title,           -- Document title for headline searches  
			content,         -- Main document body content
			category,        -- Document classification for filtering
			tokenize='unicode61 remove_diacritics 1'
		)`

	if _, err := d.db.ExecContext(ctx, createTableSQL); err != nil {
		return errors.FTS5f("failed to create FTS5 table: %w", err)
	}

	return nil
}

// Begin starts a new transaction
func (d *Database) Begin(ctx context.Context) (*sql.Tx, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Transactionf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// DB returns the underlying sql.DB for direct queries when needed
func (d *Database) DB() *sql.DB {
	return d.db
}