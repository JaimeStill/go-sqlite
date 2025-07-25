package commands

import (
	"fmt"
	"os"

	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/errors"
	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/handlers"
	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/models"
	"github.com/spf13/cobra"
)

// documentCmd stores the document command for flag registration
var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Document table operations",
	Long: `Commands for managing documents in the FTS5 table.

This includes creating the table, inserting documents, searching,
and performing CRUD operations on the document collection.`,
}

// documentGroup represents the document command group with all sub-commands
var documentGroup = &CommandGroup{
	Command: documentCmd,
	SubCommands: []*cobra.Command{
		createTableCmd,
		insertCmd,
		batchInsertCmd,
		searchCmd,
		searchCategoryCmd,
		searchFieldCmd,
		listCmd,
		updateCmd,
		deleteCmd,
	},
	FlagSetup: setupDocumentFlags,
}

// Document sub-commands

// createTableCmd represents the create-table command
var createTableCmd = &cobra.Command{
	Use:   "create-table",
	Short: "Create the FTS5 virtual table for document storage",
	Long: `Creates an FTS5 virtual table named 'documents' with columns for title, content, and category.

This demonstrates the basic FTS5 virtual table creation syntax and sets up
the table structure for document insertion and search operations.

The table uses the unicode61 tokenizer with diacritics removal for better
international text handling.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handlers.CreateDocumentsTable(); err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}
		fmt.Println("✓ FTS5 documents table created successfully")
	},
}

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a single document into the FTS5 table",
	Long: `Insert a single document with title, content, and category into the FTS5 table.

This demonstrates basic FTS5 document insertion patterns and automatic indexing.

Example usage:
  fts5-foundation document insert --title "My Document" --content "Document content here" --category "example"`,
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		category, _ := cmd.Flags().GetString("category")

		// Validate required fields
		if title == "" || content == "" || category == "" {
			fmt.Println("Error: title, content, and category are all required")
			fmt.Println("Usage:")
			fmt.Println("  --title \"Document Title\"")
			fmt.Println("  --content \"Document content...\"")
			fmt.Println("  --category \"document-category\"")
			os.Exit(1)
		}

		// Insert the document
		if err := handlers.InsertDocument(title, content, category); err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		fmt.Println("✓ Document inserted successfully")
	},
}

// batchInsertCmd represents the batch-insert command
var batchInsertCmd = &cobra.Command{
	Use:   "batch-insert",
	Short: "Insert multiple example documents for learning",
	Long: `Insert multiple predefined documents in a single transaction for better performance.

This demonstrates batch insertion patterns and transaction handling in FTS5.
The example documents cover various topics to showcase FTS5 search capabilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create sample documents for learning purposes
		documents := []models.Document{
			{
				Title:    "Introduction to Go Programming",
				Content:  "Go is a statically typed, compiled programming language designed at Google. It combines the efficiency of a compiled language with the ease of programming of an interpreted language.",
				Category: "programming",
			},
			{
				Title:    "SQLite FTS5 Full-Text Search",
				Content:  "SQLite FTS5 is an SQLite virtual table module that provides full-text search functionality. It supports advanced features like BM25 ranking, custom tokenizers, and phrase queries.",
				Category: "database",
			},
			{
				Title:    "Understanding BM25 Scoring Algorithm",
				Content:  "BM25 is a ranking function used by search engines to estimate the relevance of documents to a given search query. It considers term frequency, inverse document frequency, and document length.",
				Category: "algorithms",
			},
			{
				Title:    "Database Indexing Fundamentals",
				Content:  "Database indexes are data structures that improve the speed of data retrieval operations. They work by creating shortcuts to the data, trading storage space for faster reads.",
				Category: "database",
			},
		}

		fmt.Printf("Inserting %d example documents...\n", len(documents))

		// Perform batch insertion
		if err := handlers.BatchInsertDocuments(documents); err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		fmt.Printf("✓ Successfully inserted %d documents\n", len(documents))
	},
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search documents using FTS5 MATCH operator",
	Long: `Search documents in the FTS5 table using full-text search with BM25 scoring.

This demonstrates basic FTS5 MATCH queries and BM25 relevance ranking.
SQLite FTS5 returns negative BM25 scores where lower values indicate better matches.

Example usage:
  fts5-foundation document search "golang programming"
  fts5-foundation document search "database" --limit 5
  fts5-foundation document search "sqlite" --scores`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		limit, _ := cmd.Flags().GetInt("limit")
		showScores, _ := cmd.Flags().GetBool("scores")

		// Perform the search
		results, err := handlers.SearchDocuments(query, limit)
		if err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		// Display results
		if len(results) == 0 {
			fmt.Printf("No documents found matching: %s\n", query)
			return
		}

		fmt.Printf("Found %d document(s) matching: %s\n", len(results), query)
		fmt.Println(handlers.FormatSearchResults(results, showScores))
	},
}

