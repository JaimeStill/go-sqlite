package handlers

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/guptarohit/asciigraph"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/config"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/models"
	"github.com/spf13/cobra"
)

// Visualize is the global visualize handler instance
var Visualize VisualizeHandler

// VisualizeHandler manages visualization operations (stateless - accesses global instances)
type VisualizeHandler struct{}

// HandleDistribution handles the score distribution visualization command
func (h *VisualizeHandler) HandleDistribution(cmd *cobra.Command, args []string) error {
	// Extract flags
	query, _ := cmd.Flags().GetString("query")
	category, _ := cmd.Flags().GetString("category")
	titleWeight, _ := cmd.Flags().GetFloat64("title-weight")
	contentWeight, _ := cmd.Flags().GetFloat64("content-weight")
	categoryWeight, _ := cmd.Flags().GetFloat64("category-weight")
	buckets, _ := cmd.Flags().GetInt("buckets")
	maxResults, _ := cmd.Flags().GetInt("max-results")

	// Build search options
	options := models.DefaultSearchOptions()
	options.Query = query
	options.MaxResults = maxResults
	options.IncludeSnippet = false

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

	// Perform search to get results
	results, err := Search.Search(ctx, options)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("No results found for query: \"%s\"\n", query)
		return nil
	}

	// Generate score distribution
	distribution := h.generateScoreDistribution(results, buckets)
	
	// Display histogram
	return h.displayDistributionHistogram(distribution, options)
}

// HandleCategories handles the category comparison visualization command
func (h *VisualizeHandler) HandleCategories(cmd *cobra.Command, args []string) error {
	// Extract flags
	query, _ := cmd.Flags().GetString("query")
	filter, _ := cmd.Flags().GetString("filter")
	titleWeight, _ := cmd.Flags().GetFloat64("title-weight")
	contentWeight, _ := cmd.Flags().GetFloat64("content-weight")
	categoryWeight, _ := cmd.Flags().GetFloat64("category-weight")
	maxResults, _ := cmd.Flags().GetInt("max-results")

	// Build search options
	options := models.DefaultSearchOptions()
	options.Query = query
	options.MaxResults = maxResults
	options.IncludeSnippet = false

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

	// Perform search to get results
	results, err := Search.Search(ctx, options)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("No results found for query: \"%s\"\n", query)
		return nil
	}

	// Parse filter categories
	var filterCategories []string
	if filter != "" {
		filterCategories = strings.Split(filter, ",")
		for i, cat := range filterCategories {
			filterCategories[i] = strings.TrimSpace(cat)
		}
	}

	// Generate category comparison
	comparison := h.generateCategoryComparison(results, filterCategories)
	
	// Display category charts
	return h.displayCategoryComparison(comparison, options)
}

// HandleRange handles the score range visualization command
func (h *VisualizeHandler) HandleRange(cmd *cobra.Command, args []string) error {
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
	options.IncludeSnippet = false

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

	// Perform search to get results
	results, err := Search.Search(ctx, options)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("No results found for query: \"%s\"\n", query)
		return nil
	}

	// Generate range analysis
	rangeAnalysis := h.generateRangeAnalysis(results)
	
	// Display range visualization
	return h.displayRangeVisualization(rangeAnalysis, options)
}

// Helper methods for visualization generation

