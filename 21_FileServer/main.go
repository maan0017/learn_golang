package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct{}

func (fs *FileServer) start(port string) {
	var loopCount uint64
	coutnChan := make(chan uint64, 100)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	// accpect connection in a goroutine
	go func(coutChan chan uint64) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("accpect err: ", err)
				continue
			}

			coutnChan <- 1
			go fs.readLoop(conn)
		}
	}(coutnChan)

	for {
		select {
		case t := <-ticker.C:

			// This will run every 10 seconds
			fmt.Println("Tick at", t)
			fmt.Printf("Loop Count: %d\n", loopCount)
			// Put your task here

		case count := <-coutnChan:
			loopCount += count
		}
	}

}

func (fs *FileServer) readLoop(conn net.Conn) {
	buff := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(conn, buff, size)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buff.Bytes())
		fmt.Printf("recieved %d bytes over the network\n", n)
	}
}

func sendFile(port string, size int64) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", port)
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, size)
	n, err := io.CopyN(conn, bytes.NewReader(file), size)
	if err != nil {
		return err
	}

	fmt.Printf("written %d over the network\n", n)
	return nil
}

func main() {
	port := ":8080"

	go func() {
		time.Sleep(5 * time.Second)
		_ = sendFile(port, 54321)
	}()

	fs := &FileServer{}
	fs.start(port)
}
