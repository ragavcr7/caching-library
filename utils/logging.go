// this function is used for programmers to understand the overflow of the project by logging at certain points
package utils

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

func init() {
	// Initialize logger
	logger = log.New(os.Stdout, "[MyApp] ", log.Ldate|log.Ltime)
}

// Info logs an informational message
func Info(message string) {
	logger.Printf("[INFO] %s\n", message)
}

// Warn logs a warning message
func Warn(message string) {
	logger.Printf("[WARN] %s\n", message)
}

// Error logs an error message
func Error(message string) {
	logger.Printf("[ERROR] %s\n", message)
}

// Debug logs a debug message
func Debug(message string) {
	logger.Printf("[DEBUG] %s\n", message)
}
