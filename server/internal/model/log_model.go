package model

import "time"

// Log represents a log entry with a level, message, and timestamp.
// This struct is used to represent log messages that are generated
// by the server and sent to connected WebSocket clients.
//
// The struct includes JSON tags for serialization and follows
// standard logging conventions with severity levels.
type Log struct {
	// Level represents the severity of the log entry.
	// Common values include "INFO", "WARN", and "ERROR".
	Level string `json:"level"`

	// Message contains the actual log message text.
	// This field holds the descriptive information about the log event.
	Message string `json:"message"`

	// Timestamp records when the log entry was created.
	// This field uses Go's time.Time type for precise timestamp handling.
	Timestamp time.Time `json:"timestamp"`
}
