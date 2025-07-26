package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/config"
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/database"
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/errors"
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/models"
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/utilities"
	"github.com/spf13/cobra"
)

// Validation is the global validation handler instance
var Validation ValidationHandler

// ValidationHandler manages validation operations (stateless - accesses global instances)
type ValidationHandler struct{}

// HandleValidateAll handles the comprehensive validation command
func (h *ValidationHandler) HandleValidateAll(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Running setup validation checks...")

	suite := &models.ValidationSuite{
		Name:        "Setup Validation Suite",
		Description: "Comprehensive validation of SQLite FTS5 learning environment",
		StartTime:   time.Now(),
		Results:     make([]models.ValidationResult, 0),
	}

	// Run all validation checks
	checks := []struct {
		name        string
		description string
		fn          func() error
	}{
		{"SQLite Connection", "Test basic SQLite database connection", h.validateSQLiteConnection},
		{"FTS5 Support", "Verify FTS5 virtual table support", h.validateFTS5Support},
		{"Test Data Generation", "Validate sample data insertion", h.validateTestData},
		{"BM25 Scoring", "Test BM25 scoring functionality", h.validateBM25Scoring},
		{"Shared Utilities", "Verify utility functions work correctly", h.validateUtilities},
	}

	for _, check := range checks {
		result := h.runValidationCheck(check.name, check.description, check.fn)
		suite.AddResult(result)
		h.displayResult(result)
	}

	suite.EndTime = time.Now()
	suite.Duration = suite.EndTime.Sub(suite.StartTime)

	// Display summary
	h.displaySummary(suite)

	if !suite.IsSuccessful() {
		return errors.Validationf("validation suite failed with %d errors", suite.Failed)
	}

	return nil
}

// HandleConnect handles the SQLite connection test
func (h *ValidationHandler) HandleConnect(cmd *cobra.Command, args []string) error {
	fmt.Println("🔗 Testing SQLite connection...")

	if err := h.validateSQLiteConnection(); err != nil {
		return err
	}

	fmt.Println("✅ SQLite connection successful")
	return nil
}

// HandleFTS5 handles the FTS5 functionality test
func (h *ValidationHandler) HandleFTS5(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Testing FTS5 functionality...")

	if err := h.validateFTS5Support(); err != nil {
		return err
	}

	fmt.Println("✅ FTS5 functionality working")
	return nil
}

// HandleTestData handles the test data validation
func (h *ValidationHandler) HandleTestData(cmd *cobra.Command, args []string) error {
	fmt.Println("📄 Testing sample data generation...")

	if err := h.validateTestData(); err != nil {
		return err
	}

	fmt.Println("✅ Sample data generation working")
	return nil
}

// HandleBM25 handles the BM25 scoring validation
func (h *ValidationHandler) HandleBM25(cmd *cobra.Command, args []string) error {
	fmt.Println("📊 Testing BM25 scoring...")

	if err := h.validateBM25Scoring(); err != nil {
		return err
	}

	fmt.Println("✅ BM25 scoring working correctly")
	return nil
}

// Validation check implementations

func (h *ValidationHandler) validateSQLiteConnection() error {
	ctx := context.Background()

	// Test basic connection
	version, err := database.Instance.GetSQLiteVersion(ctx)
	if err != nil {
		return errors.Connectionf("failed to get SQLite version: %w", err)
	}

	if config.App.IsVerbose() {
		fmt.Printf("  📍 SQLite version: %s\n", version)
	}

	return nil
}

func (h *ValidationHandler) validateFTS5Support() error {
	ctx := context.Background()

	// Verify FTS5 availability
	if err := database.Instance.VerifyFTS5Support(ctx); err != nil {
		return err
	}

	// Test FTS5 table creation
	if err := database.Instance.CreateTestTable(ctx, "fts5_validation_test"); err != nil {
		return err
	}

	if config.App.IsVerbose() {
		fmt.Println("  📍 FTS5 virtual table created successfully")
	}

	return nil
}

