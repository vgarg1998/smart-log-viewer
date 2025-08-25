package websocket

import (
	"fmt"
	"log"
	"smart-log-viewer/server/internal/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection represents a WebSocket connection with a channel for log messages.
// It manages the lifecycle of a single client connection including message
// buffering, pause state, and health monitoring.
type Connection struct {
	ws       *websocket.Conn
	channel  chan model.WebSocketMessage
	lastSent time.Time    // Each connection tracks its own timing
	mu       sync.RWMutex // Protect connection's own state (read/write mutex)
	isClosed bool
	isPaused bool // Track if client is paused
}

// NewConnection creates a new WebSocket connection instance.
// It initializes the connection with default values and creates
// a buffered channel for message handling.
//
// Parameters:
//   - ws: The underlying WebSocket connection
//
// Returns:
//   - *Connection: A new connection instance
func NewConnection(ws *websocket.Conn) *Connection {
	log.Printf("Creating new WebSocket connection: %p", ws)
	return &Connection{
		ws:       ws,
		channel:  make(chan model.WebSocketMessage, 100), // Buffer for better performance
		lastSent: time.Now(),
		isClosed: false,
		isPaused: false, // Start as not paused
	}
}

// Send runs the main send goroutine for this connection.
// It continuously reads messages from the channel and sends them
// to the WebSocket client. This method runs in a separate goroutine
// and handles the complete lifecycle of message sending.
//
// The goroutine will exit when the channel is closed or an error occurs.
func (c *Connection) Send() {
	log.Printf("STARTING Send goroutine for connection: %p", c)
	defer func() {
		log.Printf("SEND GOROUTINE EXITING for connection: %p", c)
		c.ws.Close()
		c.Close()
	}()

	for message := range c.channel {
		if err := c.ws.WriteJSON(message); err != nil {
			log.Printf("Error sending message to client %p: %v", c, err)
			break
		}
		log.Printf("Successfully sent message type '%s' to connection %p", message.Type, c)
	}
	log.Printf("Send goroutine finished for connection: %p", c)
}

// Close safely closes the connection and cleans up resources.
// It marks the connection as closed, clears missed logs to prevent
// memory leaks, and safely closes the message channel.
//
// This method is thread-safe and can be called multiple times safely.
func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		log.Printf("Connection %p already closed, skipping", c)
		return
	}

	log.Printf("CLOSING CONNECTION %p", c)
	c.isClosed = true // Mark as closed

	// Close channel safely
	select {
	case <-c.channel:
		// Channel already closed, do nothing
	default:
		close(c.channel) // Close channel
	}
}

// IsClosed safely checks if the connection is closed.
// This method is thread-safe and provides a safe way to check
// the connection state without race conditions.
//
// Returns:
//   - bool: true if the connection is closed, false otherwise
func (c *Connection) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isClosed
}

// sendLog sends a log message to the client.
// It attempts to send the message and logs the result.
// If the channel is full, the message is skipped (TCP will handle backpressure).
//
// Parameters:
//   - message: The WebSocket message to send
func (c *Connection) sendLog(message model.WebSocketMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		log.Printf("Connection %p is closed, cannot send message type '%s'", c, message.Type)
		return
	}

	select {
	case c.channel <- message:
		c.lastSent = time.Now()
		log.Printf("Successfully sent message type '%s' to connection %p", message.Type, c)
	default:
		log.Printf("Connection %p channel full, skipping message type '%s' (TCP will handle backpressure)", c, message.Type)
	}
}

// shouldDrop determines if this connection should be dropped.
// It checks if the connection is closed or if paused connections
// haven't responded to ping messages within the timeout period.
//
// Returns:
//   - bool: true if the connection should be dropped, false otherwise
func (c *Connection) shouldDrop() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isClosed {
		log.Printf("Connection %p is closed, should be dropped", c)
		return true
	}

	// For paused connections, check if they respond to ping
	if c.isPaused && time.Since(c.lastSent) > 10*time.Second {
		log.Printf("Connection %p is paused and not responding to ping, should be dropped", c)
		return true
	}

	return false
}

