package models

// Document represents a document to be inserted into the FTS5 table
type Document struct {
	Title    string
	Content  string
	Category string
}

// SearchResult represents a search result from the FTS5 table
type SearchResult struct {
	RowID    int64
	Title    string
	Content  string
	Category string
	Score    float64
}

// DocumentInfo represents basic document information for listing
type DocumentInfo struct {
	RowID    int64
	Title    string
	Category string
	Preview  string // First 100 chars of content
}