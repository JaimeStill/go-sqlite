package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/config"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/database"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/errors"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/models"
	"github.com/spf13/cobra"
)

// Search is the global search handler instance
var Search SearchHandler

// SearchHandler manages search operations and BM25 analysis (stateless - accesses global instances)
type SearchHandler struct{}

// GetSearchStats generates comprehensive statistics about search results
func (h *SearchHandler) GetSearchStats(ctx context.Context, results []*models.SearchResult, query string, executionTime time.Duration) (*models.SearchStats, error) {
	if len(results) == 0 {
		return &models.SearchStats{
			Query:         query,
			TotalResults:  0,
			ExecutionTime: executionTime,
		}, nil
	}

	stats := &models.SearchStats{
		Query:             query,
		TotalResults:      len(results),
		ExecutionTime:     executionTime,
		CategoryBreakdown: make(map[string]int),
	}

	// Collect scores and categories
	scores := make([]float64, len(results))
	for i, result := range results {
		scores[i] = result.Score
		stats.CategoryBreakdown[result.Category]++
	}

	// Calculate score distribution
	stats.ScoreDistrib = h.calculateScoreDistribution(scores)

	// Populate ScoreRange from distribution results
	stats.ScoreRange.Best = scores[0]              // Highest score (best/least negative)
	stats.ScoreRange.Worst = scores[len(scores)-1] // Lowest score (worst/most negative)
	stats.ScoreRange.Mean = stats.ScoreDistrib.Mean
	stats.ScoreRange.Median = stats.ScoreDistrib.Median
	stats.ScoreRange.StdDev = stats.ScoreDistrib.StdDev

	return stats, nil
}

// HandleCompare handles the search comparison command
func (h *SearchHandler) HandleCompare(cmd *cobra.Command, args []string) error {
	query, _ := cmd.Flags().GetString("query")
	compareWeights, _ := cmd.Flags().GetString("compare-weights")
	maxResults, _ := cmd.Flags().GetInt("max-results")

	ctx := context.Background()

	// Create baseline search (default weights)
	baselineOptions := models.DefaultSearchOptions()
	baselineOptions.Query = query
	baselineOptions.MaxResults = maxResults
	baselineOptions.IncludeSnippet = false

	baselineResults, err := h.Search(ctx, baselineOptions)
	if err != nil {
		return err
	}

	// Create comparison with custom weights
	comparisonResults := baselineResults // Default to same results
	var weights map[string]float64
	
	if compareWeights != "" {
		weights, err = h.parseWeights(compareWeights)
		if err != nil {
			return errors.Validationf("invalid weight format: %w", err)
		}

		comparisonOptions := models.DefaultSearchOptions()
		comparisonOptions.Query = query
		comparisonOptions.MaxResults = maxResults
		comparisonOptions.ColumnWeights = weights
		comparisonOptions.IncludeSnippet = false

		comparisonResults, err = h.Search(ctx, comparisonOptions)
		if err != nil {
			return err
		}
	}

	// Generate comparison analysis
	comparison := h.generateComparison(query, baselineResults, comparisonResults, weights)
	
	// Display comparison
	return h.displayComparison(comparison)
}

// HandleExplain handles the search explanation command
func (h *SearchHandler) HandleExplain(cmd *cobra.Command, args []string) error {
	// Extract flags
	query, _ := cmd.Flags().GetString("query")
	category, _ := cmd.Flags().GetString("category")
	titleWeight, _ := cmd.Flags().GetFloat64("title-weight")
	contentWeight, _ := cmd.Flags().GetFloat64("content-weight")
	categoryWeight, _ := cmd.Flags().GetFloat64("category-weight")
	maxResults, _ := cmd.Flags().GetInt("max-results")

	// Build search options
	options := models.DefaultSearchOptions()
	options.Query = query
	options.MaxResults = maxResults
	options.IncludeSnippet = false // Don't need snippets for explanations
	options.ExplainScores = true

	if category != "" {
		options.CategoryFilter = category
	}

	// Set column weights if specified
	if titleWeight > 0 || contentWeight > 0 || categoryWeight > 0 {
		options.ColumnWeights = make(map[string]float64)
		if titleWeight > 0 {
			options.ColumnWeights["title"] = titleWeight
		}
		if contentWeight > 0 {
			options.ColumnWeights["content"] = contentWeight
		}
		if categoryWeight > 0 {
			options.ColumnWeights["category"] = categoryWeight
		}
	}

	ctx := context.Background()

	// Perform search
	results, err := h.Search(ctx, options)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("No results found for query: \"%s\"\n", query)
		return nil
	}

	// Generate detailed explanations
	explanations, err := h.GenerateScoreExplanations(ctx, results, options)
	if err != nil {
		return err
	}

	// Display detailed explanations
	return h.displayScoreExplanations(explanations, options)
}

