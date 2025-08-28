package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Listen for incoming messages
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Print("Peer: " + msg)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Host (h) or Connect (c)? ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "h" {
		// Act as server
		ln, _ := net.Listen("tcp", ":12345")
		fmt.Println("Listening on port 12345...")
		conn, _ := ln.Accept()
		go handleConnection(conn)

		for {
			text, _ := reader.ReadString('\n')
			conn.Write([]byte(text))
		}

	} else {
		// Act as client
		fmt.Print("Enter host (e.g., 127.0.0.1:12345): ")
		host, _ := reader.ReadString('\n')
		host = strings.TrimSpace(host)

		conn, _ := net.Dial("tcp", host)
		go handleConnection(conn)

		for {
			text, _ := reader.ReadString('\n')
			conn.Write([]byte(text))
		}
	}
}
