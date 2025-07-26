package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/config"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/database"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/errors"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/models"
	"github.com/spf13/cobra"
)

// Corpus is the global corpus handler instance
var Corpus CorpusHandler

// CorpusHandler manages corpus operations (stateless - accesses global instances)
type CorpusHandler struct{}

// HandleGenerate handles the corpus generate command
func (h *CorpusHandler) HandleGenerate(cmd *cobra.Command, args []string) error {
	// Extract flags
	size, _ := cmd.Flags().GetInt("size")
	categories, _ := cmd.Flags().GetString("categories")
	minTokens, _ := cmd.Flags().GetInt("min-tokens")
	maxTokens, _ := cmd.Flags().GetInt("max-tokens")
	titleMinTokens, _ := cmd.Flags().GetInt("title-min-tokens")
	titleMaxTokens, _ := cmd.Flags().GetInt("title-max-tokens")
	seed, _ := cmd.Flags().GetInt64("seed")
	confirmClear, _ := cmd.Flags().GetBool("confirm")

	// Start with default options
	options := models.DefaultCorpusOptions()

	// Apply size from flag or config
	if size > 0 {
		options.Size = size
	} else {
		options.Size = config.App.Corpus.Size
	}

	if options.Size < 1 {
		return errors.Validationf("corpus size must be at least 1")
	}

	// Apply flag overrides
	if categories != "" {
		options.Categories = strings.Split(categories, ",")
		for i, cat := range options.Categories {
			options.Categories[i] = strings.TrimSpace(cat)
		}
	}

	if minTokens > 0 {
		options.MinTokens = minTokens
	}

	if maxTokens > 0 {
		options.MaxTokens = maxTokens
	}

	if titleMinTokens > 0 {
		options.TitleMinTokens = titleMinTokens
	}

	if titleMaxTokens > 0 {
		options.TitleMaxTokens = titleMaxTokens
	}

	if seed != 0 {
		options.Seed = seed
	}

	// Validate configuration
	if options.MinTokens >= options.MaxTokens {
		return errors.Validationf("min-tokens (%d) must be less than max-tokens (%d)",
			options.MinTokens, options.MaxTokens)
	}

	ctx := context.Background()

	// Initialize schema
	if err := database.Instance.InitSchema(ctx); err != nil {
		return err
	}

	// Check if corpus already exists
	existingCount, err := h.GetDocumentCount(ctx)
	if err != nil {
		return err
	}

	if existingCount > 0 && !confirmClear {
		fmt.Printf("Corpus already contains %d documents.\n", existingCount)
		fmt.Print("Do you want to clear the existing corpus? (y/N): ")

		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			fmt.Println("Corpus generation cancelled.")
			return nil
		}

		if err := h.ClearDocuments(ctx); err != nil {
			return err
		}
		fmt.Println("Existing corpus cleared.")
	}

	// Generate corpus
	fmt.Printf("Generating corpus with %d documents...\n", options.Size)
	if config.App.Verbose {
		fmt.Printf("Configuration:\n")
		fmt.Printf("  Categories: %v\n", options.Categories)
		fmt.Printf("  Document length: %d-%d tokens\n", options.MinTokens, options.MaxTokens)
		fmt.Printf("  Title length: %d-%d tokens\n", options.TitleMinTokens, options.TitleMaxTokens)
		if options.Seed != 0 {
			fmt.Printf("  Random seed: %d\n", options.Seed)
		}
	}

	if err := h.GenerateCorpus(ctx, options); err != nil {
		return err
	}

	// Show results
	finalCount, err := h.GetDocumentCount(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Successfully generated %d documents\n", finalCount)

	if config.App.Verbose {
		// Show quick stats
		stats, err := h.GetCorpusStats(ctx)
		if err == nil {
			fmt.Printf("\nCorpus Statistics:\n")
			fmt.Printf("  Average document length: %.1f tokens\n", stats.AverageDocLength)
			fmt.Printf("  Categories: %d\n", len(stats.Categories))
			for _, cat := range stats.Categories {
				fmt.Printf("    %s: %d documents\n", cat, stats.CategoryCounts[cat])
			}
		}
	}

	return nil
}