// HandleQuery handles the search query command
func (h *SearchHandler) HandleQuery(cmd *cobra.Command, args []string) error {
	// Extract flags
	query, _ := cmd.Flags().GetString("query")
	maxResults, _ := cmd.Flags().GetInt("max-results")
	category, _ := cmd.Flags().GetString("category")
	titleWeight, _ := cmd.Flags().GetFloat64("title-weight")
	contentWeight, _ := cmd.Flags().GetFloat64("content-weight")
	categoryWeight, _ := cmd.Flags().GetFloat64("category-weight")
	includeSnippets, _ := cmd.Flags().GetBool("snippets")
	snippetLength, _ := cmd.Flags().GetInt("snippet-length")

	// Build search options
	options := models.DefaultSearchOptions()
	options.Query = query

	if maxResults > 0 {
		options.MaxResults = maxResults
	} else {
		options.MaxResults = config.App.Search.MaxResults
	}

	if category != "" {
		options.CategoryFilter = category
	}

	// Set column weights if specified
	if titleWeight > 0 || contentWeight > 0 || categoryWeight > 0 {
		options.ColumnWeights = make(map[string]float64)
		if titleWeight > 0 {
			options.ColumnWeights["title"] = titleWeight
		}
		if contentWeight > 0 {
			options.ColumnWeights["content"] = contentWeight
		}
		if categoryWeight > 0 {
			options.ColumnWeights["category"] = categoryWeight
		}
	}

	options.IncludeSnippet = includeSnippets
	if snippetLength > 0 {
		options.SnippetLength = snippetLength
	}

	ctx := context.Background()

	// Perform search
	startTime := time.Now()
	results, err := h.Search(ctx, options)
	if err != nil {
		return err
	}
	executionTime := time.Since(startTime)

	// Display results
	return h.displaySearchResults(results, options, executionTime)
}

// HandleStats handles the search statistics command
func (h *SearchHandler) HandleStats(cmd *cobra.Command, args []string) error {
	// Extract flags
	query, _ := cmd.Flags().GetString("query")
	category, _ := cmd.Flags().GetString("category")
	titleWeight, _ := cmd.Flags().GetFloat64("title-weight")
	contentWeight, _ := cmd.Flags().GetFloat64("content-weight")
	categoryWeight, _ := cmd.Flags().GetFloat64("category-weight")

	// Build search options
	options := models.DefaultSearchOptions()
	options.Query = query
	options.MaxResults = 1000      // Get more results for better statistics
	options.IncludeSnippet = false // Don't need snippets for stats

	if category != "" {
		options.CategoryFilter = category
	}

	// Set column weights if specified
	if titleWeight > 0 || contentWeight > 0 || categoryWeight > 0 {
		options.ColumnWeights = make(map[string]float64)
		if titleWeight > 0 {
			options.ColumnWeights["title"] = titleWeight
		}
		if contentWeight > 0 {
			options.ColumnWeights["content"] = contentWeight
		}
		if categoryWeight > 0 {
			options.ColumnWeights["category"] = categoryWeight
		}
	}

	ctx := context.Background()

	// Perform search
	startTime := time.Now()
	results, err := h.Search(ctx, options)
	if err != nil {
		return err
	}
	executionTime := time.Since(startTime)

	// Generate statistics
	stats, err := h.GetSearchStats(ctx, results, query, executionTime)
	if err != nil {
		return err
	}

	// Display statistics
	return h.displaySearchStats(stats)
}

