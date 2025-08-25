package websocket

import (
	"log"
	"smart-log-viewer/server/internal/model"
	"time"
)

// ConnectionHub manages all active WebSocket connections.
// It provides centralized connection management including registration,
// unregistration, broadcasting, and health monitoring.
//
// The hub runs in a single goroutine to avoid race conditions and
// coordinates all connection operations through channels.
type ConnectionHub struct {
	connections map[*Connection]bool
	register    chan *Connection
	unregister  chan *Connection
	Broadcast   chan model.WebSocketMessage // Capitalized to make it public
}

// NewConnectionHub creates a new connection hub instance.
// It initializes the hub with empty connection maps and
// buffered channels for connection management.
//
// Returns:
//   - *ConnectionHub: A new connection hub instance
func NewConnectionHub() *ConnectionHub {
	log.Printf("Creating new ConnectionHub")
	return &ConnectionHub{
		connections: make(map[*Connection]bool),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		Broadcast:   make(chan model.WebSocketMessage),
	}
}

// checkConnectionHealth performs health checks on all active connections.
// It identifies connections that should be dropped due to performance
// issues or lack of responsiveness and queues them for unregistration.
//
// This method runs periodically and processes connections in batches
// to avoid blocking the main hub loop.
func (h *ConnectionHub) checkConnectionHealth() {
	// Check each connection for health issues
	connectionsToDrop := make([]*Connection, 0)

	for connection := range h.connections {
		// Immediately remove closed connections
		if connection.IsClosed() {
			log.Printf("Health check: Connection %p is closed, removing immediately", connection)
			delete(h.connections, connection)
			connection.Close()
			continue
		}

		// Check if connection should be dropped
		if connection.shouldDrop() {
			log.Printf("Health check: Connection %p should be dropped, queuing for unregister", connection)
			connectionsToDrop = append(connectionsToDrop, connection)
			continue
		}

		// For paused connections, send ping to check if they're still alive
		if connection.IsPaused() {
			if err := connection.SendPing(); err != nil {
				log.Printf("Health check: Failed to ping paused connection %p, queuing for unregister", connection)
				connectionsToDrop = append(connectionsToDrop, connection)
				continue
			}
		}
	}

	// Process all connections to drop in batch (non-blocking)
	for _, conn := range connectionsToDrop {
		select {
		case h.unregister <- conn:
			log.Printf("Health check: Connection %p queued for unregister", conn)
		default:
			log.Printf("WARNING: Health check unregister channel full, connection %p dropped immediately", conn)
			conn.Close()
		}
	}
}

// Run is the main event loop for the connection hub.
// It handles all connection lifecycle events including registration,
// unregistration, health checks, and message broadcasting.
//
// This method runs in a single goroutine and coordinates all
// connection operations to prevent race conditions.
//
// The hub will continue running until the program exits or an
// unrecoverable error occurs.
func (h *ConnectionHub) Run() {
	log.Printf("Starting ConnectionHub main loop")

	// Start health check ticker
	healthTicker := time.NewTicker(2 * time.Second)
	defer healthTicker.Stop()

	for {
		select {
		case connection := <-h.register:
			h.connections[connection] = true
			log.Printf("REGISTERED Connection %p, total connections: %d", connection, len(h.connections))

		case connection := <-h.unregister:
			delete(h.connections, connection)
			log.Printf("UNREGISTERED Connection %p, total connections: %d", connection, len(h.connections))
			connection.Close()

		case <-healthTicker.C:
			// Check connection health every 2 seconds
			h.checkConnectionHealth()

		case logEntry := <-h.Broadcast:
			if len(h.connections) == 0 {
				continue
			}

			log.Printf("Broadcasting message type '%s' to %d connections", logEntry.Type, len(h.connections))

			// Filter out closed connections and broadcast to active ones
			activeConnections := make([]*Connection, 0)
			for conn := range h.connections {
				if conn.IsClosed() {
					// Immediately unregister closed connections
					select {
					case h.unregister <- conn:
						log.Printf("Broadcast: Closed connection %p queued for unregister", conn)
					default:
						log.Printf("Broadcast: Unregister channel full, dropping closed connection %p immediately", conn)
						delete(h.connections, conn)
						conn.Close()
					}
					continue
				}
				activeConnections = append(activeConnections, conn)
			}

			// Broadcast to active connections only
			for _, conn := range activeConnections {
				go func(c *Connection) {
					if !c.IsPaused() && !c.shouldDrop() {
						c.sendLog(logEntry)
					}
				}(conn)
			}
		}
	}
}
