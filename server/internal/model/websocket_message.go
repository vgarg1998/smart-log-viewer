package model

// WebSocketMessage represents a message sent over the WebSocket connection.
// This struct provides a standardized format for all WebSocket communication
// between the server and clients, including logs, control messages, and
// heartbeat responses.
//
// The message type system allows for different kinds of data to be
// transmitted while maintaining a consistent structure.
type WebSocketMessage struct {
	// Type identifies the kind of message being sent.
	// Common types include "log", "pause", "resume", "ping", and "pong".
	Type string `json:"type"`

	// Data contains the actual message payload.
	// The type of data varies based on the message type.
	// For log messages, this will be a Log struct.
	// For control messages, this may be null or a simple string.
	Data interface{} `json:"data"`
}