// Search performs FTS5 search with BM25 scoring
func (h *SearchHandler) Search(ctx context.Context, options models.SearchOptions) ([]*models.SearchResult, error) {
	// Build the search query
	query, args := h.buildSearchQuery(options)

	// Execute search
	rows, err := database.Instance.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.FTS5f("search query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.SearchResult

	for rows.Next() {
		result := &models.SearchResult{}

		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Content,
			&result.Category,
			&result.Length,
			&result.Created,
			&result.Score,
		)
		if err != nil {
			return nil, errors.Databasef("failed to scan search result: %w", err)
		}

		// Add snippet if requested
		if options.IncludeSnippet {
			result.Snippet = h.generateSnippet(result.Content, options.Query, options.SnippetLength)
		}

		// Classify relevance based on score
		result.Relevance = h.classifyRelevance(result.Score)

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Databasef("error iterating search results: %w", err)
	}

	return results, nil
}

// buildSearchQuery constructs the FTS5 search query with optional column weighting
func (h *SearchHandler) buildSearchQuery(options models.SearchOptions) (string, []interface{}) {
	var queryParts []string
	var args []interface{}

	// Base query with BM25 scoring
	baseQuery := `
		SELECT 
			d.id, d.title, d.content, d.category, d.length, d.created,
			%s as score
		FROM documents d
		JOIN documents_fts fts ON d.id = fts.rowid
		WHERE documents_fts MATCH ?`

	// Determine scoring method based on column weights
	var scoreExpr string
	if len(options.ColumnWeights) > 0 {
		// Custom column weighting
		weights := make([]string, 0, 3)
		if w, ok := options.ColumnWeights["title"]; ok {
			weights = append(weights, fmt.Sprintf("%.2f", w))
		} else {
			weights = append(weights, "1.0")
		}
		if w, ok := options.ColumnWeights["content"]; ok {
			weights = append(weights, fmt.Sprintf("%.2f", w))
		} else {
			weights = append(weights, "1.0")
		}
		if w, ok := options.ColumnWeights["category"]; ok {
			weights = append(weights, fmt.Sprintf("%.2f", w))
		} else {
			weights = append(weights, "1.0")
		}
		scoreExpr = fmt.Sprintf("bm25(documents_fts, %s)", strings.Join(weights, ", "))
	} else {
		// Default BM25 scoring
		scoreExpr = "bm25(documents_fts)"
	}

	queryParts = append(queryParts, fmt.Sprintf(baseQuery, scoreExpr))
	args = append(args, options.Query)

	// Add category filter if specified
	if options.CategoryFilter != "" {
		queryParts = append(queryParts, "AND d.category = ?")
		args = append(args, options.CategoryFilter)
	}

	// Order by BM25 score (remember: lower = better in SQLite FTS5)
	queryParts = append(queryParts, "ORDER BY score")

	// Add limit
	if options.MaxResults > 0 {
		queryParts = append(queryParts, "LIMIT ?")
		args = append(args, options.MaxResults)
	}

	query := strings.Join(queryParts, " ")
	return query, args
}

// calculateScoreDistribution computes statistical measures of score distribution
func (h *SearchHandler) calculateScoreDistribution(scores []float64) models.ScoreDistribution {
	sort.Float64s(scores) // Sort for percentile calculation

	distrib := models.ScoreDistribution{
		Percentiles: make(map[int]float64),
	}

	n := len(scores)
	if n == 0 {
		return distrib
	}

	// Calculate mean
	sum := 0.0
	for _, score := range scores {
		sum += score
	}
	distrib.Mean = sum / float64(n)

	// Calculate median (50th percentile)
	if n%2 == 0 {
		distrib.Median = (scores[n/2-1] + scores[n/2]) / 2
	} else {
		distrib.Median = scores[n/2]
	}

	// Calculate standard deviation
	sumSquareDiff := 0.0
	for _, score := range scores {
		diff := score - distrib.Mean
		sumSquareDiff += diff * diff
	}
	distrib.StdDev = math.Sqrt(sumSquareDiff / float64(n))

	// Calculate percentiles
	percentiles := []int{25, 50, 75, 90, 95, 99}
	for _, p := range percentiles {
		index := float64(p) / 100.0 * float64(n-1)
		lower := int(math.Floor(index))
		upper := int(math.Ceil(index))

		if lower == upper {
			distrib.Percentiles[p] = scores[lower]
		} else {
			// Linear interpolation
			weight := index - float64(lower)
			distrib.Percentiles[p] = scores[lower]*(1-weight) + scores[upper]*weight
		}
	}

	// Create score buckets for histogram
	distrib.Buckets = h.createScoreBuckets(scores, 10)

	return distrib
}

