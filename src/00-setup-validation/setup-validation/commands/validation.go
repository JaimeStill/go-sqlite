package commands

import (
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/handlers"
	"github.com/spf13/cobra"
)

// Validation is the public validation command group instance
var Validation = newValidationGroup()

// newValidationGroup creates the validation command group with all its subcommands
func newValidationGroup() *CommandGroup {
	// validateCmd represents the comprehensive validation command
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Run all validation checks",
		Long: `Runs a comprehensive suite of validation checks for the learning environment.

This command tests:
- SQLite database connection
- FTS5 virtual table support
- Sample data generation and insertion  
- BM25 scoring functionality
- Utility function accessibility

All checks must pass for the environment to be considered ready for FTS5 learning.`,
		RunE: handlers.Validation.HandleValidateAll,
	}

	// connectCmd tests SQLite connection
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Test SQLite database connection",
		Long:  "Validates that we can establish a connection to SQLite and query basic information.",
		RunE:  handlers.Validation.HandleConnect,
	}

	// fts5Cmd tests FTS5 functionality
	fts5Cmd := &cobra.Command{
		Use:   "fts5",
		Short: "Test FTS5 functionality",
		Long:  "Validates FTS5 virtual table creation and basic operations to ensure FTS5 support is available.",
		RunE:  handlers.Validation.HandleFTS5,
	}

	// testDataCmd tests sample data generation
	testDataCmd := &cobra.Command{
		Use:   "testdata",
		Short: "Test sample data generation",
		Long:  "Validates the utilities for generating and inserting test data into FTS5 tables.",
		RunE:  handlers.Validation.HandleTestData,
	}

	// bm25Cmd tests BM25 scoring
	bm25Cmd := &cobra.Command{
		Use:   "bm25",
		Short: "Test BM25 scoring",
		Long:  "Validates BM25 scoring functionality and SQLite's negative scoring behavior for proper FTS5 integration.",
		RunE:  handlers.Validation.HandleBM25,
	}

	// setupFlags configures flags for validation commands
	setupFlags := func() {
		// No specific flags needed for validation commands
		// They inherit global flags (verbose, format) from root
	}

	// Return the command group
	return &CommandGroup{
		Command: &cobra.Command{
			Use:   "validation",
			Short: "Validation commands",
			Long:  "Commands for validating SQLite FTS5 setup and functionality",
		},
		SubCommands: []*cobra.Command{
			validateCmd,
			connectCmd,
			fts5Cmd,
			testDataCmd,
			bm25Cmd,
		},
		FlagSetup: setupFlags,
	}
}