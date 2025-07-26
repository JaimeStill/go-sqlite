package models

import "time"

// ValidationResult represents the outcome of a validation check
type ValidationResult struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Passed      bool          `json:"passed"`
	Error       error         `json:"error,omitempty"`
	Duration    time.Duration `json:"duration"`
	Details     string        `json:"details,omitempty"`
}

// ValidationSuite represents a collection of validation checks
type ValidationSuite struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	StartTime   time.Time          `json:"start_time"`
	EndTime     time.Time          `json:"end_time"`
	Duration    time.Duration      `json:"duration"`
	Results     []ValidationResult `json:"results"`
	Passed      int                `json:"passed"`
	Failed      int                `json:"failed"`
	Total       int                `json:"total"`
}

// SystemInfo holds information about the system environment
type SystemInfo struct {
	SQLiteVersion string            `json:"sqlite_version"`
	FTS5Available bool              `json:"fts5_available"`
	GoVersion     string            `json:"go_version"`
	Platform      string            `json:"platform"`
	Environment   map[string]string `json:"environment,omitempty"`
}

// BM25TestResult represents results from BM25 scoring validation
type BM25TestResult struct {
	Query       string  `json:"query"`
	ResultCount int     `json:"result_count"`
	FirstScore  float64 `json:"first_score"`
	LastScore   float64 `json:"last_score"`
	HasNegative bool    `json:"has_negative_scores"`
	MinScore    float64 `json:"min_score"`
	MaxScore    float64 `json:"max_score"`
}

// AddResult adds a validation result to the suite
func (vs *ValidationSuite) AddResult(result ValidationResult) {
	vs.Results = append(vs.Results, result)
	vs.Total++
	
	if result.Passed {
		vs.Passed++
	} else {
		vs.Failed++
	}
}

// IsSuccessful returns true if all validations passed
func (vs *ValidationSuite) IsSuccessful() bool {
	return vs.Failed == 0 && vs.Total > 0
}

// GetSuccessRate returns the percentage of passed validations
func (vs *ValidationSuite) GetSuccessRate() float64 {
	if vs.Total == 0 {
		return 0.0
	}
	return float64(vs.Passed) / float64(vs.Total) * 100.0
}