// classifyRelevance assigns a relevance label based on BM25 score
func (h *SearchHandler) classifyRelevance(score float64) string {
	// Remember: SQLite FTS5 BM25 scores are negative (lower = better)
	switch {
	case score >= -1.0:
		return "excellent"
	case score >= -2.0:
		return "good"
	case score >= -4.0:
		return "fair"
	default:
		return "poor"
	}
}

// createScoreBuckets creates histogram buckets for score distribution
func (h *SearchHandler) createScoreBuckets(scores []float64, numBuckets int) []models.ScoreBucket {
	if len(scores) == 0 {
		return nil
	}

	min := scores[0]
	max := scores[len(scores)-1]

	// Handle case where all scores are the same
	if min == max {
		return []models.ScoreBucket{{
			Min:   min,
			Max:   max,
			Count: len(scores),
			Label: fmt.Sprintf("%.2f", min),
		}}
	}

	buckets := make([]models.ScoreBucket, numBuckets)
	bucketWidth := (max - min) / float64(numBuckets)

	// Initialize buckets
	for i := 0; i < numBuckets; i++ {
		buckets[i].Min = min + float64(i)*bucketWidth
		buckets[i].Max = min + float64(i+1)*bucketWidth
		buckets[i].Label = fmt.Sprintf("%.2f to %.2f", buckets[i].Min, buckets[i].Max)
	}

	// Count scores in each bucket
	for _, score := range scores {
		bucketIndex := int((score - min) / bucketWidth)
		if bucketIndex >= numBuckets {
			bucketIndex = numBuckets - 1
		}
		buckets[bucketIndex].Count++
	}

	return buckets
}

// displaySearchResults formats and displays search results
func (h *SearchHandler) displaySearchResults(results []*models.SearchResult, options models.SearchOptions, executionTime time.Duration) error {
	switch config.App.Format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(map[string]interface{}{
			"query":          options.Query,
			"total_results":  len(results),
			"execution_time": executionTime.String(),
			"results":        results,
		})

	case "csv":
		fmt.Println("id,title,category,score,relevance")
		for _, result := range results {
			fmt.Printf("%d,\"%s\",\"%s\",%.4f,%s\n",
				result.ID, result.Title, result.Category, result.Score, result.Relevance)
		}

	default: // text format
		fmt.Printf("Search Results for: \"%s\"\n", options.Query)
		fmt.Printf("Found %d documents in %v\n", len(results), executionTime)

		if len(options.ColumnWeights) > 0 {
			fmt.Printf("Column weights: %v\n", options.ColumnWeights)
		}

		fmt.Printf("\n")

		if len(results) == 0 {
			fmt.Println("No documents found.")
			return nil
		}

		for i, result := range results {
			fmt.Printf("%d. %s\n", i+1, result.Title)
			fmt.Printf("   Score: %.4f (%s relevance)\n", result.Score, result.Relevance)
			fmt.Printf("   Category: %s | Length: %d tokens\n", result.Category, result.Length)

			if result.Snippet != "" {
				fmt.Printf("   Snippet: %s\n", result.Snippet)
			}

			if config.App.Verbose {
				fmt.Printf("   ID: %d | Created: %s\n", result.ID, result.Created.Format("2006-01-02 15:04"))
			}

			fmt.Printf("\n")
		}
	}

	return nil
}

