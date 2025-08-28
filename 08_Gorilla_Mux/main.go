package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := NewMuxRouter()

	r.Routes()

	fmt.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", r.router))
}