// searchCategoryCmd represents the search-category command
var searchCategoryCmd = &cobra.Command{
	Use:   "search-category [query] [category]",
	Short: "Search documents within a specific category",
	Long: `Search for documents within a specific category using FTS5 field filtering.

This demonstrates FTS5 field-specific queries and category filtering patterns.

Example usage:
  fts5-foundation document search-category "programming" "database"
  fts5-foundation document search-category "algorithm" "programming" --limit 3`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		category := args[1]
		limit, _ := cmd.Flags().GetInt("limit")
		showScores, _ := cmd.Flags().GetBool("scores")

		// Perform category-filtered search
		results, err := handlers.SearchByCategory(query, category, limit)
		if err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		// Display results
		if len(results) == 0 {
			fmt.Printf("No documents found matching '%s' in category '%s'\n", query, category)
			return
		}

		fmt.Printf("Found %d document(s) matching '%s' in category '%s':\n", len(results), query, category)
		fmt.Println(handlers.FormatSearchResults(results, showScores))
	},
}

// searchFieldCmd represents the search-field command
var searchFieldCmd = &cobra.Command{
	Use:   "search-field [query] [field]",
	Short: "Search within a specific field (title, content, or category)",
	Long: `Search for documents within a specific field using FTS5 field-specific queries.

This demonstrates FTS5 column-specific search patterns.

Valid fields: title, content, category

Example usage:
  fts5-foundation document search-field "Introduction" "title"
  fts5-foundation document search-field "algorithm" "content"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		field := args[1]
		limit, _ := cmd.Flags().GetInt("limit")
		showScores, _ := cmd.Flags().GetBool("scores")

		// Perform field-specific search
		results, err := handlers.SearchByField(query, field, limit)
		if err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		// Display results
		if len(results) == 0 {
			fmt.Printf("No documents found matching '%s' in field '%s'\n", query, field)
			return
		}

		fmt.Printf("Found %d document(s) matching '%s' in field '%s':\n", len(results), query, field)
		fmt.Println(handlers.FormatSearchResults(results, showScores))
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all documents in the FTS5 table",
	Long: `List all documents with their IDs, titles, categories, and content previews.

This is useful for seeing what documents exist and finding their row IDs
for update and delete operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")

		// List documents
		documents, err := handlers.ListDocuments(limit)
		if err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		// Display results
		if len(documents) == 0 {
			fmt.Println("No documents found in the database")
			return
		}

		fmt.Println(handlers.FormatDocumentList(documents))
	},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [rowid]",
	Short: "Update an existing document",
	Long: `Update an existing document by row ID. You can update the title, content, 
and/or category. Only the fields you specify will be changed.

The FTS5 index will be automatically updated when the document is modified.

Example usage:
  fts5-foundation document update 1 --title "New Title"
  fts5-foundation document update 2 --content "New content here" --category "updated"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse row ID
		var rowID int64
		if _, err := fmt.Sscanf(args[0], "%d", &rowID); err != nil {
			fmt.Printf("Invalid row ID '%s': must be a number\n", args[0])
			os.Exit(1)
		}

		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		category, _ := cmd.Flags().GetString("category")

		// Update the document
		if err := handlers.UpdateDocument(rowID, title, content, category); err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		fmt.Printf("✓ Document %d updated successfully\n", rowID)
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [rowid]",
	Short: "Delete a document from the FTS5 table",
	Long: `Delete a document by its row ID. The FTS5 index will be automatically 
updated when the document is removed.

Use the 'list' command to find document row IDs.

Example usage:
  fts5-foundation document delete 1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse row ID
		var rowID int64
		if _, err := fmt.Sscanf(args[0], "%d", &rowID); err != nil {
			fmt.Printf("Invalid row ID '%s': must be a number\n", args[0])
			os.Exit(1)
		}

		// Delete the document
		if err := handlers.DeleteDocument(rowID); err != nil {
			errors.DisplayError(err)
			os.Exit(1)
		}

		fmt.Printf("✓ Document %d deleted successfully\n", rowID)
	},
}

// Flag setup function

func setupDocumentFlags() {
	// Insert command flags
	insertCmd.Flags().StringP("title", "t", "", "Document title")
	insertCmd.Flags().StringP("content", "c", "", "Document content")
	insertCmd.Flags().StringP("category", "g", "", "Document category")

	// Search command flags
	searchCmd.Flags().IntP("limit", "l", 10, "Maximum number of results to return")
	searchCmd.Flags().BoolP("scores", "s", false, "Show BM25 scores and ranking information")

	// Search-category command flags
	searchCategoryCmd.Flags().IntP("limit", "l", 10, "Maximum number of results to return")
	searchCategoryCmd.Flags().BoolP("scores", "s", false, "Show BM25 scores and ranking information")

	// Search-field command flags
	searchFieldCmd.Flags().IntP("limit", "l", 10, "Maximum number of results to return")
	searchFieldCmd.Flags().BoolP("scores", "s", false, "Show BM25 scores and ranking information")

	// List command flags
	listCmd.Flags().IntP("limit", "l", 50, "Maximum number of documents to list")

	// Update command flags
	updateCmd.Flags().StringP("title", "t", "", "New document title")
	updateCmd.Flags().StringP("content", "c", "", "New document content")
	updateCmd.Flags().StringP("category", "g", "", "New document category")
}