// displaySearchStats formats and displays search statistics
func (h *SearchHandler) displaySearchStats(stats *models.SearchStats) error {
	switch config.App.Format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(stats)

	case "csv":
		fmt.Println("metric,value")
		fmt.Printf("query,\"%s\"\n", stats.Query)
		fmt.Printf("total_results,%d\n", stats.TotalResults)
		fmt.Printf("execution_time,%v\n", stats.ExecutionTime)
		fmt.Printf("score_best,%.4f\n", stats.ScoreRange.Best)
		fmt.Printf("score_worst,%.4f\n", stats.ScoreRange.Worst)
		fmt.Printf("score_mean,%.4f\n", stats.ScoreRange.Mean)
		fmt.Printf("score_median,%.4f\n", stats.ScoreRange.Median)
		fmt.Printf("score_stddev,%.4f\n", stats.ScoreRange.StdDev)

	default: // text format
		fmt.Printf("Search Statistics for: \"%s\"\n", stats.Query)
		fmt.Printf("========================================\n\n")

		fmt.Printf("Results: %d documents in %v\n\n", stats.TotalResults, stats.ExecutionTime)

		if stats.TotalResults == 0 {
			return nil
		}

		fmt.Printf("Score Distribution:\n")
		fmt.Printf("  Range:     %.4f to %.4f\n", stats.ScoreRange.Best, stats.ScoreRange.Worst)
		fmt.Printf("  Mean:      %.4f\n", stats.ScoreRange.Mean)
		fmt.Printf("  Median:    %.4f\n", stats.ScoreRange.Median)
		fmt.Printf("  Std Dev:   %.4f\n\n", stats.ScoreRange.StdDev)

		fmt.Printf("Percentiles:\n")
		for _, p := range []int{25, 50, 75, 90, 95, 99} {
			if score, ok := stats.ScoreDistrib.Percentiles[p]; ok {
				fmt.Printf("  %2dth:     %.4f\n", p, score)
			}
		}
		fmt.Printf("\n")

		if len(stats.CategoryBreakdown) > 0 {
			fmt.Printf("Category Breakdown:\n")
			for category, count := range stats.CategoryBreakdown {
				percentage := float64(count) * 100.0 / float64(stats.TotalResults)
				fmt.Printf("  %-15s: %3d documents (%.1f%%)\n", category, count, percentage)
			}
			fmt.Printf("\n")
		}

		if len(stats.ScoreDistrib.Buckets) > 0 {
			fmt.Printf("Score Distribution Buckets:\n")
			for _, bucket := range stats.ScoreDistrib.Buckets {
				if bucket.Count > 0 {
					fmt.Printf("  %s: %d documents\n", bucket.Label, bucket.Count)
				}
			}
		}
	}

	return nil
}

// generateSnippet creates a contextual snippet around search terms
func (h *SearchHandler) generateSnippet(content, query string, maxLength int) string {
	if maxLength <= 0 {
		maxLength = 200
	}

	// Simple snippet generation - find first occurrence of query terms
	queryTerms := strings.Fields(strings.ToLower(query))
	contentLower := strings.ToLower(content)

	var earliestPos int = len(content)
	for _, term := range queryTerms {
		if pos := strings.Index(contentLower, term); pos != -1 && pos < earliestPos {
			earliestPos = pos
		}
	}

	// If no terms found, return beginning of content
	if earliestPos == len(content) {
		earliestPos = 0
	}

	// Calculate snippet boundaries
	start := earliestPos - maxLength/4
	if start < 0 {
		start = 0
	}

	end := start + maxLength
	if end > len(content) {
		end = len(content)
		start = end - maxLength
		if start < 0 {
			start = 0
		}
	}

	snippet := content[start:end]

	// Add ellipsis if truncated
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	return strings.TrimSpace(snippet)
}

// parseWeights parses weight specification string
func (h *SearchHandler) parseWeights(weightStr string) (map[string]float64, error) {
	if weightStr == "" {
		return nil, nil
	}

	weights := make(map[string]float64)
	pairs := strings.Split(weightStr, ",")

	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid weight pair: %s", pair)
		}

		field := strings.TrimSpace(parts[0])
		weightStr := strings.TrimSpace(parts[1])

		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid weight value for %s: %w", field, err)
		}

		weights[field] = weight
	}

	return weights, nil
}

