package models

// ScoreAnalysis contains detailed BM25 score analysis results
type ScoreAnalysis struct {
	Query           string        `json:"query"`
	TotalResults    int           `json:"total_results"`
	ScoreRange      ScoreRange    `json:"score_range"`
	Distribution    []ScoreBucket `json:"distribution"`
	TopTerms        []TermFreq    `json:"top_terms"`
	CategoryBreakdown map[string]CategoryStats `json:"category_breakdown,omitempty"`
}

// ScoreRange provides statistical information about score distribution
type ScoreRange struct {
	Best   float64 `json:"best"`    // Least negative (best match)
	Worst  float64 `json:"worst"`   // Most negative (worst match)
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	StdDev float64 `json:"std_dev"`
}

// ScoreBucket represents a histogram bucket for score distribution
type ScoreBucket struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int     `json:"count"`
	Label string  `json:"label"`
}

// TermFreq represents term frequency information
type TermFreq struct {
	Term      string  `json:"term"`
	Frequency int     `json:"frequency"`
	Documents int     `json:"documents"` // Number of documents containing this term
	IDF       float64 `json:"idf"`       // Inverse document frequency
}

// CategoryStats provides category-specific scoring statistics
type CategoryStats struct {
	DocumentCount int        `json:"document_count"`
	AverageScore  float64    `json:"average_score"`
	ScoreRange    ScoreRange `json:"score_range"`
}

// SearchComparison compares results from different search strategies
type SearchComparison struct {
	Query       string                    `json:"query"`
	Strategies  map[string]SearchStrategy `json:"strategies"`
	CommonDocs  []SearchResult         `json:"common_docs"`  // Documents in all result sets
	UniqueDocs  map[string][]SearchResult `json:"unique_docs"` // Documents unique to each strategy
}

// SearchStrategy represents a specific search configuration and its results
type SearchStrategy struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Config      StrategyConfig   `json:"config"`
	Results     []SearchResult `json:"results"`
	Analysis    ScoreAnalysis    `json:"analysis"`
}

// StrategyConfig holds configuration for a search strategy
type StrategyConfig struct {
	ColumnWeights map[string]float64 `json:"column_weights,omitempty"`
	MaxResults    int                `json:"max_results"`
	FieldFilter   string             `json:"field_filter,omitempty"`
}

// ScoreExplanation provides detailed BM25 score breakdown
type ScoreExplanation struct {
	DocumentID    int64                    `json:"document_id"`
	TotalScore    float64                  `json:"total_score"`
	FieldScores   map[string]FieldScore    `json:"field_scores"`
	QueryTerms    []TermScore              `json:"query_terms"`
	DocumentStats DocumentStats            `json:"document_stats"`
}

// FieldScore breaks down scoring by field (title, content, category)
type FieldScore struct {
	Score  float64 `json:"score"`
	Weight float64 `json:"weight"`
	Terms  []TermScore `json:"terms"`
}

// TermScore provides per-term scoring details
type TermScore struct {
	Term       string  `json:"term"`
	TF         float64 `json:"tf"`          // Term frequency in document
	IDF        float64 `json:"idf"`         // Inverse document frequency
	FieldTF    float64 `json:"field_tf"`    // Term frequency in specific field
	Score      float64 `json:"score"`       // Final BM25 score for this term
}

// DocumentStats provides document-level statistics affecting BM25
type DocumentStats struct {
	Length        int     `json:"length"`         // Document length in tokens
	AvgLength     float64 `json:"avg_length"`     // Average document length in corpus
	LengthNorm    float64 `json:"length_norm"`    // BM25 length normalization factor
	FieldLengths  map[string]int `json:"field_lengths"` // Length of each field
}

// ScoreDistribution provides percentile analysis
type ScoreDistribution struct {
	Mean       float64            `json:"mean"`
	Median     float64            `json:"median"`
	StdDev     float64            `json:"std_dev"`
	Percentiles map[int]float64   `json:"percentiles"` // 25th, 50th, 75th, 90th, 95th, 99th
	Buckets    []ScoreBucket     `json:"buckets"`
}