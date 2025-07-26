package commands

import (
	"fmt"
	"os"

	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/config"
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	format  string
)

// rootCmd stores the root command for flag registration
var rootCmd = &cobra.Command{
	Use:   "setup-validation",
	Short: "Phase 0: SQLite FTS5 Setup Validation Tool",
	Long: `A CLI tool for validating that the Go-SQLite FTS5 learning environment 
is properly configured and all utilities are working correctly.

This validation tool ensures:
- SQLite connection and FTS5 support
- Sample data generation and insertion
- BM25 scoring functionality
- Utility function accessibility

Phase 0 establishes a solid foundation for subsequent learning phases.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration after flags are parsed
		config.App = &config.Config{}
		config.App.Init()
		
		// Initialize database connection (in-memory for validation)
		if err := database.Init(":memory:"); err != nil {
			fmt.Fprintf(os.Stderr, "Database initialization error: %v\n", err)
			os.Exit(1)
		}
		
		// Handlers are stateless - no initialization needed
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SQLite FTS5 Setup Validation Tool")
		fmt.Println("Use 'setup-validation --help' to see available commands")
		fmt.Println()
		fmt.Println("Quick start:")
		fmt.Println("  setup-validation validate    # Run all validation checks")
		fmt.Println("  setup-validation connect     # Test SQLite connection")
		fmt.Println("  setup-validation fts5        # Test FTS5 functionality")
		fmt.Println("  setup-validation bm25        # Test BM25 scoring")
	},
}

// Root represents the root command group with all child groups initialized
var Root = &CommandGroup{
	Command:     rootCmd,
	ChildGroups: []*CommandGroup{
		Validation,
	},
	FlagSetup: setupGlobalFlags,
}

// setupGlobalFlags registers global flags that are available to all commands
func setupGlobalFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.setup-validation.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output with detailed explanations")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "output format (text, json)")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
}