// GenerateScoreExplanations creates detailed BM25 score explanations for search results
func (h *SearchHandler) GenerateScoreExplanations(ctx context.Context, results []*models.SearchResult, options models.SearchOptions) ([]*models.ScoreExplanation, error) {
	explanations := make([]*models.ScoreExplanation, 0, len(results))
	
	// Get corpus statistics for average document length
	avgDocLength, err := h.getAverageDocumentLength(ctx)
	if err != nil {
		return nil, err
	}

	// Parse query terms
	queryTerms := strings.Fields(strings.ToLower(options.Query))
	
	for _, result := range results {
		explanation := &models.ScoreExplanation{
			DocumentID: result.ID,
			TotalScore: result.Score,
			FieldScores: make(map[string]models.FieldScore),
			QueryTerms: make([]models.TermScore, 0, len(queryTerms)),
			DocumentStats: models.DocumentStats{
				Length:       result.Length,
				AvgLength:    avgDocLength,
				LengthNorm:   h.calculateLengthNormalization(result.Length, avgDocLength),
				FieldLengths: map[string]int{
					"title":    len(strings.Fields(result.Title)),
					"content":  len(strings.Fields(result.Content)),
					"category": len(strings.Fields(result.Category)),
				},
			},
		}

		// Calculate term scores and field contributions
		for _, term := range queryTerms {
			termScore := h.calculateTermScore(term, result, avgDocLength)
			explanation.QueryTerms = append(explanation.QueryTerms, termScore)
		}

		// Calculate field scores
		explanation.FieldScores["title"] = h.calculateFieldScore("title", result.Title, queryTerms, options.ColumnWeights)
		explanation.FieldScores["content"] = h.calculateFieldScore("content", result.Content, queryTerms, options.ColumnWeights)
		explanation.FieldScores["category"] = h.calculateFieldScore("category", result.Category, queryTerms, options.ColumnWeights)

		explanations = append(explanations, explanation)
	}

	return explanations, nil
}

// displayScoreExplanations formats and displays detailed score explanations
func (h *SearchHandler) displayScoreExplanations(explanations []*models.ScoreExplanation, options models.SearchOptions) error {
	fmt.Printf("Score Explanations for: \"%s\"\n", options.Query)
	fmt.Printf("=====================================\n\n")

	if len(options.ColumnWeights) > 0 {
		fmt.Printf("Custom column weights: %v\n", options.ColumnWeights)
	} else {
		fmt.Printf("Using default FTS5 column weights (all fields weighted equally)\n")
	}
	
	fmt.Printf("BM25 parameters: k1=1.2, b=0.75 (SQLite FTS5 defaults)\n\n")

	for i, explanation := range explanations {
		fmt.Printf("Document %d (ID: %d)\n", i+1, explanation.DocumentID)
		fmt.Printf("Total Score: %.4f\n", explanation.TotalScore)
		fmt.Printf("Document Length: %d tokens (avg: %.1f)\n", 
			explanation.DocumentStats.Length, explanation.DocumentStats.AvgLength)
		fmt.Printf("Length Normalization Factor: %.3f\n", explanation.DocumentStats.LengthNorm)
		fmt.Printf("\n")

		// Display field contributions
		fmt.Printf("Field Contributions:\n")
		for fieldName, fieldScore := range explanation.FieldScores {
			fmt.Printf("  %s: score=%.4f, weight=%.2f, length=%d tokens\n",
				fieldName, fieldScore.Score, fieldScore.Weight,
				explanation.DocumentStats.FieldLengths[fieldName])
		}
		fmt.Printf("\n")

		// Display term analysis
		fmt.Printf("Query Term Analysis:\n")
		for _, termScore := range explanation.QueryTerms {
			fmt.Printf("  \"%s\": tf=%.3f, idf=%.3f, score=%.4f\n",
				termScore.Term, termScore.TF, termScore.IDF, termScore.Score)
		}
		
		fmt.Print("\n" + strings.Repeat("-", 50) + "\n\n")
	}

	return nil
}

// Helper methods for BM25 calculations

func (h *SearchHandler) getAverageDocumentLength(ctx context.Context) (float64, error) {
	query := "SELECT AVG(length) FROM documents"
	var avgLength float64
	err := database.Instance.DB().QueryRowContext(ctx, query).Scan(&avgLength)
	if err != nil {
		return 0, errors.Databasef("failed to get average document length: %w", err)
	}
	return avgLength, nil
}

func (h *SearchHandler) calculateLengthNormalization(docLength int, avgLength float64) float64 {
	// BM25 length normalization: k1 * ((1 - b) + b * (|d| / avgdl))
	k1 := 1.2
	b := 0.75
	return k1 * ((1 - b) + b * (float64(docLength) / avgLength))
}

