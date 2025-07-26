package commands

import (
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/handlers"
	"github.com/spf13/cobra"
)

// Search is the public search command group instance
var Search = newSearchGroup()

// newSearchGroup creates the search command group with all its subcommands
func newSearchGroup() *CommandGroup {
	// searchCmd represents the search command group
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search documents using BM25 full-text search",
		Long: `The search command group provides BM25-powered full-text search capabilities
with detailed score analysis and relevance insights.

Key Features:
- BM25 relevance scoring with SQLite FTS5
- Column weighting for multi-field ranking
- Score distribution analysis and explanation  
- Category filtering and result statistics
- Snippet generation with term highlighting

Understanding BM25 scores:
- SQLite FTS5 returns NEGATIVE scores (lower = better match)
- Scores around -1.0 indicate excellent relevance
- Scores below -4.0 typically indicate poor relevance`,
	}

	// queryCmd performs basic BM25 search
	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Search documents with BM25 scoring",
		Long: `Search the document corpus using BM25 full-text search.

Examples:
  # Basic search
  bm25-fundamentals search query --query "database optimization"
  
  # Search with custom column weights (title=2x, content=1x, category=0.5x)
  bm25-fundamentals search query --query "database" --title-weight 2.0 --content-weight 1.0 --category-weight 0.5
  
  # Filter by category
  bm25-fundamentals search query --query "algorithm" --category "programming"
  
  # Show detailed results with snippets
  bm25-fundamentals search query --query "optimization" --max-results 10 --snippets`,
		RunE: handlers.Search.HandleQuery,
	}

	// statsCmd shows search result statistics
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Analyze BM25 score distribution and search statistics",
		Long: `Analyze the statistical distribution of BM25 scores for a search query.

This command provides insights into:
- Score range and distribution characteristics
- Percentile analysis (25th, 50th, 75th, 90th, 95th, 99th)
- Category breakdown of results
- Score histogram buckets for visualization

Examples:
  # Basic score analysis
  bm25-fundamentals search stats --query "database optimization"
  
  # Export statistics as JSON
  bm25-fundamentals search stats --query "algorithm" --format json`,
		RunE: handlers.Search.HandleStats,
	}

	// compareCmd compares search results with different configurations
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare BM25 scores across different search configurations",
		Long: `Compare how different column weightings affect BM25 search results.

This helps understand the impact of field weighting on relevance ranking:
- Compare default vs custom column weights
- Analyze ranking changes between configurations
- Identify optimal weighting strategies

Examples:
  # Compare default vs title-weighted search
  bm25-fundamentals search compare --query "database" --compare-weights "title:2.0,content:1.0"`,
		RunE: handlers.Search.HandleCompare,
	}

	// explainCmd provides detailed BM25 score explanations
	explainCmd := &cobra.Command{
		Use:   "explain",
		Short: "Explain BM25 scoring mechanics for search results",
		Long: `Provide detailed explanations of how BM25 scores are calculated for search results.

This educational command breaks down:
- Term frequency (TF) contributions from each field
- Inverse document frequency (IDF) calculations
- Document length normalization effects
- Field weighting impact on final scores

Examples:
  # Explain scoring for a specific query
  bm25-fundamentals search explain --query "database optimization"
  
  # Explain with custom field weights
  bm25-fundamentals search explain --query "database" --title-weight 2.0 --content-weight 1.0`,
		RunE: handlers.Search.HandleExplain,
	}

	// setupFlags configures flags for search commands
	setupFlags := func() {
		// Query command flags
		queryCmd.Flags().StringP("query", "q", "", "search query (required)")
		queryCmd.Flags().IntP("max-results", "n", 0, "maximum results to return (0 = use config default)")
		queryCmd.Flags().StringP("category", "c", "", "filter by category")
		queryCmd.Flags().Float64P("title-weight", "", 0, "title field weight (0 = default)")
		queryCmd.Flags().Float64P("content-weight", "", 0, "content field weight (0 = default)")
		queryCmd.Flags().Float64P("category-weight", "", 0, "category field weight (0 = default)")
		queryCmd.Flags().BoolP("snippets", "s", false, "include content snippets")
		queryCmd.Flags().IntP("snippet-length", "", 0, "snippet length in characters (0 = default)")
		queryCmd.MarkFlagRequired("query")

		// Stats command flags
		statsCmd.Flags().StringP("query", "q", "", "search query (required)")
		statsCmd.Flags().StringP("category", "c", "", "filter by category")
		statsCmd.Flags().Float64P("title-weight", "", 0, "title field weight (0 = default)")
		statsCmd.Flags().Float64P("content-weight", "", 0, "content field weight (0 = default)")
		statsCmd.Flags().Float64P("category-weight", "", 0, "category field weight (0 = default)")
		statsCmd.MarkFlagRequired("query")

		// Compare command flags
		compareCmd.Flags().StringP("query", "q", "", "search query (required)")
		compareCmd.Flags().StringP("compare-weights", "", "", "weights to compare (format: field:weight,field:weight)")
		compareCmd.Flags().IntP("max-results", "n", 10, "maximum results for comparison")
		compareCmd.MarkFlagRequired("query")

		// Explain command flags
		explainCmd.Flags().StringP("query", "q", "", "search query (required)")
		explainCmd.Flags().StringP("category", "c", "", "filter by category")
		explainCmd.Flags().Float64P("title-weight", "", 0, "title field weight (0 = default)")
		explainCmd.Flags().Float64P("content-weight", "", 0, "content field weight (0 = default)")
		explainCmd.Flags().Float64P("category-weight", "", 0, "category field weight (0 = default)")
		explainCmd.Flags().IntP("max-results", "n", 5, "maximum results to explain (default: 5)")
		explainCmd.MarkFlagRequired("query")
	}

	// Return the command group
	return &CommandGroup{
		Command: searchCmd,
		SubCommands: []*cobra.Command{
			queryCmd,
			statsCmd,
			compareCmd,
			explainCmd,
		},
		FlagSetup: setupFlags,
	}
}
