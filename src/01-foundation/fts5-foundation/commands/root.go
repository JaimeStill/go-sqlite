package commands

import (
	"fmt"
	"os"

	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/errors"
	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	dbPath  string
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

Phase 1 focuses on establishing the foundation concepts of FTS5.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SQLite FTS5 Foundation Learning Tool")
		fmt.Println("Use 'fts5-foundation --help' to see available commands")

		// Verify FTS5 support
		if err := handlers.VerifyFTS5Support(); err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}
		fmt.Println("✓ SQLite FTS5 support verified")
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

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))
}