package models

import (
	"time"
)

// Document represents a document in the corpus
type Document struct {
	ID       int64     `json:"id" db:"id"`
	Title    string    `json:"title" db:"title"`
	Content  string    `json:"content" db:"content"`
	Category string    `json:"category" db:"category"`
	Length   int       `json:"length" db:"length"`         // Document length in tokens
	Created  time.Time `json:"created" db:"created"`
}

// DocumentInfo provides metadata about a document for analysis
type DocumentInfo struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Category      string    `json:"category"`
	TokenCount    int       `json:"token_count"`
	UniqueTerms   int       `json:"unique_terms"`
	AvgTermLength float64   `json:"avg_term_length"`
	Created       time.Time `json:"created"`
}