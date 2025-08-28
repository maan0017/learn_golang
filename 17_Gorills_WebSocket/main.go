package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	manager := NewManager()
	go manager.Run()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.LoginHandler)

	fmt.Println("web socket is live at ws://localhost:8080/ws")
	fmt.Println("server is live at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
