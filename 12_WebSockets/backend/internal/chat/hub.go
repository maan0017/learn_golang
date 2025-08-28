package chat

import "time"

type Message struct {
	ID        string `json:"id"`
	SenderID  string `json:"senderId"`
	Sender    string `json:"sender,omitempty"` // optional: username
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	System    bool   `json:"system"`
}

// Hub coordinates all clients and broadcasts messages.
type Hub struct {
	// Registered clients
	Clients map[*Client]struct{}
	// Clients map[string]*Client
	// Inbound messages from clients
	Broadcast chan *Message
	// Register requests from clients
	Register chan *Client
	// Unregister requests from clients
	Unregister chan *Client

	Closed chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]struct{}),
		Broadcast:  make(chan *Message, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Closed:     make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = struct{}{}
			welcome := &Message{
				ID:        c.ID,      // unique connection id
				SenderID:  "server",  // mark as server
				Message:   "WELCOME", // just a flag message
				Timestamp: time.Now().Format(time.RFC3339),
				System:    true,
			}
			c.Send <- welcome

		case c := <-h.Unregister:
			if _, ok := h.Clients[c]; ok {
				delete(h.Clients, c)
				close(c.Send)
			}
		case msg := <-h.Broadcast:
			for c := range h.Clients {
				select {
				case c.Send <- msg:
				default:
					delete(h.Clients, c)
					close(c.Send)
				}
			}
		case <-h.Closed:
			for c := range h.Clients {
				close(c.Send)
				delete(h.Clients, c)
			}
			return
		}
	}
}

func (h *Hub) Close() {
	close(h.Closed)
}
