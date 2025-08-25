// Package main is the entry point for the Smart Log Viewer Server.
// This server provides a WebSocket endpoint for real-time log streaming
// and generates mock log entries for demonstration purposes.
package main

import (
	"log"
	"net/http"
	"smart-log-viewer/server/internal/loggenerator"
	"smart-log-viewer/server/internal/model"
	"smart-log-viewer/server/internal/websocket"
	"strconv"
	"time"
)

// main is the entry point for the Smart Log Viewer Server application.
// It initializes the WebSocket connection hub, starts log generation,
// and sets up HTTP endpoints for WebSocket upgrades and server status.
//
// The server runs on port 8080 and provides:
// - WebSocket endpoint at /ws for real-time log streaming
// - Status endpoint at / for server health checks
// - Mock log generation every second for demonstration
//
// The function runs indefinitely until the program is terminated
// or an unrecoverable error occurs.
func main() {
	log.Printf("Starting Smart Log Viewer Server...")

	// Create connection hub
	hub := websocket.NewConnectionHub()

	// Start hub in background
	go hub.Run()

	// Start log generation in background
	go func() {
		count := 0
		for {
			count++
			time.Sleep(1 * time.Second)

			// Create WebSocket message
			message := model.WebSocketMessage{
				Type: "log",
				Data: loggenerator.GenerateMockLog(" - Test message " + strconv.Itoa(count)),
			}

			// Send to broadcast channel
			log.Printf("Sending log #%d to broadcast channel...", count)
			hub.Broadcast <- message
			log.Printf("Sent log #%d to broadcast channel", count)
		}
	}()

	// HTTP handler for WebSocket upgrade
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(w, r, hub)
	})

	// Simple status endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Smart Log Viewer Server is running!\nConnect to /ws for WebSocket"))
	})

	log.Printf("Server starting on :8080")
	log.Printf("WebSocket endpoint: ws://localhost:8080/ws")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
