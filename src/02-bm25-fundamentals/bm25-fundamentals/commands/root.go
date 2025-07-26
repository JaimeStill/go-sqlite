package commands

import (
	"fmt"
	"os"

	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/config"
	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	dbPath  string
	format  string
)

// rootCmd stores the root command for flag registration
var rootCmd = &cobra.Command{
	Use:   "bm25-fundamentals",
	Short: "Phase 2: BM25 Fundamentals Learning Tool",
	Long: `A CLI tool for mastering BM25 scoring mechanics and interpretation in SQLite FTS5.

This educational tool demonstrates:
- BM25 score analysis and distribution
- Document length normalization effects
- Column weighting experiments
- Ranking comparison and optimization

Phase 2 builds on FTS5 foundation to explore relevance scoring in depth.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration after flags are parsed
		config.App.Init()
		
		// Initialize database connection
		if err := database.Init(config.App.GetDatabasePath()); err != nil {
			fmt.Fprintf(os.Stderr, "Database initialization error: %v\n", err)
			os.Exit(1)
		}
		
		// Handlers are stateless - no initialization needed
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BM25 Fundamentals Learning Tool")
		fmt.Println("Use 'bm25-fundamentals --help' to see available commands")
		fmt.Println()
		fmt.Println("Key learning areas:")
		fmt.Println("  • Understanding negative BM25 scores (lower = better)")
		fmt.Println("  • Impact of document length on scoring")
		fmt.Println("  • Multi-field relevance tuning")
		fmt.Println("  • Score distribution analysis")
	},
}

// Root represents the root command group with all child groups initialized
var Root = &CommandGroup{
	Command:     rootCmd,
	ChildGroups: []*CommandGroup{
		Corpus,
		Search,
		Visualize,
	},
	FlagSetup: setupGlobalFlags,
}

// setupGlobalFlags registers global flags that are available to all commands
func setupGlobalFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bm25-fundamentals.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output with detailed explanations")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "database", "d", ":memory:", "database path (default: in-memory)")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "output format (text, json, csv)")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
}