// generateScoreDistribution creates histogram buckets for score distribution
func (h *VisualizeHandler) generateScoreDistribution(results []*models.SearchResult, numBuckets int) []models.ScoreBucket {
	if len(results) == 0 {
		return nil
	}

	scores := make([]float64, len(results))
	for i, result := range results {
		scores[i] = result.Score
	}

	sort.Float64s(scores)
	
	min := scores[0]
	max := scores[len(scores)-1]

	// Handle case where all scores are the same
	if min == max {
		return []models.ScoreBucket{{
			Min:   min,
			Max:   max,
			Count: len(scores),
			Label: fmt.Sprintf("%.3f", min),
		}}
	}

	buckets := make([]models.ScoreBucket, numBuckets)
	bucketWidth := (max - min) / float64(numBuckets)

	// Initialize buckets
	for i := 0; i < numBuckets; i++ {
		buckets[i].Min = min + float64(i)*bucketWidth
		buckets[i].Max = min + float64(i+1)*bucketWidth
		buckets[i].Label = fmt.Sprintf("%.3f to %.3f", buckets[i].Min, buckets[i].Max)
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

// CategoryData holds score data for a specific category
type CategoryData struct {
	Name     string
	Scores   []float64
	Count    int
	AvgScore float64
	MinScore float64
	MaxScore float64
}

// generateCategoryComparison creates category-wise score analysis
func (h *VisualizeHandler) generateCategoryComparison(results []*models.SearchResult, filter []string) map[string]*CategoryData {
	categoryMap := make(map[string]*CategoryData)

	// Group results by category
	for _, result := range results {
		category := result.Category
		
		// Apply filter if specified
		if len(filter) > 0 {
			found := false
			for _, filteredCat := range filter {
				if category == filteredCat {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if _, exists := categoryMap[category]; !exists {
			categoryMap[category] = &CategoryData{
				Name:     category,
				Scores:   make([]float64, 0),
				MinScore: math.MaxFloat64,
				MaxScore: -math.MaxFloat64,
			}
		}

		data := categoryMap[category]
		data.Scores = append(data.Scores, result.Score)
		data.Count++
		
		if result.Score < data.MinScore {
			data.MinScore = result.Score
		}
		if result.Score > data.MaxScore {
			data.MaxScore = result.Score
		}
	}

	// Calculate averages
	for _, data := range categoryMap {
		if data.Count > 0 {
			sum := 0.0
			for _, score := range data.Scores {
				sum += score
			}
			data.AvgScore = sum / float64(data.Count)
		}
	}

	return categoryMap
}

// RangeAnalysis holds comprehensive range and percentile data
type RangeAnalysis struct {
	Scores      []float64
	Min         float64
	Max         float64
	Mean        float64
	Median      float64
	StdDev      float64
	Percentiles map[int]float64
	Quartiles   []float64 // Q1, Q2, Q3
}

// generateRangeAnalysis creates comprehensive score range analysis
func (h *VisualizeHandler) generateRangeAnalysis(results []*models.SearchResult) *RangeAnalysis {
	scores := make([]float64, len(results))
	for i, result := range results {
		scores[i] = result.Score
	}

	sort.Float64s(scores)

	analysis := &RangeAnalysis{
		Scores:      scores,
		Min:         scores[0],
		Max:         scores[len(scores)-1],
		Percentiles: make(map[int]float64),
		Quartiles:   make([]float64, 3),
	}

	// Calculate mean
	sum := 0.0
	for _, score := range scores {
		sum += score
	}
	analysis.Mean = sum / float64(len(scores))

	// Calculate median
	n := len(scores)
	if n%2 == 0 {
		analysis.Median = (scores[n/2-1] + scores[n/2]) / 2
	} else {
		analysis.Median = scores[n/2]
	}

	// Calculate standard deviation
	sumSquareDiff := 0.0
	for _, score := range scores {
		diff := score - analysis.Mean
		sumSquareDiff += diff * diff
	}
	analysis.StdDev = math.Sqrt(sumSquareDiff / float64(n))

	// Calculate percentiles
	percentiles := []int{25, 50, 75, 90, 95, 99}
	for _, p := range percentiles {
		index := float64(p) / 100.0 * float64(n-1)
		lower := int(math.Floor(index))
		upper := int(math.Ceil(index))

		if lower == upper {
			analysis.Percentiles[p] = scores[lower]
		} else {
			weight := index - float64(lower)
			analysis.Percentiles[p] = scores[lower]*(1-weight) + scores[upper]*weight
		}
	}

	// Set quartiles
	analysis.Quartiles[0] = analysis.Percentiles[25] // Q1
	analysis.Quartiles[1] = analysis.Percentiles[50] // Q2 (median)
	analysis.Quartiles[2] = analysis.Percentiles[75] // Q3

	return analysis
}

// Display methods for different visualizations

// displayDistributionHistogram shows an ASCII histogram of score distribution using asciigraph
func (h *VisualizeHandler) displayDistributionHistogram(buckets []models.ScoreBucket, options models.SearchOptions) error {
	fmt.Printf("Score Distribution Histogram for: \"%s\"\n", options.Query)
	fmt.Printf("==========================================\n\n")

	if len(options.ColumnWeights) > 0 {
		fmt.Printf("Column weights: %v\n\n", options.ColumnWeights)
	}

	if len(buckets) == 0 {
		fmt.Println("No data to display.")
		return nil
	}

	// Prepare data for asciigraph
	data := make([]float64, len(buckets))
	labels := make([]string, len(buckets))
	totalCount := 0
	
	for i, bucket := range buckets {
		data[i] = float64(bucket.Count)
		labels[i] = fmt.Sprintf("%.3f", (bucket.Min+bucket.Max)/2) // Use bucket midpoint
		totalCount += bucket.Count
	}

	// Create histogram using asciigraph
	graph := asciigraph.Plot(data, asciigraph.Height(15), asciigraph.Width(70), 
		asciigraph.Caption("Document Count per Score Range"))
	
	fmt.Println(graph)
	fmt.Printf("\n")

	// Display detailed bucket information
	fmt.Printf("Detailed Distribution:\n")
	fmt.Printf("%-20s │ %8s │ %10s\n", "Score Range", "Count", "Percentage")
	fmt.Printf("%s\n", strings.Repeat("-", 42))
	
	for _, bucket := range buckets {
		if bucket.Count > 0 {
			percentage := float64(bucket.Count) * 100.0 / float64(totalCount)
			fmt.Printf("%-20s │ %8d │ %9.1f%%\n", bucket.Label, bucket.Count, percentage)
		}
	}

	fmt.Printf("\nTotal documents: %d\n", totalCount)
	fmt.Printf("Note: Lower scores indicate better relevance (SQLite FTS5 uses negative BM25)\n")

	return nil
}

// displayCategoryComparison shows category-wise score comparison using asciigraph
func (h *VisualizeHandler) displayCategoryComparison(categories map[string]*CategoryData, options models.SearchOptions) error {
	fmt.Printf("Category Score Comparison for: \"%s\"\n", options.Query)
	fmt.Printf("==========================================\n\n")

	if len(options.ColumnWeights) > 0 {
		fmt.Printf("Column weights: %v\n\n", options.ColumnWeights)
	}

	if len(categories) == 0 {
		fmt.Println("No categories to display.")
		return nil
	}

	// Sort categories by average score (best to worst)
	type categoryPair struct {
		name string
		data *CategoryData
	}
	
	sortedCategories := make([]categoryPair, 0, len(categories))
	for name, data := range categories {
		sortedCategories = append(sortedCategories, categoryPair{name, data})
	}
	
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].data.AvgScore > sortedCategories[j].data.AvgScore // Higher scores first (less negative)
	})

	// Prepare data for average score comparison chart
	avgScores := make([]float64, len(sortedCategories))
	categoryNames := make([]string, len(sortedCategories))
	
	for i, pair := range sortedCategories {
		avgScores[i] = pair.data.AvgScore
		categoryNames[i] = pair.data.Name
	}

	// Create average score comparison chart
	fmt.Printf("Average Score by Category:\n")
	graph := asciigraph.Plot(avgScores, asciigraph.Height(12), asciigraph.Width(60),
		asciigraph.Caption("Average BM25 Scores (higher is better)"))
	
	fmt.Println(graph)
	fmt.Printf("\n")

	// Show category labels and data
	fmt.Printf("Category Details:\n")
	fmt.Printf("%-15s │ %8s │ %5s │ %12s │ %8s\n", "Category", "Avg Score", "Count", "Score Range", "Std Dev")
	fmt.Printf("%s\n", strings.Repeat("-", 70))

	for _, pair := range sortedCategories {
		data := pair.data
		
		// Calculate standard deviation for this category
		var stdDev float64
		if len(data.Scores) > 1 {
			sum := 0.0
			for _, score := range data.Scores {
				sum += score
			}
			mean := sum / float64(len(data.Scores))
			
			sumSquareDiff := 0.0
			for _, score := range data.Scores {
				diff := score - mean
				sumSquareDiff += diff * diff
			}
			stdDev = math.Sqrt(sumSquareDiff / float64(len(data.Scores)))
		}
		
		fmt.Printf("%-15s │ %8.3f │ %5d │ %5.3f-%5.3f │ %8.3f\n",
			data.Name, data.AvgScore, data.Count, data.MinScore, data.MaxScore, stdDev)
	}

	// Create a second chart showing document counts by category
	if len(sortedCategories) > 1 {
		fmt.Printf("\nDocument Count by Category:\n")
		countData := make([]float64, len(sortedCategories))
		for i, pair := range sortedCategories {
			countData[i] = float64(pair.data.Count)
		}
		
		countGraph := asciigraph.Plot(countData, asciigraph.Height(8), asciigraph.Width(60),
			asciigraph.Caption("Number of Documents per Category"))
		
		fmt.Println(countGraph)
		fmt.Printf("\n")
	}

	fmt.Printf("Note: Higher positioned categories have better average relevance\n")
	
	return nil
}

// displayRangeVisualization shows score range and percentile analysis
func (h *VisualizeHandler) displayRangeVisualization(analysis *RangeAnalysis, options models.SearchOptions) error {
	fmt.Printf("Score Range Analysis for: \"%s\"\n", options.Query)
	fmt.Printf("=====================================\n\n")

	if len(options.ColumnWeights) > 0 {
		fmt.Printf("Column weights: %v\n\n", options.ColumnWeights)
	}

	// Display statistics
	fmt.Printf("Statistical Summary:\n")
	fmt.Printf("  Range:     %.4f to %.4f (span: %.4f)\n", analysis.Min, analysis.Max, analysis.Max-analysis.Min)
	fmt.Printf("  Mean:      %.4f\n", analysis.Mean)
	fmt.Printf("  Median:    %.4f\n", analysis.Median)
	fmt.Printf("  Std Dev:   %.4f\n\n", analysis.StdDev)

	// Display percentiles
	fmt.Printf("Percentiles:\n")
	percentiles := []int{25, 50, 75, 90, 95, 99}
	for _, p := range percentiles {
		if score, ok := analysis.Percentiles[p]; ok {
			fmt.Printf("  %2dth:     %.4f\n", p, score)
		}
	}
	fmt.Printf("\n")

	// Visual range representation using asciigraph
	fmt.Printf("Score Range Visualization:\n")
	
	if scoreRange := analysis.Max - analysis.Min; scoreRange > 0 {
		// Create data points for key statistics to visualize as a line chart
		// We'll plot percentiles and mean as a distribution curve
		dataPoints := []float64{}
		labels := []string{}
		
		// Add key percentile points for visualization
		keyPercentiles := []int{25, 50, 75, 90, 95}
		for _, p := range keyPercentiles {
			if score, ok := analysis.Percentiles[p]; ok {
				dataPoints = append(dataPoints, score)
				labels = append(labels, fmt.Sprintf("P%d", p))
			}
		}
		
		// Add mean
		dataPoints = append(dataPoints, analysis.Mean)
		labels = append(labels, "Mean")
		
		// Create a line graph showing the score distribution
		if len(dataPoints) > 1 {
			graph := asciigraph.Plot(dataPoints, 
				asciigraph.Height(8), 
				asciigraph.Width(60),
				asciigraph.Caption("Score Distribution (Percentiles and Mean)"))
			fmt.Println(graph)
			fmt.Printf("\n")
		}
		
		// Enhanced textual range representation
		fmt.Printf("Range Breakdown:\n")
		fmt.Printf("  Min (worst):  %.4f\n", analysis.Min)
		fmt.Printf("  Q1 (25th):    %.4f\n", analysis.Percentiles[25])
		fmt.Printf("  Median:       %.4f\n", analysis.Median)
		fmt.Printf("  Q3 (75th):    %.4f\n", analysis.Percentiles[75])
		fmt.Printf("  Max (best):   %.4f\n", analysis.Max)
		fmt.Printf("  Mean:         %.4f\n", analysis.Mean)
		fmt.Printf("  Range span:   %.4f\n\n", scoreRange)
	} else {
		fmt.Printf("  All scores are identical: %.4f\n\n", analysis.Min)
	}

	// Quartile analysis
	if len(analysis.Quartiles) >= 3 {
		fmt.Printf("Quartile Analysis:\n")
		fmt.Printf("  Q1 (25th): %.4f\n", analysis.Quartiles[0])
		fmt.Printf("  Q2 (50th): %.4f (Median)\n", analysis.Quartiles[1])
		fmt.Printf("  Q3 (75th): %.4f\n", analysis.Quartiles[2])
		
		iqr := analysis.Quartiles[2] - analysis.Quartiles[0]
		fmt.Printf("  IQR:       %.4f\n", iqr)
		
		// Outlier boundaries
		lowerBound := analysis.Quartiles[0] - 1.5*iqr
		upperBound := analysis.Quartiles[2] + 1.5*iqr
		fmt.Printf("  Outliers:  < %.4f or > %.4f\n", lowerBound, upperBound)
	}

	if config.App.Verbose {
		fmt.Printf("\nDetailed Score List (first 20):\n")
		limit := 20
		if len(analysis.Scores) < limit {
			limit = len(analysis.Scores)
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("  %2d: %.4f\n", i+1, analysis.Scores[i])
		}
		if len(analysis.Scores) > 20 {
			fmt.Printf("  ... and %d more scores\n", len(analysis.Scores)-20)
		}
	}

	return nil
}