// SetPaused changes the pause state of this connection.
// When paused, the connection will not receive new log broadcasts
// but will still respond to ping messages for health checks.
//
// Parameters:
//   - paused: true to pause the connection, false to resume
func (c *Connection) SetPaused(paused bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		log.Printf("Connection %p is closed, cannot change pause state", c)
		return
	}

	oldPaused := c.isPaused
	c.isPaused = paused
	log.Printf("Connection %p pause state changed: %v → %v", c, oldPaused, paused)
}

// IsPaused checks if this connection is currently paused.
// A paused connection will not receive new log broadcasts
// but maintains the WebSocket connection for health checks.
//
// Returns:
//   - bool: true if the connection is paused, false otherwise
func (c *Connection) IsPaused() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isPaused
}

// SendPing sends a ping message to the client to check connection health.
// This is used for paused connections to verify they are still alive
// and responding to messages.
//
// Returns:
//   - error: nil if ping was sent successfully, error otherwise
func (c *Connection) SendPing() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isClosed {
		return fmt.Errorf("connection is closed")
	}

	pingMessage := model.WebSocketMessage{
		Type: "ping",
		Data: "heartbeat",
	}

	if err := c.ws.WriteJSON(pingMessage); err != nil {
		log.Printf("Failed to send ping to connection %p: %v", c, err)
		return err
	}

	log.Printf("Sent ping to connection %p", c)
	return nil
}

// HandleMessages runs the message handler goroutine for this connection.
// It continuously reads messages from the WebSocket client and processes
// them according to their type (pause, resume, ping, pong).
//
// This method runs in a separate goroutine and handles the complete
// lifecycle of client message processing.
func (c *Connection) HandleMessages() {
	log.Printf("STARTING message handler for connection %p", c)
	defer func() {
		log.Printf("MESSAGE HANDLER EXITING for connection %p", c)
		c.Close()
	}()

	for {
		// Read message from WebSocket
		var message model.WebSocketMessage
		if err := c.ws.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error for connection %p: %v", c, err)
			} else {
				log.Printf("WebSocket closed normally for connection %p: %v", c, err)
			}
			break
		}

		log.Printf("📨 RECEIVED MESSAGE from connection %p: Type=%s, Data=%+v", c, message.Type, message.Data)

		// Handle different message types
		switch message.Type {
		case "ping":
			// Client heartbeat ping, respond with pong
			c.mu.Lock()
			c.lastSent = time.Now()
			c.mu.Unlock()

			pongMessage := model.WebSocketMessage{
				Type: "pong",
				Data: "heartbeat",
			}

			if err := c.ws.WriteJSON(pongMessage); err != nil {
				log.Printf("Failed to send pong to connection %p: %v", c, err)
			} else {
				log.Printf("Connection %p sent ping, responded with pong", c)
			}
		case "pong":
			// Update last sent time when we receive pong (client is alive)
			c.mu.Lock()
			c.lastSent = time.Now()
			c.mu.Unlock()
			log.Printf("Connection %p responded to ping with pong", c)
		case "pause":
			// Update last sent time when client sends pause (client is alive)
			c.mu.Lock()
			c.lastSent = time.Now()
			c.mu.Unlock()
			log.Printf("PROCESSING PAUSE for connection %p", c)
			c.SetPaused(true)
			log.Printf("Connection %p PAUSED successfully", c)
		case "resume":
			// Update last sent time when client sends resume (client is alive)
			c.mu.Lock()
			c.lastSent = time.Now()
			c.mu.Unlock()
			log.Printf("PROCESSING RESUME for connection %p", c)
			c.SetPaused(false)
			log.Printf("Connection %p RESUMED successfully", c)
		default:
			log.Printf("Unknown message type from connection %p: %s", c, message.Type)
		}
	}

	log.Printf("Message handler finished for connection %p", c)
}
