package errors

import (
	"fmt"
	"log"
	"os"
)

// CustomError wraps the original error with a message
type CustomError struct {
	Message string
	Err     error
}

// New creates a new CustomError
func New(message string, err error) error {
	if err == nil {
		return nil
	}
	return &CustomError{
		Message: message,
		Err:     err,
	}
}

// Error returns the error message
func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// LogError logs the error message to the console or a file
func LogError(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

// FatalError logs the error and exits the program
func FatalError(err error) {
	if err != nil {
		log.Println("Fatal error:", err)
		os.Exit(1)
	}
}
