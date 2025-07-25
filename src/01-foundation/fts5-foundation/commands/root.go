package commands

import (
	"fmt"
	"os"

	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/config"
	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/database"
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
	Use:   "fts5-foundation",
	Short: "Phase 1: SQLite FTS5 Foundation Learning Tool",
	Long: `A CLI tool for learning SQLite FTS5 fundamentals and BM25 scoring.

This educational tool demonstrates:
- FTS5 virtual table creation
- Basic document insertion and indexing
- Simple search queries with BM25 relevance scoring
- CRUD operations with automatic indexing

Phase 1 focuses on establishing the foundation concepts of FTS5.`,
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
		fmt.Println("SQLite FTS5 Foundation Learning Tool")
		fmt.Println("Use 'fts5-foundation --help' to see available commands")
		fmt.Println()
		fmt.Println("Key learning areas:")
		fmt.Println("  • FTS5 virtual table creation and management")
		fmt.Println("  • Document insertion with automatic indexing")
		fmt.Println("  • Basic search operations with MATCH operator")
		fmt.Println("  • BM25 scoring and relevance ranking")
	},
}

// Root represents the root command group with all child groups initialized
var Root = &CommandGroup{
	Command:     rootCmd,
	ChildGroups: []*CommandGroup{documentGroup},
	FlagSetup:   setupGlobalFlags,
}

// setupGlobalFlags registers global flags that are available to all commands
func setupGlobalFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fts5-foundation.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "database", "d", ":memory:", "database path (default: in-memory)")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "output format (text, json)")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
}