package shared

import (
	"database/sql"
	"fmt"
)

// Document represents a simple document for FTS5 experiments
type Document struct {
	ID      int
	Title   string
	Content string
}

// SampleDocuments returns a set of test documents for learning experiments
func SampleDocuments() []Document {
	return []Document{
		{1, "SQLite Introduction", "SQLite is a lightweight database engine that supports full-text search through FTS5."},
		{2, "BM25 Algorithm", "BM25 is a ranking function used by search engines to estimate relevance of documents to queries."},
		{3, "Full-Text Search", "Full-text search allows users to search for documents containing specific words or phrases."},
		{4, "Database Indexing", "Indexes improve query performance by creating efficient data structures for searching."},
		{5, "Information Retrieval", "Information retrieval systems help users find relevant documents from large collections."},
	}
}

// InsertSampleDocuments inserts test documents into an FTS5 table
func InsertSampleDocuments(db *sql.DB, tableName string) error {
	docs := SampleDocuments()
	
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (id, title, content) VALUES (?, ?, ?)", tableName))
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()
	
	for _, doc := range docs {
		if _, err := stmt.Exec(doc.ID, doc.Title, doc.Content); err != nil {
			return fmt.Errorf("failed to insert document %d: %w", doc.ID, err)
		}
	}
	
	return nil
}