func (h *SearchHandler) calculateTermScore(term string, result *models.SearchResult, avgDocLength float64) models.TermScore {
	// Simplified term frequency calculation (count occurrences in title + content)
	text := strings.ToLower(result.Title + " " + result.Content + " " + result.Category)
	termCount := strings.Count(text, term)
	tf := float64(termCount)
	
	// Simplified IDF calculation (would need corpus-wide term frequency in real implementation)
	// For now, use a mock IDF based on term length and frequency
	idf := 3.0 - (tf * 0.1) // Higher frequency terms get lower IDF
	if idf < 0.1 {
		idf = 0.1
	}
	
	// Simplified BM25 score for this term
	lengthNorm := h.calculateLengthNormalization(result.Length, avgDocLength)
	score := (idf * tf) / lengthNorm
	
	return models.TermScore{
		Term:    term,
		TF:      tf,
		IDF:     idf,
		FieldTF: tf, // Simplified - would normally be per-field
		Score:   score,
	}
}

func (h *SearchHandler) calculateFieldScore(fieldName, fieldContent string, queryTerms []string, weights map[string]float64) models.FieldScore {
	weight := 1.0
	if weights != nil {
		if w, ok := weights[fieldName]; ok {
			weight = w
		}
	}
	
	// Calculate field-specific score
	fieldText := strings.ToLower(fieldContent)
	score := 0.0
	termScores := make([]models.TermScore, 0, len(queryTerms))
	
	for _, term := range queryTerms {
		termCount := strings.Count(fieldText, term)
		tf := float64(termCount)
		
		// Field contributes to score based on term frequency and weight
		fieldContrib := tf * weight * 0.1 // Simplified scoring
		score += fieldContrib
		
		termScores = append(termScores, models.TermScore{
			Term:    term,
			TF:      tf,
			FieldTF: tf,
			Score:   fieldContrib,
		})
	}
	
	return models.FieldScore{
		Score:  score,
		Weight: weight,
		Terms:  termScores,
	}
}

// generateComparison creates a comparison analysis between baseline and custom weighted results
func (h *SearchHandler) generateComparison(query string, baseline, comparison []*models.SearchResult, weights map[string]float64) *models.SearchComparison {
	comp := &models.SearchComparison{
		Query:      query,
		Strategies: make(map[string]models.SearchStrategy),
		CommonDocs: make([]models.SearchResult, 0),
		UniqueDocs: make(map[string][]models.SearchResult),
	}

	// Create baseline strategy
	comp.Strategies["baseline"] = models.SearchStrategy{
		Name:        "Default FTS5",
		Description: "Standard FTS5 BM25 scoring with equal field weights",
		Config: models.StrategyConfig{
			ColumnWeights: nil,
			MaxResults:    len(baseline),
		},
		Results:  make([]models.SearchResult, len(baseline)),
		Analysis: models.ScoreAnalysis{}, // Would populate with actual analysis
	}

	// Copy baseline results (dereference pointers)
	for i, result := range baseline {
		comp.Strategies["baseline"].Results[i] = *result
	}

	// Create comparison strategy if weights provided
	if weights != nil {
		comp.Strategies["weighted"] = models.SearchStrategy{
			Name:        "Custom Weighted",
			Description: fmt.Sprintf("BM25 scoring with custom field weights: %v", weights),
			Config: models.StrategyConfig{
				ColumnWeights: weights,
				MaxResults:    len(comparison),
			},
			Results:  make([]models.SearchResult, len(comparison)),
			Analysis: models.ScoreAnalysis{},
		}

		// Copy comparison results
		for i, result := range comparison {
			comp.Strategies["weighted"].Results[i] = *result
		}
	}

	// Find common and unique documents
	h.analyzeResultOverlap(comp)

	return comp
}

