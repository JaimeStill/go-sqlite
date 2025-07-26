package errors

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Sentinel errors for type checking
var (
	ErrValidation   = errors.New("validation failed")
	ErrDatabase     = errors.New("database operation failed")
	ErrFTS5         = errors.New("FTS5 operation failed")
	ErrNotFound     = errors.New("not found")
	ErrTransaction  = errors.New("transaction failed")
	ErrAnalysis     = errors.New("analysis failed")
	ErrVisualization = errors.New("visualization failed")
)

// Error creation helpers
func Validationf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrValidation}, args...)...)
}

func Databasef(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrDatabase}, args...)...)
}

func FTS5f(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrFTS5}, args...)...)
}

func NotFoundf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrNotFound}, args...)...)
}

func Transactionf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrTransaction}, args...)...)
}

func Analysisf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrAnalysis}, args...)...)
}

func Visualizationf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrVisualization}, args...)...)
}

// DisplayError shows the error in an appropriate format based on verbose flag
func DisplayError(err error) {
	if err == nil {
		return
	}

	if viper.GetBool("verbose") {
		displayVerboseError(err)
	} else {
		displaySimpleError(err)
	}
}

// displaySimpleError shows a user-friendly error message
func displaySimpleError(err error) {
	switch {
	case errors.Is(err, ErrValidation):
		fmt.Fprintf(os.Stderr, "Validation Error: %v\n", unwrapError(err))
	case errors.Is(err, ErrDatabase):
		fmt.Fprintf(os.Stderr, "Database Error: %v\n", unwrapError(err))
		fmt.Fprintln(os.Stderr, "Hint: Check that the database file exists and is accessible")
	case errors.Is(err, ErrFTS5):
		fmt.Fprintf(os.Stderr, "FTS5 Error: %v\n", unwrapError(err))
		fmt.Fprintln(os.Stderr, "Hint: Ensure SQLite is compiled with FTS5 support (go build -tags fts5)")
	case errors.Is(err, ErrNotFound):
		fmt.Fprintf(os.Stderr, "Not Found: %v\n", unwrapError(err))
	case errors.Is(err, ErrTransaction):
		fmt.Fprintf(os.Stderr, "Transaction Error: %v\n", unwrapError(err))
		fmt.Fprintln(os.Stderr, "Hint: The operation was rolled back; no changes were made")
	case errors.Is(err, ErrAnalysis):
		fmt.Fprintf(os.Stderr, "Analysis Error: %v\n", unwrapError(err))
		fmt.Fprintln(os.Stderr, "Hint: Check that you have search results to analyze")
	case errors.Is(err, ErrVisualization):
		fmt.Fprintf(os.Stderr, "Visualization Error: %v\n", unwrapError(err))
		fmt.Fprintln(os.Stderr, "Hint: Check terminal width and visualization settings")
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

// displayVerboseError shows the full error chain for debugging
func displayVerboseError(err error) {
	fmt.Fprintln(os.Stderr, "=== VERBOSE ERROR OUTPUT ===")
	fmt.Fprintf(os.Stderr, "Error Type: %T\n", err)
	
	// Display full error chain
	fmt.Fprintln(os.Stderr, "\nError Chain:")
	currentErr := err
	depth := 0
	for currentErr != nil {
		indent := strings.Repeat("  ", depth)
		fmt.Fprintf(os.Stderr, "%s└─ %v\n", indent, currentErr)
		currentErr = errors.Unwrap(currentErr)
		depth++
	}
	
	// Display categorized error type
	fmt.Fprintln(os.Stderr, "\nError Category:")
	switch {
	case errors.Is(err, ErrValidation):
		fmt.Fprintln(os.Stderr, "  Validation Error - Input validation failed")
	case errors.Is(err, ErrDatabase):
		fmt.Fprintln(os.Stderr, "  Database Error - SQLite operation failed")
	case errors.Is(err, ErrFTS5):
		fmt.Fprintln(os.Stderr, "  FTS5 Error - Full-text search operation failed")
	case errors.Is(err, ErrNotFound):
		fmt.Fprintln(os.Stderr, "  Not Found Error - Requested resource does not exist")
	case errors.Is(err, ErrTransaction):
		fmt.Fprintln(os.Stderr, "  Transaction Error - Database transaction failed")
	case errors.Is(err, ErrAnalysis):
		fmt.Fprintln(os.Stderr, "  Analysis Error - Score analysis operation failed")
	case errors.Is(err, ErrVisualization):
		fmt.Fprintln(os.Stderr, "  Visualization Error - Chart rendering failed")
	default:
		fmt.Fprintln(os.Stderr, "  Uncategorized Error")
	}
	
	fmt.Fprintln(os.Stderr, "\n=== END VERBOSE OUTPUT ===")
}

// unwrapError gets the deepest error message without the wrapping context
func unwrapError(err error) string {
	if err == nil {
		return ""
	}
	
	// Get the wrapped error message by removing the prefix
	parts := strings.SplitN(err.Error(), ": ", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return err.Error()
}