package shared

import (
	"fmt"
	"os"
)

// HandleError provides consistent error handling for CLI applications
func HandleError(err error, context string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s: %v\n", context, err)
		os.Exit(1)
	}
}

// CheckError is a simple error check that panics on error (for learning examples)
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}