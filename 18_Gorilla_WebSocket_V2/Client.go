package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	id   string
}

// readPump handles reading messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v\n", err)
			}
			log.Println("read error:", err)
			break
		}

		var req Request
		if err := json.Unmarshal(payload, &req); err != nil {
			log.Println("invalid request:", err)
			continue
		}

		switch req.Type {
		case "message":
			var msg MessagePayload
			if err := json.Unmarshal(req.Payload, &msg); err != nil {
				log.Println("invalid message payload:", err)
				continue
			}
			fmt.Printf("Got message: %+v\n", msg)

			// Example: broadcast it back
			res := Response{
				Type:    "message",
				Payload: msg,
			}

			c.hub.broadcast <- res.ToJSON()

		case "coords":
			var coords CoordsPayload
			if err := json.Unmarshal(req.Payload, &coords); err != nil {
				log.Println("invalid message payload:", err)
				continue
			}
			fmt.Printf("Got coords: %+v\n", coords)

			// Example: broadcast it back
			res := Response{
				Type:    "coords",
				Payload: coords,
			}

			c.hub.broadcast <- res.ToJSON()
		}

		// log.Printf("Received from %s: %s\n", c.id, payload)

		// Echo the message back to all clients
		// response := fmt.Sprintf(`{"type":"message","payload":{"msg":"%s","from":"%s"}}`, string(payload), c.id)
		// c.hub.broadcast <- []byte(response)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	// for {
	// 	select {
	// 	case message, ok := <-c.send:
	// 		if !ok {
	// 			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	// 			return
	// 		}
	// 		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 	}
	// }

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
			return
		}
	}
	// Channel closed, send close message
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})

}
