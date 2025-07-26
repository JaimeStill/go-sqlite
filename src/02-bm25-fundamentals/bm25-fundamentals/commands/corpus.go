package commands

import (
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/handlers"
	"github.com/spf13/cobra"
)

// Corpus is the public corpus command group instance
var Corpus = newCorpusGroup()

// newCorpusGroup creates the corpus command group with all its subcommands
func newCorpusGroup() *CommandGroup {
	// corpusCmd represents the corpus command group
	corpusCmd := &cobra.Command{
		Use:   "corpus",
		Short: "Manage document corpus for BM25 experiments",
		Long: `The corpus command group provides tools for managing the document collection
used in BM25 scoring experiments. This includes generating synthetic documents,
importing external content, and analyzing corpus statistics.

Understanding your corpus characteristics is crucial for BM25 analysis:
- Document length distribution affects score normalization
- Term frequency patterns influence ranking behavior
- Category diversity enables comparative analysis`,
	}

	// generateCmd generates a synthetic corpus
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate synthetic documents for BM25 experimentation",
		Long: `Generate a synthetic corpus with controlled characteristics for BM25 analysis.

The generated corpus includes:
- Varied document lengths to study BM25 length normalization
- Multiple categories for comparative ranking analysis  
- Realistic term distributions for scoring experiments
- Reproducible generation with optional seed parameter

Examples:
  # Generate default corpus (100 documents)
  bm25-fundamentals corpus generate
  
  # Generate specific size
  bm25-fundamentals corpus generate --size 500
  
  # Generate with custom categories
  bm25-fundamentals corpus generate --categories "tech,science,business"`,
		RunE: handlers.Corpus.HandleGenerate,
	}

	// statsCmd shows corpus statistics
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Display comprehensive corpus statistics",
		Long: `Show detailed statistics about the current corpus including:

- Document count and length distribution
- Category breakdown and diversity
- Term frequency characteristics
- Time range of document creation

These statistics help understand how BM25 scoring will behave with your corpus.`,
		RunE: handlers.Corpus.HandleStats,
	}

	// clearCmd removes all documents
	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Remove all documents from the corpus",
		Long: `Clear all documents from the corpus and optimize the FTS5 index.

This operation:
- Removes all documents from the database
- Clears the FTS5 search index
- Resets document ID sequence
- Optimizes index storage

Use this to start fresh with new corpus experiments.`,
		RunE: handlers.Corpus.HandleClear,
	}

	// setupFlags configures flags for corpus commands
	setupFlags := func() {
		// Generate command flags
		generateCmd.Flags().IntP("size", "s", 0, "number of documents to generate (0 = use config default)")
		generateCmd.Flags().StringP("categories", "c", "", "comma-separated list of categories")
		generateCmd.Flags().IntP("min-tokens", "m", 0, "minimum document length in tokens")
		generateCmd.Flags().IntP("max-tokens", "M", 0, "maximum document length in tokens")
		generateCmd.Flags().IntP("title-min-tokens", "", 0, "minimum title length in tokens")
		generateCmd.Flags().IntP("title-max-tokens", "", 0, "maximum title length in tokens")
		generateCmd.Flags().Int64P("seed", "", 0, "random seed for reproducible generation (0 = use current time)")

		// Clear command flags
		clearCmd.Flags().BoolP("confirm", "y", false, "confirm corpus deletion without prompt")
	}

	// Return the command group
	return &CommandGroup{
		Command: corpusCmd,
		SubCommands: []*cobra.Command{
			generateCmd,
			statsCmd,
			clearCmd,
		},
		FlagSetup: setupFlags,
	}
}
