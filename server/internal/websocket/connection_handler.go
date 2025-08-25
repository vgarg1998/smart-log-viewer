package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader is the WebSocket upgrader used to convert HTTP connections
// to WebSocket connections. It's configured to accept all origins
// for development purposes.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket handles incoming WebSocket connection requests.
// It upgrades the HTTP connection to a WebSocket connection,
// creates a new Connection instance, and registers it with the hub.
//
// The function includes timeout protection to prevent deadlocks
// during hub registration and starts the necessary goroutines
// for message handling and sending.
//
// Parameters:
//   - w: HTTP response writer for the upgrade response
//   - r: HTTP request containing the upgrade request
//   - hub: The connection hub to register the new connection with
func HandleWebSocket(w http.ResponseWriter, r *http.Request, hub *ConnectionHub) {
	log.Printf("New WebSocket connection request from %s", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket from %s: %v", r.RemoteAddr, err)
		return
	}

	log.Printf("WebSocket upgraded successfully for %s", r.RemoteAddr)
	connection := NewConnection(conn)

	// Non-blocking registration with timeout to prevent deadlock
	select {
	case hub.register <- connection:
		log.Printf("Connection %p queued for registration", connection)
	case <-time.After(5 * time.Second):
		log.Printf("ERROR: Hub registration timeout after 5 seconds, dropping connection %p", connection)
		connection.Close()
		return
	}

	// Start message handler goroutine
	go connection.HandleMessages()
	// Start send goroutine
	go connection.Send()

	log.Printf("WebSocket connection %p started for %s", connection, r.RemoteAddr)
}
