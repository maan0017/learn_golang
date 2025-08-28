package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var connectionUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		// In production, you should restrict this to your frontend domain
		return true
	},
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client %s connected. Total clients: %d", client.id, len(h.clients))

			// Send welcome message to the new client
			welcome := fmt.Sprintf(`{"type":"welcome","payload":{"message":"Welcome to the WebSocket server!","clientId":"%s"}}`, client.id)

			select {
			case client.send <- []byte(welcome):
			default:
				close(client.send)
				delete(h.clients, client)
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client %s disconnected. Total clients: %d", client.id, len(h.clients))
			}

		// Broadcast message to all clients
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// handleWebSocket handles WebSocket requests from clients
func (h *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := connectionUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Generate a simple client ID (in production, use UUID or similar)
	clientID := fmt.Sprintf("client_%d", len(h.clients)+1)

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
		id:   clientID,
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}
