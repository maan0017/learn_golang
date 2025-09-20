package main

import (
	"log"
	"net/http"
)

type UserCred struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	port := ":8080"
	mux := http.NewServeMux()

	log.Println("server is live at port ", port)

	mux.HandleFunc("GET /", MainHandler)

	mux.HandleFunc("GET /{id}", GetUserByID)

	wrappedMux := LoggingMiddleware(mux)

	if err := http.ListenAndServe(port, wrappedMux); err != nil {
		log.Fatal(err)
	}
}
