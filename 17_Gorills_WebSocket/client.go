package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	name     string
	chatroom string
	conn     *websocket.Conn
	manager  *Manager
	send     chan Event // single channel for outgoing messages
	CoordX   int16
	CoordY   int16
}

var (
	PongWait     = 10 * time.Second
	PingInterval = (PongWait * 9) / 10
)

func (c *Client) readMessages() {
	defer func() {
		c.manager.unregister <- c
		fmt.Println("connection closed (read loop)")
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(PongWait)); err != nil {
		fmt.Printf("not recieved pong: %+v\n", err)
		return
	}

	c.conn.SetReadLimit(1024)

	c.conn.SetPongHandler(func(appData string) error {
		fmt.Println("pong msg recieved: ", appData)
		return c.conn.SetReadDeadline(time.Now().Add(PongWait))
	})

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			break
		}

		// Debug
		fmt.Println("payload:", string(payload))

		// Example: broadcast to all clients
		// response := fmt.Sprintf(`{"type":"message","data":"%s","from":"%s"}`, payload, c.id)
		// c.manager.broadcast <- response
		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			fmt.Printf("error marshalling event: %+v\n", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			fmt.Printf("error handling msg: %+v\n", err)
		}

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.conn.Close() // ✅ ensure socket closes
		fmt.Println("connection closed (write loop)")
	}()

	ticker := time.NewTicker(PingInterval)

	for {
		select {
		case event := <-c.send:
			data, err := json.Marshal(event)
			if err != nil {
				fmt.Printf("marshal error: %+v\n", err)
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to write message: %v", err)
				return
			}

		case <-ticker.C:
			// send a ping to the client
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				fmt.Printf("ping write error: %+v\n", err)
				return
			}
			fmt.Println("ping msg sent")
		}
	}
}
