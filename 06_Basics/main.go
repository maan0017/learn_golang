package main

import (
	"fmt"
	"sync"
	"time"
)

type Book struct {
	Name   string
	Author string

	Published time.Time
	Wg        sync.WaitGroup
}

func (book *Book) setName(name string, t time.Duration) {
	defer book.Wg.Done()
	time.Sleep(t * time.Second)
	book.Name = name
}

func (book *Book) setAuthor(name string) {
	defer book.Wg.Done()
	time.Sleep(3 * time.Second)
	book.Author = name
}

// func (book *Book) getName() string {
// 	return book.Name
// }

func main() {
	book1 := Book{}
	book1.Wg.Add(2) // tell WaitGroup we have 1 goroutine to wait for

	book2 := Book{}
	book2.Wg.Add(1)

	go book1.setName("raman", 1)
	go book1.setAuthor("Mr. Raman Nauadu")

	go book2.setName("Hallo", 5)
	book2.Wg.Wait()
	fmt.Println("Book 2 Name:", book2.Name)

	book1.Wg.Wait()

	fmt.Println("Book 1 Name:", book1.Name)
	fmt.Println("Book 1 Author Name:", book1.Author)
}