// HandleStats handles the corpus stats command
func (h *CorpusHandler) HandleStats(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get corpus statistics
	stats, err := h.GetCorpusStats(ctx)
	if err != nil {
		return err
	}

	if stats.TotalDocuments == 0 {
		fmt.Println("No documents in corpus. Use 'corpus generate' to create a synthetic corpus.")
		return nil
	}

	// Display statistics based on format
	switch config.App.Format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(stats)

	case "csv":
		fmt.Println("metric,value")
		fmt.Printf("total_documents,%d\n", stats.TotalDocuments)
		fmt.Printf("total_tokens,%d\n", stats.TotalTokens)
		fmt.Printf("average_doc_length,%.2f\n", stats.AverageDocLength)
		fmt.Printf("median_doc_length,%.2f\n", stats.MedianDocLength)
		fmt.Printf("min_doc_length,%d\n", stats.MinDocLength)
		fmt.Printf("max_doc_length,%d\n", stats.MaxDocLength)
		fmt.Printf("unique_terms,%d\n", stats.UniqueTerms)
		fmt.Printf("categories,%d\n", len(stats.Categories))

	default: // text format
		fmt.Printf("Corpus Statistics\n")
		fmt.Printf("=================\n\n")

		fmt.Printf("Document Count: %d\n", stats.TotalDocuments)
		fmt.Printf("Total Tokens: %d\n", stats.TotalTokens)
		fmt.Printf("Unique Terms: %d\n", stats.UniqueTerms)
		fmt.Printf("\n")

		fmt.Printf("Document Length Distribution:\n")
		fmt.Printf("  Average: %.1f tokens\n", stats.AverageDocLength)
		fmt.Printf("  Median:  %.1f tokens\n", stats.MedianDocLength)
		fmt.Printf("  Range:   %d - %d tokens\n", stats.MinDocLength, stats.MaxDocLength)
		fmt.Printf("\n")

		if len(stats.Categories) > 0 {
			fmt.Printf("Categories (%d):\n", len(stats.Categories))
			for _, cat := range stats.Categories {
				count := stats.CategoryCounts[cat]
				percentage := float64(count) * 100.0 / float64(stats.TotalDocuments)
				fmt.Printf("  %-15s: %5d documents (%.1f%%)\n", cat, count, percentage)
			}
			fmt.Printf("\n")
		}

		if !stats.CreatedRange.Start.IsZero() {
			fmt.Printf("Creation Time Range:\n")
			fmt.Printf("  From: %s\n", stats.CreatedRange.Start.Format("2006-01-02 15:04:05"))
			fmt.Printf("  To:   %s\n", stats.CreatedRange.End.Format("2006-01-02 15:04:05"))
			fmt.Printf("\n")
		}

		fmt.Printf("Last Updated: %s\n", stats.LastUpdated.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// HandleClear handles the corpus clear command
func (h *CorpusHandler) HandleClear(cmd *cobra.Command, args []string) error {
	confirmClear, _ := cmd.Flags().GetBool("confirm")

	ctx := context.Background()

	// Check current document count
	count, err := h.GetDocumentCount(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		fmt.Println("Corpus is already empty.")
		return nil
	}

	// Confirm deletion
	if !confirmClear {
		fmt.Printf("This will delete all %d documents from the corpus.\n", count)
		fmt.Print("Are you sure? (y/N): ")

		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Clear the corpus
	fmt.Printf("Clearing %d documents...\n", count)
	if err := h.ClearDocuments(ctx); err != nil {
		return err
	}

	fmt.Println("✓ Corpus cleared successfully")
	return nil
}

// InsertDocument adds a single document to the corpus
func (h *CorpusHandler) InsertDocument(ctx context.Context, doc *models.Document) error {
	// Calculate document length in tokens (simple whitespace tokenization)
	doc.Length = len(strings.Fields(doc.Title + " " + doc.Content))

	query := `
		INSERT INTO documents (title, content, category, length, created) 
		VALUES (?, ?, ?, ?, ?)
		RETURNING id`

	err := database.Instance.DB().QueryRowContext(ctx, query,
		doc.Title, doc.Content, doc.Category, doc.Length, doc.Created,
	).Scan(&doc.ID)

	if err != nil {
		return errors.Databasef("failed to insert document: %w", err)
	}

	return nil
}

// BatchInsertDocuments efficiently inserts multiple documents
func (h *CorpusHandler) BatchInsertDocuments(ctx context.Context, docs []*models.Document) error {
	if len(docs) == 0 {
		return nil
	}

	tx, err := database.Instance.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO documents (title, content, category, length, created) 
		 VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return errors.Databasef("failed to prepare batch insert: %w", err)
	}
	defer stmt.Close()

	for _, doc := range docs {
		// Calculate document length
		doc.Length = len(strings.Fields(doc.Title + " " + doc.Content))

		if _, err := stmt.ExecContext(ctx,
			doc.Title, doc.Content, doc.Category, doc.Length, doc.Created); err != nil {
			return errors.Databasef("failed to insert document in batch: %w", err)
		}
	}

	return tx.Commit()
}

// GetDocumentCount returns the total number of documents
func (h *CorpusHandler) GetDocumentCount(ctx context.Context) (int, error) {
	var count int
	err := database.Instance.DB().QueryRowContext(ctx, "SELECT COUNT(*) FROM documents").Scan(&count)
	if err != nil {
		return 0, errors.Databasef("failed to get document count: %w", err)
	}
	return count, nil
}

// ClearDocuments removes all documents from the corpus
func (h *CorpusHandler) ClearDocuments(ctx context.Context) error {
	tx, err := database.Instance.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear documents table (triggers will handle FTS5 cleanup)
	if _, err := tx.ExecContext(ctx, "DELETE FROM documents"); err != nil {
		return errors.Databasef("failed to clear documents: %w", err)
	}

	// Reset auto-increment counter
	if _, err := tx.ExecContext(ctx, "DELETE FROM sqlite_sequence WHERE name='documents'"); err != nil {
		// This might fail if no auto-increment has occurred yet, which is fine
	}

	// Optimize FTS5 index
	if _, err := tx.ExecContext(ctx, "INSERT INTO documents_fts(documents_fts) VALUES('optimize')"); err != nil {
		return errors.FTS5f("failed to optimize FTS5 index: %w", err)
	}

	return tx.Commit()
}

// GetCorpusStats calculates comprehensive statistics about the corpus
func (h *CorpusHandler) GetCorpusStats(ctx context.Context) (*models.CorpusStats, error) {
	stats := &models.CorpusStats{
		CategoryCounts: make(map[string]int),
		LastUpdated:    time.Now(),
	}

	// Get basic counts and length statistics
	query := `
		SELECT 
			COUNT(*) as total_docs,
			SUM(length) as total_tokens,
			AVG(length) as avg_length,
			MIN(length) as min_length,
			MAX(length) as max_length,
			MIN(created) as earliest,
			MAX(created) as latest
		FROM documents`

	var earliest, latest sql.NullString
	err := database.Instance.DB().QueryRowContext(ctx, query).Scan(
		&stats.TotalDocuments,
		&stats.TotalTokens,
		&stats.AverageDocLength,
		&stats.MinDocLength,
		&stats.MaxDocLength,
		&earliest,
		&latest,
	)

	if err != nil {
		return nil, errors.Databasef("failed to get basic corpus stats: %w", err)
	}

	if earliest.Valid {
		if t, err := time.Parse("2006-01-02 15:04:05", earliest.String); err == nil {
			stats.CreatedRange.Start = t
		}
	}
	if latest.Valid {
		if t, err := time.Parse("2006-01-02 15:04:05", latest.String); err == nil {
			stats.CreatedRange.End = t
		}
	}

	// Get median document length
	medianQuery := `
		SELECT length 
		FROM documents 
		ORDER BY length 
		LIMIT 1 
		OFFSET (SELECT (COUNT(*) - 1) / 2 FROM documents)`

	err = database.Instance.DB().QueryRowContext(ctx, medianQuery).Scan(&stats.MedianDocLength)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Databasef("failed to get median document length: %w", err)
	}

	// Get category breakdown
	categoryQuery := `
		SELECT category, COUNT(*) 
		FROM documents 
		GROUP BY category 
		ORDER BY COUNT(*) DESC`

	rows, err := database.Instance.DB().QueryContext(ctx, categoryQuery)
	if err != nil {
		return nil, errors.Databasef("failed to get category breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, errors.Databasef("failed to scan category data: %w", err)
		}
		stats.Categories = append(stats.Categories, category)
		stats.CategoryCounts[category] = count
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Databasef("error iterating category data: %w", err)
	}

	// Get unique terms count (approximate)
	uniqueTermsQuery := `
		SELECT COUNT(DISTINCT term) 
		FROM documents_fts_data 
		WHERE col = '*'`

	err = database.Instance.DB().QueryRowContext(ctx, uniqueTermsQuery).Scan(&stats.UniqueTerms)
	if err != nil {
		// If this fails, it's not critical - just set to 0
		stats.UniqueTerms = 0
	}

	return stats, nil
}

// GenerateCorpus creates a synthetic corpus for BM25 experimentation
func (h *CorpusHandler) GenerateCorpus(ctx context.Context, options models.CorpusOptions) error {
	// Set up random seed for reproducible generation
	if options.Seed == 0 {
		options.Seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(options.Seed))

	// Generate documents
	docs := make([]*models.Document, 0, options.Size)
	generator := &corpusGenerator{
		rng:     rng,
		options: options,
	}

	for i := 0; i < options.Size; i++ {
		doc := generator.generateDocument()
		docs = append(docs, doc)
	}

	// Insert in batches for efficiency
	return h.BatchInsertDocuments(ctx, docs)
}

// corpusGenerator handles synthetic document generation
type corpusGenerator struct {
	rng     *rand.Rand
	options models.CorpusOptions
}

// generateDocument creates a single synthetic document
func (g *corpusGenerator) generateDocument() *models.Document {
	category := g.options.Categories[g.rng.Intn(len(g.options.Categories))]

	title := g.generateTitle(category)
	content := g.generateContent(category)

	// Random creation time within the last 30 days
	createdOffset := time.Duration(g.rng.Intn(30*24*60)) * time.Minute
	created := time.Now().Add(-createdOffset)

	return &models.Document{
		Title:    title,
		Content:  content,
		Category: category,
		Created:  created,
	}
}

// generateTitle creates a synthetic title based on category
func (g *corpusGenerator) generateTitle(category string) string {
	templates := map[string][]string{
		"technology": {
			"Advanced %s Development Techniques",
			"Understanding %s Architecture",
			"Modern %s Best Practices",
			"Introduction to %s Programming",
			"%s Performance Optimization",
		},
		"science": {
			"Research in %s Methods",
			"Scientific Analysis of %s",
			"Experimental %s Studies",
			"Theoretical %s Frameworks",
			"Applications of %s Theory",
		},
		"programming": {
			"Mastering %s Algorithms",
			"Efficient %s Implementation",
			"Advanced %s Patterns",
			"Learning %s Programming",
			"%s Code Optimization",
		},
		"database": {
			"Optimizing %s Queries",
			"Advanced %s Indexing",
			"%s Transaction Management",
			"Scaling %s Systems",
			"%s Performance Tuning",
		},
		"algorithms": {
			"Efficient %s Algorithms",
			"Complex %s Analysis",
			"Optimized %s Solutions",
			"Advanced %s Techniques",
			"Comparative %s Study",
		},
	}

	terms := map[string][]string{
		"technology":  {"Cloud", "Mobile", "Web", "AI", "Blockchain", "IoT"},
		"science":     {"Data", "Machine Learning", "Statistics", "Analytics", "Research"},
		"programming": {"Object-Oriented", "Functional", "Concurrent", "Distributed", "Reactive"},
		"database":    {"SQL", "NoSQL", "Relational", "Graph", "Time-Series"},
		"algorithms":  {"Sorting", "Search", "Graph", "Dynamic Programming", "Greedy"},
	}

	categoryTemplates := templates[category]
	if len(categoryTemplates) == 0 {
		categoryTemplates = templates["technology"] // fallback
	}

	categoryTerms := terms[category]
	if len(categoryTerms) == 0 {
		categoryTerms = terms["technology"] // fallback
	}

	template := categoryTemplates[g.rng.Intn(len(categoryTemplates))]
	term := categoryTerms[g.rng.Intn(len(categoryTerms))]

	return fmt.Sprintf(template, term)
}

// generateContent creates synthetic content with controlled length
func (g *corpusGenerator) generateContent(category string) string {
	targetTokens := g.options.MinTokens + g.rng.Intn(g.options.MaxTokens-g.options.MinTokens+1)

	// Base vocabulary for different categories
	vocabulary := map[string][]string{
		"technology":  {"system", "development", "architecture", "framework", "platform", "solution", "design", "implementation", "scalable", "efficient", "robust", "secure", "modern", "advanced", "innovative"},
		"science":     {"research", "analysis", "methodology", "hypothesis", "experiment", "data", "results", "conclusion", "theory", "evidence", "statistical", "empirical", "quantitative", "qualitative", "validation"},
		"programming": {"function", "variable", "algorithm", "optimization", "performance", "debugging", "testing", "refactoring", "maintainable", "readable", "efficient", "scalable", "object", "method", "interface"},
		"database":    {"query", "index", "transaction", "optimization", "performance", "schema", "normalization", "relational", "primary", "foreign", "key", "table", "column", "constraint", "integrity"},
		"algorithms":  {"complexity", "efficiency", "optimization", "iteration", "recursion", "sorting", "searching", "traversal", "comparison", "analysis", "space", "time", "linear", "logarithmic", "polynomial"},
	}

	// Common connecting words
	connectors := []string{"and", "the", "of", "in", "to", "for", "with", "by", "from", "on", "at", "as", "is", "are", "can", "will", "this", "that", "these", "those"}

	words := make([]string, 0, targetTokens)
	categoryWords := vocabulary[category]
	if len(categoryWords) == 0 {
		categoryWords = vocabulary["technology"] // fallback
	}

	for len(words) < targetTokens {
		if g.rng.Float64() < 0.3 { // 30% chance for category-specific word
			word := categoryWords[g.rng.Intn(len(categoryWords))]
			words = append(words, word)
		} else { // 70% chance for connector word
			word := connectors[g.rng.Intn(len(connectors))]
			words = append(words, word)
		}
	}

	return strings.Join(words, " ")
}
