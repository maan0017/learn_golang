package main

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventMessage     = "message"
	EventSendMessage = "send_message"
	EventWelcome     = "welcome"
	EventCoords      = "coords"
	EventChatRoom    = "chat_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type WelcomeEvent struct {
	ClientId string `json:"clientId"`
	Message  string `json:"message"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChatRoomChangeEvent struct {
	Name string `json:"name"`
}
