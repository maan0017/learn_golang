package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusBadRequest)
	}

	fmt.Fprintf(w, "hello")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseError() err: %+v\n", err)
	}
	fmt.Fprintf(w, "Post Successful\n")

	name := r.FormValue("name")
	email := r.FormValue("email")
	fmt.Fprintf(w, "Name: %s\n", name)
	fmt.Fprintf(w, "Address: %s\n", email)
}

func main() {
	server := NewServer(":3000")
	fileServer := NewFileServer("./static")

	fileHandler := fileServer.StartFileServer()
	// http://localhost:3000/ or http://localhost:3000/index.html -> loads index.html page
	// http://localhost:3000/form.html -> loads form.html page
	http.Handle("/", fileHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server initalized.")

	if err := server.StartServer(); err != nil {
		log.Fatal("failed to start server.")
	}

	fmt.Println("Server is live at http://localhost:3000")
}
