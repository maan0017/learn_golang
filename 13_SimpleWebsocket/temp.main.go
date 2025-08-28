package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"

// 	"golang.org/x/net/websocket"
// )

// type Server struct {
// 	conns map[*websocket.Conn]bool
// }

// func NewServer() *Server {
// 	return &Server{
// 		conns: make(map[*websocket.Conn]bool),
// 	}
// }

// func (s *Server) handleWS(ws *websocket.Conn) {
// 	fmt.Println("new incoming connnection from client:", ws.RemoteAddr())

// 	s.conns[ws] = true

// 	s.readLoop(ws)
// }

// func (s *Server) readLoop(ws *websocket.Conn) {
// 	buff := make([]byte, 1024)

// 	// for loop keeps looping infinitely //
// 	for {
// 		n, err := ws.Read(buff)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			fmt.Println("read error: ", err)
// 			continue
// 		}

// 		msg := buff[:n]
// 		// fmt.Println(string(msg))

// 		// ws.Write([]byte("thank you for the msg"))

// 		s.broadcast(msg)
// 	}
// }

// func (s *Server) broadcast(b []byte) {
// 	for ws := range s.conns {
// 		go func(ws *websocket.Conn) {
// 			if _, err := ws.Write(b); err != nil {
// 				fmt.Println("broadcast error: ", err)
// 			}
// 		}(ws)
// 	}
// }

// func main() {
// 	fmt.Println("hello mf")

// 	s := NewServer()
// 	http.Handle("/ws", websocket.Handler(s.handleWS))
// 	if err := http.ListenAndServe(":3000", nil); err != nil {
// 		log.Fatal("failed to start the server: ", err)
// 	}
// 	fmt.Println("server is live at port :3000")
// }
