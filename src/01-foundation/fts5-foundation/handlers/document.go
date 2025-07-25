package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/errors"
	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/models"
	"github.com/spf13/viper"
)

// InsertDocument inserts a single document into the FTS5 table
func InsertDocument(title, content, category string) error {
	// Input validation
	if strings.TrimSpace(title) == "" {
		return errors.Validationf("title cannot be empty")
	}
	if strings.TrimSpace(content) == "" {
		return errors.Validationf("content cannot be empty")
	}
	if strings.TrimSpace(category) == "" {
		return errors.Validationf("category cannot be empty")
	}

	// Open database connection
	dbPath := viper.GetString("database")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// Insert the document
	insertSQL := `INSERT INTO documents (title, content, category) VALUES (?, ?, ?)`
	result, err := db.Exec(insertSQL, title, content, category)
	if err != nil {
		return errors.Databasef("failed to insert document: %w", err)
	}

	// Get the row ID for confirmation
	rowID, err := result.LastInsertId()
	if err != nil {
		return errors.Databasef("failed to get last insert ID: %w", err)
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Successfully inserted document with ID: %d\n", rowID)
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Category: %s\n", category)
		fmt.Printf("Content length: %d characters\n", len(content))
	}

	return nil
}

// BatchInsertDocuments inserts multiple documents in a single transaction
func BatchInsertDocuments(documents []models.Document) error {
	if len(documents) == 0 {
		return errors.Validationf("no documents provided for batch insertion")
	}

	// Validate all documents before starting transaction
	for i, doc := range documents {
		if strings.TrimSpace(doc.Title) == "" {
			return errors.Validationf("document %d: title cannot be empty", i+1)
		}
		if strings.TrimSpace(doc.Content) == "" {
			return errors.Validationf("document %d: content cannot be empty", i+1)
		}
		if strings.TrimSpace(doc.Category) == "" {
			return errors.Validationf("document %d: category cannot be empty", i+1)
		}
	}

	// Open database connection
	dbPath := viper.GetString("database")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// Begin transaction for batch operations
	tx, err := db.Begin()
	if err != nil {
		return errors.Transactionf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare statement within transaction
	stmt, err := tx.Prepare("INSERT INTO documents (title, content, category) VALUES (?, ?, ?)")
	if err != nil {
		return errors.Databasef("failed to prepare batch insert statement: %w", err)
	}
	defer stmt.Close()

	// Insert all documents
	insertedCount := 0
	for i, doc := range documents {
		_, err := stmt.Exec(doc.Title, doc.Content, doc.Category)
		if err != nil {
			return errors.Databasef("failed to insert document %d: %w", i+1, err)
		}
		insertedCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Transactionf("failed to commit batch insert transaction: %w", err)
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Successfully inserted %d documents in batch operation\n", insertedCount)

		// Count documents by category for summary
		categoryCount := make(map[string]int)
		for _, doc := range documents {
			categoryCount[doc.Category]++
		}

		fmt.Printf("Documents by category:\n")
		for category, count := range categoryCount {
			fmt.Printf("  %s: %d documents\n", category, count)
		}
	}

	return nil
}

// SearchDocuments performs a full-text search using the MATCH operator
func SearchDocuments(query string, limit int) ([]models.SearchResult, error) {
	// Input validation
	if strings.TrimSpace(query) == "" {
		return nil, errors.Validationf("search query cannot be empty")
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Open database connection
	dbPath := viper.GetString("database")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// Prepare the FTS5 search query with BM25 scoring
	searchSQL := `
		SELECT 
			rowid,
			title,
			content,
			category,
			bm25(documents) as score
		FROM documents 
		WHERE documents MATCH ? 
		ORDER BY score 
		LIMIT ?`

	rows, err := db.Query(searchSQL, query, limit)
	if err != nil {
		return nil, errors.Databasef("failed to execute search query: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var result models.SearchResult
		err := rows.Scan(
			&result.RowID,
			&result.Title,
			&result.Content,
			&result.Category,
			&result.Score,
		)
		if err != nil {
			return nil, errors.Databasef("failed to scan search result: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Databasef("error iterating search results: %w", err)
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Search query: %s\n", query)
		fmt.Printf("Found %d results (limit: %d)\n", len(results), limit)
		if len(results) > 0 {
			fmt.Printf("Best match score: %.4f (lower is better in SQLite FTS5)\n", results[0].Score)
		}
	}

	return results, nil
}

// SearchByCategory performs a category-filtered search
func SearchByCategory(query, category string, limit int) ([]models.SearchResult, error) {
	// Input validation
	if strings.TrimSpace(query) == "" {
		return nil, errors.Validationf("search query cannot be empty")
	}
	if strings.TrimSpace(category) == "" {
		return nil, errors.Validationf("category cannot be empty")
	}

	if limit <= 0 {
		limit = 10
	}

	// Construct FTS5 query with category filter
	ftsQuery := fmt.Sprintf("category:%s AND %s", category, query)

	// Use the main search function with the filtered query
	return SearchDocuments(ftsQuery, limit)
}

// SearchByField performs a field-specific search (title, content, or category)
func SearchByField(query, field string, limit int) ([]models.SearchResult, error) {
	// Input validation
	if strings.TrimSpace(query) == "" {
		return nil, errors.Validationf("search query cannot be empty")
	}

	validFields := []string{"title", "content", "category"}
	fieldValid := false
	for _, validField := range validFields {
		if field == validField {
			fieldValid = true
			break
		}
	}

	if !fieldValid {
		return nil, errors.Validationf("invalid field '%s'. Valid fields: %s", field, strings.Join(validFields, ", "))
	}

	if limit <= 0 {
		limit = 10
	}

	// Construct FTS5 field-specific query
	ftsQuery := fmt.Sprintf("%s:%s", field, query)

	// Use the main search function with the field-specific query
	return SearchDocuments(ftsQuery, limit)
}

// ListDocuments retrieves all documents from the FTS5 table
func ListDocuments(limit int) ([]models.DocumentInfo, error) {
	if limit <= 0 {
		limit = 50 // Default limit for listings
	}

	// Open database connection
	dbPath := viper.GetString("database")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// Query all documents
	listSQL := `
		SELECT 
			rowid,
			title,
			content,
			category
		FROM documents 
		ORDER BY rowid 
		LIMIT ?`

	rows, err := db.Query(listSQL, limit)
	if err != nil {
		return nil, errors.Databasef("failed to list documents: %w", err)
	}
	defer rows.Close()

	var documents []models.DocumentInfo
	for rows.Next() {
		var doc models.DocumentInfo
		var fullContent string

		err := rows.Scan(
			&doc.RowID,
			&doc.Title,
			&fullContent,
			&doc.Category,
		)
		if err != nil {
			return nil, errors.Databasef("failed to scan document: %w", err)
		}

		// Create preview (first 100 characters)
		if len(fullContent) > 100 {
			doc.Preview = fullContent[:100] + "..."
		} else {
			doc.Preview = fullContent
		}

		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Databasef("error iterating documents: %w", err)
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Found %d documents (limit: %d)\n", len(documents), limit)
	}

	return documents, nil
}

// UpdateDocument updates an existing document in the FTS5 table
func UpdateDocument(rowID int64, title, content, category string) error {
	// Input validation
	if rowID <= 0 {
		return errors.Validationf("invalid rowid: %d", rowID)
	}

	// At least one field must be provided for update
	if title == "" && content == "" && category == "" {
		return errors.Validationf("at least one field (title, content, category) must be provided for update")
	}

	// Open database connection
	dbPath := viper.GetString("database")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// First, check if the document exists and get current values
	var currentTitle, currentContent, currentCategory string
	checkSQL := `SELECT title, content, category FROM documents WHERE rowid = ?`
	err = db.QueryRow(checkSQL, rowID).Scan(&currentTitle, &currentContent, &currentCategory)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NotFoundf("document with rowid %d", rowID)
		}
		return errors.Databasef("failed to check existing document: %w", err)
	}

	// Use current values for fields not being updated
	newTitle := currentTitle
	newContent := currentContent
	newCategory := currentCategory

	if strings.TrimSpace(title) != "" {
		newTitle = strings.TrimSpace(title)
	}
	if strings.TrimSpace(content) != "" {
		newContent = strings.TrimSpace(content)
	}
	if strings.TrimSpace(category) != "" {
		newCategory = strings.TrimSpace(category)
	}

	// Update the document (FTS5 will automatically update the index)
	updateSQL := `UPDATE documents SET title = ?, content = ?, category = ? WHERE rowid = ?`
	result, err := db.Exec(updateSQL, newTitle, newContent, newCategory, rowID)
	if err != nil {
		return errors.Databasef("failed to update document: %w", err)
	}

	// Verify the update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Databasef("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NotFoundf("no document was updated (rowid %d may not exist)", rowID)
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Successfully updated document %d\n", rowID)
		fmt.Printf("New values:\n")
		fmt.Printf("  Title: %s\n", newTitle)
		fmt.Printf("  Category: %s\n", newCategory)
		fmt.Printf("  Content length: %d characters\n", len(newContent))
		fmt.Printf("FTS5 index automatically updated\n")
	}

	return nil
}

// DeleteDocument removes a document from the FTS5 table
func DeleteDocument(rowID int64) error {
	// Input validation
	if rowID <= 0 {
		return errors.Validationf("invalid rowid: %d", rowID)
	}

	// Open database connection
	dbPath := viper.GetString("database")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Databasef("failed to open database: %w", err)
	}
	defer db.Close()

	// First, get document info for confirmation (optional but helpful)
	var title, category string
	checkSQL := `SELECT title, category FROM documents WHERE rowid = ?`
	err = db.QueryRow(checkSQL, rowID).Scan(&title, &category)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NotFoundf("document with rowid %d", rowID)
		}
		return errors.Databasef("failed to check existing document: %w", err)
	}

	// Delete the document (FTS5 will automatically update the index)
	deleteSQL := `DELETE FROM documents WHERE rowid = ?`
	result, err := db.Exec(deleteSQL, rowID)
	if err != nil {
		return errors.Databasef("failed to delete document: %w", err)
	}

	// Verify the deletion was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Databasef("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NotFoundf("no document was deleted (rowid %d may not exist)", rowID)
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Successfully deleted document %d\n", rowID)
		fmt.Printf("Deleted document: %s (category: %s)\n", title, category)
		fmt.Printf("FTS5 index automatically updated\n")
	}

	return nil
}