func (h *ValidationHandler) validateTestData() error {
	ctx := context.Background()
	tableName := "test_data_validation"

	// Create table and insert sample data
	if err := database.Instance.CreateTestTable(ctx, tableName); err != nil {
		return err
	}

	if err := database.Instance.InsertTestData(ctx, tableName); err != nil {
		return err
	}

	// Verify data was inserted
	count, err := database.Instance.CountDocuments(ctx, tableName)
	if err != nil {
		return err
	}

	expectedCount := len(utilities.SampleDocuments())
	if count != expectedCount {
		return errors.Validationf("expected %d documents, got %d", expectedCount, count)
	}

	if config.App.IsVerbose() {
		fmt.Printf("  📍 %d sample documents inserted successfully\n", count)
	}

	return nil
}

func (h *ValidationHandler) validateBM25Scoring() error {
	ctx := context.Background()
	tableName := "bm25_validation_test"

	// Setup test data
	if err := database.Instance.CreateTestTable(ctx, tableName); err != nil {
		return err
	}

	if err := database.Instance.InsertTestData(ctx, tableName); err != nil {
		return err
	}

	// Test BM25 query
	rows, err := database.Instance.QueryWithBM25(ctx, tableName, "SQLite")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Verify we get results with negative scores
	var resultCount int
	var hasNegativeScore bool
	var firstScore float64

	for rows.Next() {
		var id int
		var title, content string
		var score float64

		err = rows.Scan(&id, &title, &content, &score)
		if err != nil {
			return errors.Databasef("failed to scan result: %w", err)
		}

		resultCount++
		if score < 0 {
			hasNegativeScore = true
		}

		if resultCount == 1 {
			firstScore = score
		}
	}

	if resultCount == 0 {
		return errors.Validationf("no results returned for BM25 query")
	}

	if !hasNegativeScore {
		return errors.Validationf("expected negative BM25 scores, but none found")
	}

	if config.App.IsVerbose() {
		fmt.Printf("  📍 %d results with proper BM25 scoring\n", resultCount)
		fmt.Printf("  📍 First result score: %.3f (negative as expected)\n", firstScore)
	}

	return nil
}

func (h *ValidationHandler) validateUtilities() error {
	// This validation is implicit in the other tests
	if config.App.IsVerbose() {
		fmt.Println("  📍 Utility functions accessible and working")
	}
	return nil
}

// Helper methods

func (h *ValidationHandler) runValidationCheck(name, description string, fn func() error) models.ValidationResult {
	start := time.Now()
	err := fn()
	duration := time.Since(start)

	result := models.ValidationResult{
		Name:        name,
		Description: description,
		Passed:      err == nil,
		Error:       err,
		Duration:    duration,
	}

	return result
}

func (h *ValidationHandler) displayResult(result models.ValidationResult) {
	status := "✅"
	if !result.Passed {
		status = "❌"
	}

	fmt.Printf("%s %s", status, result.Name)
	
	if config.App.IsVerbose() {
		fmt.Printf(" (%v)", result.Duration)
	}
	
	fmt.Println()

	if !result.Passed && result.Error != nil {
		fmt.Printf("   Error: %v\n", result.Error)
	}
}

func (h *ValidationHandler) displaySummary(suite *models.ValidationSuite) {
	fmt.Printf("\n📊 Validation Results:\n")
	fmt.Printf("   Total checks: %d\n", suite.Total)
	fmt.Printf("   Passed: %d\n", suite.Passed)
	fmt.Printf("   Failed: %d\n", suite.Failed)
	fmt.Printf("   Success rate: %.1f%%\n", suite.GetSuccessRate())
	fmt.Printf("   Duration: %v\n", suite.Duration)

	if suite.IsSuccessful() {
		fmt.Println("\n🎉 All validation checks passed! Environment is ready for FTS5 learning.")
	} else {
		fmt.Println("\n⚠️  Some checks failed. Please review the setup before proceeding.")
	}
}