package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("starting websocket server...")

	hub := NewHub()
	go hub.Run()

	// WebSocket endpoint
	http.HandleFunc("/ws", hub.handleWebSocket)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running"))
	})

	// Enable CORS for all endpoints
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	port := ":8080"
	fmt.Printf("WebSocket server starting on port %s\n", port)
	fmt.Printf("WebSocket endpoint: ws://localhost%s/ws\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
