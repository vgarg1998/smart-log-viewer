// Package loggenerator provides functionality for creating mock log entries.
// This package is used for testing and demonstration purposes to simulate
// real-world log generation scenarios.
package loggenerator

import (
	"math/rand"
	"time"

	"smart-log-viewer/server/internal/model"
)

// levels contains the available log severity levels for mock log generation.
// These levels follow standard logging conventions and are randomly
// selected when creating mock log entries.
var levels = []string{"INFO", "WARN", "ERROR"}

// GenerateMockLog creates a mock log entry with a random log level (INFO, WARN, or ERROR),
// the provided message, and current timestamp. This function is used for testing and
// demonstration purposes to simulate log generation.
//
// The function generates realistic log entries by randomly selecting
// severity levels and combining them with the provided message text.
//
// Parameters:
//   - message: The message text to append to the mock log entry
//
// Returns:
//   - model.Log: A mock log entry with random level, message, and current timestamp
func GenerateMockLog(message string) model.Log {
	return model.Log{
		Level:     levels[rand.Intn(len(levels))],
		Message:   "This is a mock log message" + message,
		Timestamp: time.Now(),
	}
}
