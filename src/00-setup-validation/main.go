package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "setup-validation",
	Short: "Validate SQLite FTS5 and shared utilities setup",
	Long: `A CLI tool to validate that the Go-SQLite FTS5 learning environment 
is properly configured and all shared utilities are working correctly.`,
}

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

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(fts5Cmd)
	rootCmd.AddCommand(testDataCmd)
	rootCmd.AddCommand(bm25Cmd)
	
	// Global flags
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}