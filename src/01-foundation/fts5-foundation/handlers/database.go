package handlers

import (
	"database/sql"
	"fmt"

	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/errors"
	"github.com/spf13/viper"
)

// CreateDocumentsTable creates the FTS5 virtual table for document storage
func CreateDocumentsTable() error {
	// Use the database path from global config
	dbPath := viper.GetString("database")
	
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// Create the FTS5 virtual table
	createTableSQL := `
		CREATE VIRTUAL TABLE IF NOT EXISTS documents USING fts5(
			title,           -- Document title for headline searches  
			content,         -- Main document body content
			category,        -- Document classification for filtering
			tokenize='unicode61 remove_diacritics 1'
		);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return errors.FTS5f("failed to create table: %w", err)
	}

	// Verify the table was created successfully
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='documents'").Scan(&tableName)
	if err != nil {
		return errors.Databasef("failed to verify table creation: %w", err)
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Created FTS5 virtual table: %s\n", tableName)
		fmt.Println("Schema: title, content, category with unicode61 tokenizer")
	}

	return nil
}

// VerifyFTS5Support checks if SQLite was compiled with FTS5 support
func VerifyFTS5Support() error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// Check if FTS5 is available
	rows, err := db.Query("PRAGMA compile_options")
	if err != nil {
		return errors.Databasef("failed to query compile options: %w", err)
	}
	defer rows.Close()

	fts5Found := false
	for rows.Next() {
		var option string
		if err := rows.Scan(&option); err != nil {
			continue
		}
		if option == "ENABLE_FTS5" {
			fts5Found = true
			break
		}
	}

	if !fts5Found {
		return errors.FTS5f("SQLite not compiled with FTS5 support. Please rebuild with -tags fts5")
	}

	// Test creating an FTS5 table
	_, err = db.Exec("CREATE VIRTUAL TABLE test_fts USING fts5(content)")
	if err != nil {
		return errors.FTS5f("failed to create test table: %w", err)
	}

	return nil
}