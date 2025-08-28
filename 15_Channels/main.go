package main

import (
	"fmt"
	"sync"
	"time"
)

func sendMessage(num int, msgChan chan<- string) {
	fmt.Printf("Sending message %d\n", num)

	time.Sleep(time.Second * time.Duration(num)) // Simulate some work

	msg := fmt.Sprintf("✅ Message %d sent!", num)

	msgChan <- msg
}

func receiveMessage(msgs <-chan string) {
	fmt.Println("Waiting for message")

	for msg := range msgs {
		fmt.Println("Received:", msg)
	}
}

func main() {
	fmt.Println("Hello Hi")

	msgChan := make(chan string)
	defer close(msgChan)
	// chan<- means that we are going to send data to channel
	// <-chan means that we are going to recieve data from channel

	wg := sync.WaitGroup{}
	wg.Add(2)

	go sendMessage(10, msgChan)
	go sendMessage(5, msgChan)

	go func() {
		defer wg.Done()
		receiveMessage(msgChan)
	}()

	wg.Wait()

	fmt.Println("Executed successfully.")

}
