package websocket

import (
	"fmt"
	"go/web-sockets/internal/chat"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var clients = make(map[string]*chat.Client)

func registerClient(c *chat.Client) {
	clients[c.ID] = c
}

// Upgrade HTTP connection → WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins, for dev
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	// defer conn.Close()

	// Generate a unique ID
	clientID := uuid.NewString()

	client := &chat.Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan *chat.Message),
	}

	// Store client in your hub/map
	registerClient(client)

	// Listen for messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Println("Received:", string(msg))

		// Echo message back to client
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}
