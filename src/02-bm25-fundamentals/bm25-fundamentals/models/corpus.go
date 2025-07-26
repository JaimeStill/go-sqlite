package models

import (
	"time"
)

// CorpusStats provides statistical information about the document corpus
type CorpusStats struct {
	TotalDocuments    int       `json:"total_documents"`
	TotalTokens       int64     `json:"total_tokens"`
	AverageDocLength  float64   `json:"average_doc_length"`
	MedianDocLength   float64   `json:"median_doc_length"`
	MinDocLength      int       `json:"min_doc_length"`
	MaxDocLength      int       `json:"max_doc_length"`
	UniqueTerms       int       `json:"unique_terms"`
	Categories        []string  `json:"categories"`
	CategoryCounts    map[string]int `json:"category_counts"`
	CreatedRange      TimeRange `json:"created_range"`
	LastUpdated       time.Time `json:"last_updated"`
}

// TimeRange represents a time span
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// CorpusOptions holds options for corpus generation
type CorpusOptions struct {
	Size           int      `json:"size"`
	Categories     []string `json:"categories"`
	MinTokens      int      `json:"min_tokens"`
	MaxTokens      int      `json:"max_tokens"`
	TitleMinTokens int      `json:"title_min_tokens"`
	TitleMaxTokens int      `json:"title_max_tokens"`
	Seed           int64    `json:"seed,omitempty"` // For reproducible generation
}

// DefaultCorpusOptions returns sensible default options
func DefaultCorpusOptions() CorpusOptions {
	return CorpusOptions{
		Size:           100,
		Categories:     []string{"technology", "science", "programming", "database", "algorithms"},
		MinTokens:      50,
		MaxTokens:      500,
		TitleMinTokens: 2,
		TitleMaxTokens: 8,
		Seed:           0, // 0 means use current time
	}
}