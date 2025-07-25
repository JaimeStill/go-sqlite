package errors

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Sentinel errors for common cases
var (
	// ErrNotFound indicates a requested resource was not found
	ErrNotFound = errors.New("not found")
	
	// ErrValidation indicates input validation failed
	ErrValidation = errors.New("validation failed")
	
	// ErrDatabase indicates a database operation failed
	ErrDatabase = errors.New("database operation failed")
	
	// ErrFTS5 indicates an FTS5-specific operation failed
	ErrFTS5 = errors.New("FTS5 operation failed")
	
	// ErrTransaction indicates a transaction failed
	ErrTransaction = errors.New("transaction failed")
)

// Helper functions for creating contextual errors

// NotFoundf creates a formatted not found error
func NotFoundf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrNotFound}, args...)...)
}

// Validationf creates a formatted validation error
func Validationf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrValidation}, args...)...)
}

// Databasef creates a formatted database error
func Databasef(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrDatabase}, args...)...)
}

// FTS5f creates a formatted FTS5 error
func FTS5f(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrFTS5}, args...)...)
}

// Transactionf creates a formatted transaction error
func Transactionf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrTransaction}, args...)...)
}

// Type checking functions

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsValidation checks if an error is a validation error
func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

// IsDatabase checks if an error is a database error
func IsDatabase(err error) bool {
	return errors.Is(err, ErrDatabase)
}

// IsFTS5 checks if an error is an FTS5 error
func IsFTS5(err error) bool {
	return errors.Is(err, ErrFTS5)
}

// IsTransaction checks if an error is a transaction error
func IsTransaction(err error) bool {
	return errors.Is(err, ErrTransaction)
}

// Display functions for standardized error output

// DisplayError displays an error with appropriate formatting based on its type
// Automatically checks verbose flag and displays full error chain if enabled
func DisplayError(err error) {
	if viper.GetBool("verbose") {
		displayVerbose(err)
	} else {
		displaySimple(err)
	}
}

// displayVerbose displays error with full error chain information
func displayVerbose(err error) {
	displaySimple(err)
	fmt.Printf("\nFull error chain: %+v\n", err)
}

// displaySimple displays error with appropriate formatting based on its type
func displaySimple(err error) {
	if IsValidation(err) {
		fmt.Printf("Validation Error: %v\n", err)
	} else if IsDatabase(err) {
		fmt.Printf("Database Error: %v\n", err)
	} else if IsFTS5(err) {
		fmt.Printf("FTS5 Error: %v\n", err)
		fmt.Println("Hint: Ensure SQLite is compiled with FTS5 support (go build -tags fts5)")
	} else if IsNotFound(err) {
		fmt.Printf("Not Found: %v\n", err)
	} else if IsTransaction(err) {
		fmt.Printf("Transaction Error: %v\n", err)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}