// FormatSearchResults formats search results for display
func FormatSearchResults(results []models.SearchResult, showScores bool) string {
	if len(results) == 0 {
		return "No results found."
	}

	var output strings.Builder

	for i, result := range results {
		output.WriteString(fmt.Sprintf("\n--- Result #%d ---\n", i+1))
		output.WriteString(fmt.Sprintf("Title: %s\n", result.Title))
		output.WriteString(fmt.Sprintf("Category: %s\n", result.Category))

		// Truncate content for display
		content := result.Content
		if len(content) > 150 {
			content = content[:150] + "..."
		}
		output.WriteString(fmt.Sprintf("Content: %s\n", content))

		if showScores {
			output.WriteString(fmt.Sprintf("BM25 Score: %.4f (lower is better)\n", result.Score))
		}
	}

	return output.String()
}

// FormatDocumentList formats document listing for display
func FormatDocumentList(documents []models.DocumentInfo) string {
	if len(documents) == 0 {
		return "No documents found."
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("\nFound %d document(s):\n", len(documents)))

	for _, doc := range documents {
		output.WriteString(fmt.Sprintf("\n--- Document ID: %d ---\n", doc.RowID))
		output.WriteString(fmt.Sprintf("Title: %s\n", doc.Title))
		output.WriteString(fmt.Sprintf("Category: %s\n", doc.Category))
		output.WriteString(fmt.Sprintf("Preview: %s\n", doc.Preview))
	}

	return output.String()
}
