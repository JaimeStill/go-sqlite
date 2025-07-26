package errors

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Sentinel errors for type-safe error checking
var (
	ErrValidation = errors.New("validation failed")
	ErrDatabase   = errors.New("database operation failed")
	ErrFTS5       = errors.New("FTS5 operation failed")
	ErrConnection = errors.New("connection failed")
)

// Helper functions for creating typed errors

func Validationf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrValidation}, args...)...)
}

func Databasef(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrDatabase}, args...)...)
}

func FTS5f(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrFTS5}, args...)...)
}

func Connectionf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{ErrConnection}, args...)...)
}

// DisplayError provides consistent error display with automatic verbose handling
func DisplayError(err error) {
	if viper.GetBool("verbose") {
		displayVerbose(err)
	} else {
		displaySimple(err)
	}
	os.Exit(1)
}

// displaySimple shows user-friendly error messages
func displaySimple(err error) {
	switch {
	case errors.Is(err, ErrValidation):
		fmt.Fprintf(os.Stderr, "❌ Validation Error: %v\n", extractMessage(err))
		fmt.Fprintf(os.Stderr, "💡 Tip: Use --verbose for detailed error information\n")
	case errors.Is(err, ErrDatabase):
		fmt.Fprintf(os.Stderr, "❌ Database Error: %v\n", extractMessage(err))
		fmt.Fprintf(os.Stderr, "💡 Tip: Check that SQLite is properly installed\n")
	case errors.Is(err, ErrFTS5):
		fmt.Fprintf(os.Stderr, "❌ FTS5 Error: %v\n", extractMessage(err))
		fmt.Fprintf(os.Stderr, "💡 Tip: Ensure SQLite was compiled with FTS5 support\n")
	case errors.Is(err, ErrConnection):
		fmt.Fprintf(os.Stderr, "❌ Connection Error: %v\n", extractMessage(err))
		fmt.Fprintf(os.Stderr, "💡 Tip: Verify database path and permissions\n")
	default:
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
	}
}

// displayVerbose shows detailed technical error information
func displayVerbose(err error) {
	fmt.Fprintf(os.Stderr, "❌ Detailed Error Information:\n")
	fmt.Fprintf(os.Stderr, "   Error: %v\n", err)
	fmt.Fprintf(os.Stderr, "   Type: %T\n", err)
	
	// Show error chain
	current := err
	depth := 0
	for current != nil {
		if depth > 0 {
			fmt.Fprintf(os.Stderr, "   Caused by [%d]: %v\n", depth, current)
		}
		if unwrapped := errors.Unwrap(current); unwrapped != nil {
			current = unwrapped
			depth++
		} else {
			break
		}
	}
}

// extractMessage extracts the user-friendly part of an error message
func extractMessage(err error) string {
	msg := err.Error()
	
	// Remove the sentinel error prefix for cleaner display
	if idx := findLastColon(msg); idx != -1 {
		return msg[idx+2:] // +2 to skip ": "
	}
	
	return msg
}

// findLastColon finds the last occurrence of ": " in a string
func findLastColon(s string) int {
	for i := len(s) - 2; i >= 0; i-- {
		if s[i] == ':' && s[i+1] == ' ' {
			return i
		}
	}
	return -1
}