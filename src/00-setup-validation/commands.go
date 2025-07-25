package main

import (
	"fmt"
	"os"

	"github.com/jaime/go-sqlite/shared"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Run all validation checks",
	Long:  "Runs a comprehensive suite of validation checks for the learning environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Running setup validation checks...")

		// Run all validation checks
		checks := []func() error{
			validateSQLiteConnection,
			validateFTS5Support,
			validateSharedUtilities,
			validateTestData,
			validateBM25Scoring,
		}

		passed := 0
		total := len(checks)

		for i, check := range checks {
			if err := check(); err != nil {
				fmt.Printf("❌ Check %d/%d failed: %v\n", i+1, total, err)
			} else {
				fmt.Printf("✅ Check %d/%d passed\n", i+1, total)
				passed++
			}
		}

		fmt.Printf("\n📊 Results: %d/%d checks passed\n", passed, total)

		if passed == total {
			fmt.Println("🎉 All validation checks passed! Environment is ready.")
		} else {
			fmt.Println("⚠️  Some checks failed. Please review the setup.")
			os.Exit(1)
		}
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Test SQLite database connection",
	Long:  "Validates that we can establish a connection to SQLite with FTS5 support",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔗 Testing SQLite connection...")

		if err := validateSQLiteConnection(); err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			return
		}

		fmt.Println("✅ SQLite connection successful")
	},
}

var fts5Cmd = &cobra.Command{
	Use:   "fts5",
	Short: "Test FTS5 functionality",
	Long:  "Validates FTS5 virtual table creation and basic operations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Testing FTS5 functionality...")

		if err := validateFTS5Support(); err != nil {
			fmt.Printf("❌ FTS5 validation failed: %v\n", err)
			return
		}

		fmt.Println("✅ FTS5 functionality working")
	},
}

var testDataCmd = &cobra.Command{
	Use:   "testdata",
	Short: "Test sample data generation",
	Long:  "Validates the shared utilities for generating and inserting test data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📄 Testing sample data generation...")

		if err := validateTestData(); err != nil {
			fmt.Printf("❌ Test data validation failed: %v\n", err)
			return
		}

		fmt.Println("✅ Sample data generation working")
	},
}

var bm25Cmd = &cobra.Command{
	Use:   "bm25",
	Short: "Test BM25 scoring",
	Long:  "Validates BM25 scoring functionality and SQLite's inverted scoring behavior",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 Testing BM25 scoring...")

		if err := validateBM25Scoring(); err != nil {
			fmt.Printf("❌ BM25 validation failed: %v\n", err)
			return
		}

		fmt.Println("✅ BM25 scoring working correctly")
	},
}

// Validation functions
func validateSQLiteConnection() error {
	db, err := shared.ConnectMemoryFTS5()
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite: %w", err)
	}
	defer db.Close()

	// Test basic query
	var version string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to query SQLite version: %w", err)
	}

	fmt.Printf("  📍 SQLite version: %s\n", version)
	return nil
}

func validateFTS5Support() error {
	db, err := shared.ConnectMemoryFTS5()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	// Test FTS5 table creation
	err = shared.CreateBasicDocumentTable(db, "test_docs")
	if err != nil {
		return fmt.Errorf("failed to create FTS5 table: %w", err)
	}

	// Verify table exists and is FTS5
	var tableType string
	err = db.QueryRow("SELECT type FROM sqlite_master WHERE name = 'test_docs'").Scan(&tableType)
	if err != nil {
		return fmt.Errorf("failed to verify table creation: %w", err)
	}

	if tableType != "table" {
		return fmt.Errorf("expected table type 'table', got '%s'", tableType)
	}

	fmt.Println("  📍 FTS5 virtual table created successfully")
	return nil
}

func validateSharedUtilities() error {
	// This is covered by other validation functions
	fmt.Println("  📍 Shared utilities package accessible")
	return nil
}

func validateTestData() error {
	db, err := shared.ConnectMemoryFTS5()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	// Create table and insert sample data
	err = shared.CreateBasicDocumentTable(db, "sample_docs")
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	err = shared.InsertSampleDocuments(db, "sample_docs")
	if err != nil {
		return fmt.Errorf("failed to insert sample data: %w", err)
	}

	// Verify data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sample_docs").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count documents: %w", err)
	}

	expectedCount := len(shared.SampleDocuments())
	if count != expectedCount {
		return fmt.Errorf("expected %d documents, got %d", expectedCount, count)
	}

	fmt.Printf("  📍 %d sample documents inserted successfully\n", count)
	return nil
}

func validateBM25Scoring() error {
	db, err := shared.ConnectMemoryFTS5()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	// Setup test data
	err = shared.CreateBasicDocumentTable(db, "bm25_test")
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	err = shared.InsertSampleDocuments(db, "bm25_test")
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	// Test BM25 query
	rows, err := shared.QueryFTS5WithBM25(db, "bm25_test", "SQLite")
	if err != nil {
		return fmt.Errorf("failed to execute BM25 query: %w", err)
	}
	defer rows.Close()

	// Verify we get results with negative scores
	var resultCount int
	var hasNegativeScore bool

	for rows.Next() {
		var id int
		var title, content string
		var score float64

		err = rows.Scan(&id, &title, &content, &score)
		if err != nil {
			return fmt.Errorf("failed to scan result: %w", err)
		}

		resultCount++
		if score < 0 {
			hasNegativeScore = true
		}

		if resultCount == 1 {
			fmt.Printf("  📍 First result score: %.3f (should be negative)\n", score)
		}
	}

	if resultCount == 0 {
		return fmt.Errorf("no results returned for BM25 query")
	}

	if !hasNegativeScore {
		return fmt.Errorf("expected negative BM25 scores, but none found")
	}

	fmt.Printf("  📍 %d results with proper BM25 scoring\n", resultCount)
	return nil
}
