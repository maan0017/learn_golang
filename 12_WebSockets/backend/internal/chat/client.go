package chat

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Tune these as needed
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4 << 10 // 4KB
)

var newline = []byte{'\n'}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID   string
	Conn *websocket.Conn
	Hub  *Hub
	Send chan *Message
	// mu   sync.Mutex
}

func NewClient(h *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ID:   uuid.NewString(),
		Conn: conn,
		Hub:  h,
		Send: make(chan *Message, 256),
	}
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			// expected on close
			break
		}

		msg := &Message{
			ID:        uuid.NewString(),
			SenderID:  c.ID,
			Message:   string(bytes.TrimSpace(message)),
			Timestamp: time.Now().Format(time.RFC3339),
		}

		c.Hub.Broadcast <- msg
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(msg) // serialize once here
			if err != nil {
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
