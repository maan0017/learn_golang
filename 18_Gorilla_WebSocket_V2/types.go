package main

import "encoding/json"

type Request struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	// json.RawMessage means “don’t parse payload yet, just keep it as raw JSON until we know the type.”
}

type Response struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (r Response) ToJSON() []byte {
	data, _ := json.Marshal(r)
	return data
}

type MessagePayload struct {
	ID         string `json:"id"`
	Msg        string `json:"msg"`
	SenderID   string `json:"senderId"`
	RecieverID string `json:"recieverId"`
	Timestamp  string `json:"timestamp"`
}

type CoordsPayload struct {
	ClientID   string `json:"clientId"`
	ClientName string `json:"clientName"`
	CoordX     int    `json:"CoordX"`
	CoordY     int    `json:"CoordY"`
}

type NotificationPayload struct {
	Category   string `json:"category"`
	Heading    string `json:"heading"`
	Content    string `json:"content"`
	SenderID   string `json:"senderId"`
	RecieverID string `json:"recieverId"`
	Timestamp  string `json:"timestamp"`
}
