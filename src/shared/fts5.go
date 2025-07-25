package shared

import (
	"database/sql"
	"fmt"
)

// CreateFTS5Table creates a basic FTS5 virtual table
func CreateFTS5Table(db *sql.DB, tableName string, columns []string) error {
	columnList := ""
	for i, col := range columns {
		if i > 0 {
			columnList += ", "
		}
		columnList += col
	}
	
	query := fmt.Sprintf("CREATE VIRTUAL TABLE %s USING fts5(%s)", tableName, columnList)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create FTS5 table %s: %w", tableName, err)
	}
	
	return nil
}

// CreateBasicDocumentTable creates a simple FTS5 table for document experiments
func CreateBasicDocumentTable(db *sql.DB, tableName string) error {
	return CreateFTS5Table(db, tableName, []string{"id", "title", "content"})
}

// QueryFTS5 performs a basic FTS5 MATCH query
func QueryFTS5(db *sql.DB, tableName, searchTerm string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT id, title, content FROM %s WHERE %s MATCH ? ORDER BY rank", tableName, tableName)
	rows, err := db.Query(query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to query FTS5 table: %w", err)
	}
	return rows, nil
}

// QueryFTS5WithBM25 performs an FTS5 query with BM25 scoring
func QueryFTS5WithBM25(db *sql.DB, tableName, searchTerm string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT id, title, content, bm25(%s) as score FROM %s WHERE %s MATCH ? ORDER BY bm25(%s)", 
		tableName, tableName, tableName, tableName)
	rows, err := db.Query(query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to query FTS5 table with BM25: %w", err)
	}
	return rows, nil
}