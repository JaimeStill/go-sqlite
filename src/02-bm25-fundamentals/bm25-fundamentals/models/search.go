package models

import (
	"time"
)

// SearchResult represents a document with its BM25 score
type SearchResult struct {
	Document
	Score     float64 `json:"score"`
	Snippet   string  `json:"snippet,omitempty"`
	Relevance string  `json:"relevance,omitempty"` // "excellent", "good", "fair", "poor"
}

// SearchOptions holds parameters for search queries
type SearchOptions struct {
	Query          string            `json:"query"`
	MaxResults     int               `json:"max_results"`
	ColumnWeights  map[string]float64 `json:"column_weights,omitempty"`
	CategoryFilter string            `json:"category_filter,omitempty"`
	IncludeSnippet bool              `json:"include_snippet"`
	SnippetLength  int               `json:"snippet_length"`
	ExplainScores  bool              `json:"explain_scores"`
}

// DefaultSearchOptions returns sensible defaults for search
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		MaxResults:     20,
		ColumnWeights:  nil, // Use FTS5 defaults
		IncludeSnippet: true,
		SnippetLength:  200,
		ExplainScores:  false,
	}
}

// SearchStats provides statistics about search results
type SearchStats struct {
	Query           string        `json:"query"`
	TotalResults    int           `json:"total_results"`
	ExecutionTime   time.Duration `json:"execution_time"`
	ScoreRange      ScoreRange    `json:"score_range"`
	ScoreDistrib    ScoreDistribution `json:"score_distribution"`
	CategoryBreakdown map[string]int `json:"category_breakdown"`
}