// analyzeResultOverlap identifies common and unique documents between strategies
func (h *SearchHandler) analyzeResultOverlap(comp *models.SearchComparison) {
	if len(comp.Strategies) < 2 {
		return
	}

	baselineResults := comp.Strategies["baseline"].Results
	weightedResults := comp.Strategies["weighted"].Results

	// Create maps for quick lookup
	baselineMap := make(map[int64]models.SearchResult)
	weightedMap := make(map[int64]models.SearchResult)

	for _, result := range baselineResults {
		baselineMap[result.ID] = result
	}

	for _, result := range weightedResults {
		weightedMap[result.ID] = result
	}

	// Find common documents
	for id, baselineResult := range baselineMap {
		if weightedResult, exists := weightedMap[id]; exists {
			// Document appears in both - compare scores
			if baselineResult.Score != weightedResult.Score {
				// Use the baseline version but note the score difference
				common := baselineResult
				comp.CommonDocs = append(comp.CommonDocs, common)
			} else {
				comp.CommonDocs = append(comp.CommonDocs, baselineResult)
			}
		}
	}

	// Find unique documents
	comp.UniqueDocs["baseline"] = make([]models.SearchResult, 0)
	comp.UniqueDocs["weighted"] = make([]models.SearchResult, 0)

	// Documents only in baseline
	for id, result := range baselineMap {
		if _, exists := weightedMap[id]; !exists {
			comp.UniqueDocs["baseline"] = append(comp.UniqueDocs["baseline"], result)
		}
	}

	// Documents only in weighted results
	for id, result := range weightedMap {
		if _, exists := baselineMap[id]; !exists {
			comp.UniqueDocs["weighted"] = append(comp.UniqueDocs["weighted"], result)
		}
	}
}

// displayComparison formats and displays search strategy comparison
func (h *SearchHandler) displayComparison(comp *models.SearchComparison) error {
	fmt.Printf("Search Strategy Comparison for: \"%s\"\n", comp.Query)
	fmt.Printf("===============================================\n\n")

	// Display strategies
	for name, strategy := range comp.Strategies {
		fmt.Printf("Strategy: %s (%s)\n", strategy.Name, name)
		fmt.Printf("Description: %s\n", strategy.Description)
		if len(strategy.Config.ColumnWeights) > 0 {
			fmt.Printf("Weights: %v\n", strategy.Config.ColumnWeights)
		}
		fmt.Printf("Results: %d documents\n\n", len(strategy.Results))
	}

	// Show ranking comparison if we have both strategies
	if baseline, hasBaseline := comp.Strategies["baseline"]; hasBaseline {
		if weighted, hasWeighted := comp.Strategies["weighted"]; hasWeighted {
			fmt.Printf("Ranking Comparison:\n")
			fmt.Printf("%-4s %-30s %-12s %-12s %-10s\n", "Rank", "Document", "Baseline", "Weighted", "Change")
			fmt.Printf("%s\n", strings.Repeat("-", 70))

			maxLen := len(baseline.Results)
			if len(weighted.Results) > maxLen {
				maxLen = len(weighted.Results)
			}

			for i := 0; i < maxLen; i++ {
				var baseDoc, weightedDoc *models.SearchResult
				var baseScore, weightedScore float64 = 0, 0
				var title string = "N/A"

				if i < len(baseline.Results) {
					baseDoc = &baseline.Results[i]
					baseScore = baseDoc.Score
					title = baseDoc.Title
					if len(title) > 25 {
						title = title[:25] + "..."
					}
				}

				if i < len(weighted.Results) {
					weightedDoc = &weighted.Results[i]
					weightedScore = weightedDoc.Score
					if baseDoc == nil {
						title = weightedDoc.Title
						if len(title) > 25 {
							title = title[:25] + "..."
						}
					}
				}

				var change string
				if baseDoc != nil && weightedDoc != nil {
					if baseDoc.ID == weightedDoc.ID {
						scoreDiff := weightedScore - baseScore
						if scoreDiff > 0.001 {
							change = fmt.Sprintf("+%.3f", scoreDiff)
						} else if scoreDiff < -0.001 {
							change = fmt.Sprintf("%.3f", scoreDiff)
						} else {
							change = "same"
						}
					} else {
						change = "reordered"
					}
				} else if baseDoc != nil {
					change = "dropped"
				} else {
					change = "new"
				}

				fmt.Printf("%-4d %-30s %-12.4f %-12.4f %-10s\n", 
					i+1, title, baseScore, weightedScore, change)
			}
			fmt.Printf("\n")
		}
	}

	// Summary statistics
	fmt.Printf("Summary:\n")
	fmt.Printf("Common documents: %d\n", len(comp.CommonDocs))
	
	if len(comp.UniqueDocs) > 0 {
		for strategy, uniqueDocs := range comp.UniqueDocs {
			if len(uniqueDocs) > 0 {
				fmt.Printf("Unique to %s: %d documents\n", strategy, len(uniqueDocs))
			}
		}
	}

	return nil
}
