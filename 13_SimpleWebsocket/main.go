package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
}

type Server struct {
	Addr    string
	clients map[string]*Client
	mu      sync.RWMutex
}

func NewServer(addr string) *Server {
	return &Server{
		Addr:    addr,
		clients: make(map[string]*Client),
	}
}

func (s *Server) handleWebSockets(c *websocket.Conn) {
	// runs when a new websocket connections happens
	fmt.Println("✅ New client connected")

	// Track this connection
	id := uuid.NewString()
	client := &Client{
		Id:   id,
		Conn: c,
	}
	s.mu.Lock()
	s.clients[id] = client
	s.mu.Unlock()

	// Start reading messages
	go s.readLoop(client)
}

func (s *Server) readLoop(client *Client) {
	// When this function exits, clean up
	defer func() {
		fmt.Println("❌ Client disconnected")
		s.mu.Lock()
		delete(s.clients, client.Id)
		s.mu.Unlock()
		client.Conn.Close()
	}()

	for {
		var msg string
		// Read message
		if err := websocket.Message.Receive(client.Conn, &msg); err != nil {
			fmt.Println("⚠️ Error reading:", err)
			break
		}

		fmt.Printf("📩 [%s] %s\n", client.Id, msg)

		// Echo back with ID
		reply := fmt.Sprintf("Hello %s 👋, you said: %s", client.Id, msg)
		if err := websocket.Message.Send(client.Conn, reply); err != nil {
			fmt.Println("⚠️ Error writing:", err)
			break
		}
	}
}

// --- Send to one client ---
func (s *Server) sendToClient(id string, msg string) {
	if client, ok := s.clients[id]; ok {
		_ = websocket.Message.Send(client.Conn, msg)
	} else {
		fmt.Printf("⚠️ No client with ID: %s\n", id)
	}
}

// --- Broadcast to all clients ---
func (s *Server) broadcast(msg string, excludeID string) {
	s.mu.RLock()
	for id, client := range s.clients {
		if id == excludeID { // skip sender
			continue
		}
		_ = websocket.Message.Send(client.Conn, msg)
	}
	s.mu.RUnlock()
}

func main() {
	server := NewServer(":8080")

	http.Handle("/ws", websocket.Handler(server.handleWebSockets))

	fmt.Println("🚀 WebSocket server started on ws://localhost:8080/ws")
	if err := http.ListenAndServe(server.Addr, nil); err != nil {
		log.Fatal(err)
	}
}
