package commands

import (
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/handlers"
	"github.com/spf13/cobra"
)

// Visualize is the public visualize command group instance
var Visualize = newVisualizeGroup()

// newVisualizeGroup creates the visualize command group with all its subcommands
func newVisualizeGroup() *CommandGroup {
	// visualizeCmd represents the visualize command group
	visualizeCmd := &cobra.Command{
		Use:   "visualize",
		Short: "Visualize BM25 score distributions and analysis results",
		Long: `The visualize command group provides ASCII-based visualizations for BM25 analysis.

Key Features:
- Score distribution histograms
- Category comparison charts  
- Percentile range visualizations
- Ranking difference analysis

All visualizations are designed for terminal display and provide educational
insights into BM25 scoring patterns and search result distributions.`,
	}

	// distributionCmd shows score distribution histogram
	distributionCmd := &cobra.Command{
		Use:   "distribution",
		Short: "Display score distribution histogram for search results",
		Long: `Generate an ASCII histogram showing the distribution of BM25 scores for search results.

This visualization helps understand:
- Score concentration patterns
- Outliers and score ranges
- Distribution shape (normal, skewed, etc.)
- Relative frequency of different score ranges

Examples:
  # Basic score distribution
  bm25-fundamentals visualize distribution --query "database optimization"
  
  # Distribution with custom weights
  bm25-fundamentals visualize distribution --query "algorithm" --title-weight 2.0`,
		RunE: handlers.Visualize.HandleDistribution,
	}

	// categoriesCmd compares scores across categories
	categoriesCmd := &cobra.Command{
		Use:   "categories",
		Short: "Compare BM25 score distributions across document categories",
		Long: `Generate category-wise score comparison charts showing how different
document categories perform for a given search query.

This visualization reveals:
- Which categories tend to score higher/lower
- Category-specific score ranges and patterns
- Relative performance between content types
- Potential bias in search results

Examples:
  # Compare categories for a query
  bm25-fundamentals visualize categories --query "optimization performance"
  
  # Focus on specific categories
  bm25-fundamentals visualize categories --query "database" --filter "technology,science"`,
		RunE: handlers.Visualize.HandleCategories,
	}

	// rangeCmd shows score range and percentile visualization
	rangeCmd := &cobra.Command{
		Use:   "range",
		Short: "Display score range and percentile visualization",
		Long: `Generate a visual representation of score ranges, percentiles, and distribution statistics.

This shows:
- Score range (min to max)
- Percentile markers (25th, 50th, 75th, 90th, 95th, 99th)
- Mean and median indicators
- Outlier identification

Examples:
  # Score range analysis
  bm25-fundamentals visualize range --query "machine learning algorithms"
  
  # Range with statistical details
  bm25-fundamentals visualize range --query "database" --verbose`,
		RunE: handlers.Visualize.HandleRange,
	}

	// setupFlags configures flags for visualize commands
	setupFlags := func() {
		// Distribution command flags
		distributionCmd.Flags().StringP("query", "q", "", "search query (required)")
		distributionCmd.Flags().StringP("category", "c", "", "filter by category")
		distributionCmd.Flags().Float64P("title-weight", "", 0, "title field weight (0 = default)")
		distributionCmd.Flags().Float64P("content-weight", "", 0, "content field weight (0 = default)")
		distributionCmd.Flags().Float64P("category-weight", "", 0, "category field weight (0 = default)")
		distributionCmd.Flags().IntP("buckets", "b", 10, "number of histogram buckets")
		distributionCmd.Flags().IntP("max-results", "n", 100, "maximum results to analyze")
		distributionCmd.MarkFlagRequired("query")

		// Categories command flags
		categoriesCmd.Flags().StringP("query", "q", "", "search query (required)")
		categoriesCmd.Flags().String("filter", "", "comma-separated list of categories to include")
		categoriesCmd.Flags().Float64P("title-weight", "", 0, "title field weight (0 = default)")
		categoriesCmd.Flags().Float64P("content-weight", "", 0, "content field weight (0 = default)")
		categoriesCmd.Flags().Float64P("category-weight", "", 0, "category field weight (0 = default)")
		categoriesCmd.Flags().IntP("max-results", "n", 100, "maximum results to analyze")
		categoriesCmd.MarkFlagRequired("query")

		// Range command flags
		rangeCmd.Flags().StringP("query", "q", "", "search query (required)")
		rangeCmd.Flags().StringP("category", "c", "", "filter by category")
		rangeCmd.Flags().Float64P("title-weight", "", 0, "title field weight (0 = default)")
		rangeCmd.Flags().Float64P("content-weight", "", 0, "content field weight (0 = default)")
		rangeCmd.Flags().Float64P("category-weight", "", 0, "category field weight (0 = default)")
		rangeCmd.Flags().IntP("max-results", "n", 100, "maximum results to analyze")
		rangeCmd.MarkFlagRequired("query")
	}

	// Return the command group
	return &CommandGroup{
		Command: visualizeCmd,
		SubCommands: []*cobra.Command{
			distributionCmd,
			categoriesCmd,
			rangeCmd,
		},
		FlagSetup: setupFlags,
	}
}