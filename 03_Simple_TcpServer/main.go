package main

import (
	"fmt"
	"net"
)

func main() {
	ls, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err) // stop program if we can't bind to port
	}
	defer ls.Close()

	fmt.Println("TCP server running on port 8080...")

	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go HandleConn(conn)
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	defer fmt.Println("1st defer")
	defer fmt.Println("2nd defer")
	defer fmt.Println("3rd defer")

	fmt.Printf("New connection: %+v\n", conn.RemoteAddr())
	fmt.Fprintln(conn, "This is a TCP server")
}
