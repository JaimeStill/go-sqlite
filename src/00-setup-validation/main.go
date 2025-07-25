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
