package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true, // Vite dev server
	"http://localhost:3000": true, // React dev on 3000
	"http://localhost:8080": true, // go server on 8080
	"https://yourapp.com":   true, // production frontend
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	return allowedOrigins[origin]
}

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}
)

type Manager struct {
	clients       map[*Client]bool
	register      chan *Client
	unregister    chan *Client
	broadcast     chan Event
	otps          RetentionMap
	eventHandlers map[string]EventHandler
	// sync.RWMutex
}

func NewManager() *Manager {
	m := &Manager{
		clients:       make(map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan Event),
		otps:          NewRetentionMap(5 * time.Second),
		eventHandlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.eventHandlers[EventSendMessage] = SendMessage
	// m.eventHandlers[EventWelcome] = WelcomeMessage
	// m.eventHandlers[EventChangeRoom] = ChatRoomHandler
}

func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}
func ChatRoomHandler(event Event, c *Client) error {
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	// check if the event type
	if handler, ok := m.eventHandlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	// verify otp
	// otp := r.URL.Query().Get("otp")
	// if otp == "" && !m.otps.VerifyOtp(otp) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// upgrade regular http connection into regular websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("connection upgrade error: %s \n", err.Error())
		return
	}

	client := &Client{
		id:       uuid.NewString(),
		name:     fmt.Sprintf("user_%d", len(m.clients)+1),
		chatroom: "",
		conn:     conn,
		manager:  m,
		send:     make(chan Event),
		CoordX:   0,
		CoordY:   0,
	}
	// m.addClient(client)
	m.register <- client

	// start client processes
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.clients[client] = true
			fmt.Println("new incomming connection: ", client.conn.RemoteAddr())
			fmt.Printf("Client %s connected. Total clients: %d", client.id, len(m.clients))

			welcomeEvent := WelcomeEvent{
				ClientId: client.id,
				Message:  "Welcome to the WebSocket server!",
			}
			welcomePayload, err := json.Marshal(welcomeEvent)
			if err != nil {
				log.Printf("Failed to marshal welcome event: %v", err)
				return
			}
			e := Event{
				Type:    "welcome",
				Payload: json.RawMessage(welcomePayload),
			}

			// Send welcome message to the new client
			// welcome := fmt.Sprintf(`{"type":"welcome","message":"Welcome to the WebSocket server!","clientId":"%s"}`, client.id)

			client.send <- e

		case client := <-m.unregister:
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send) // ✅ important: close send channel so writeMessages can stop
				client.conn.Close()
				log.Printf("Client %s disconnected. Total clients: %d", client.id, len(m.clients))
			}

		// Broadcast message to all clients
		case event := <-m.broadcast:
			for client := range m.clients {
				select {
				case client.send <- event:
				default:
					delete(m.clients, client)
					close(client.send)
					client.conn.Close()
				}
			}
		}
	}
}

// func (m *Manager) addClient(client *Client) {
// 	m.Lock()
// 	defer m.Unlock()

// 	m.clients[client] = true
// }

// func (m *Manager) removeClient(client *Client) {
// 	m.Lock()
// 	defer m.Unlock()

// 	if _, ok := m.clients[client]; ok {
// 		client.conn.Close()
// 		delete(m.clients, client)
// 	}
// }

func (m *Manager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type UserLoginCreds struct {
		Usernaem string `json:"username"`
		Password string `json:"password"`
	}

	var userLoginCred UserLoginCreds

	if err := json.NewDecoder(r.Body).Decode(&userLoginCred); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if userLoginCred.Usernaem != "user" && userLoginCred.Password != "1234" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type SuccessFullLoginResponse struct {
		Otp string `json:"otp"`
	}

	otp := m.otps.NewOtp()

	resp := SuccessFullLoginResponse{
		Otp: otp.Key,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("auth marshal error: %+v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
