package main

import (
	"fmt"
	"sync"
	"time"
)

type Order struct {
	TableNumber int
	PrepTime    time.Duration
}

func processOrder(order Order) {
	// Simulate cooking time
	fmt.Printf("Preparing order for table %d...\n", order.TableNumber)

	time.Sleep(order.PrepTime)

	fmt.Printf("Order ready for table %d!\n\n", order.TableNumber)
}

func main() {
	fmt.Println("Hello mf")

	orders := []Order{
		{TableNumber: 1, PrepTime: 3 * time.Second},
		{TableNumber: 2, PrepTime: 9 * time.Second},
		{TableNumber: 3, PrepTime: 5 * time.Second},
		{TableNumber: 4, PrepTime: 2 * time.Second},
		{TableNumber: 5, PrepTime: 4 * time.Second},
		{TableNumber: 6, PrepTime: 8 * time.Second},
		{TableNumber: 7, PrepTime: 1 * time.Second},
		{TableNumber: 8, PrepTime: 3 * time.Second},
		{TableNumber: 9, PrepTime: 15 * time.Second},
		{TableNumber: 10, PrepTime: 2 * time.Second},
		{TableNumber: 11, PrepTime: 4 * time.Second},
		{TableNumber: 12, PrepTime: 11 * time.Second},
		{TableNumber: 13, PrepTime: 3 * time.Second},
		{TableNumber: 14, PrepTime: 5 * time.Second},
		{TableNumber: 15, PrepTime: 2 * time.Second},
		{TableNumber: 16, PrepTime: 4 * time.Second},
		{TableNumber: 17, PrepTime: 7 * time.Second},
		{TableNumber: 18, PrepTime: 5 * time.Second},
	}

	wg := sync.WaitGroup{}
	// wg.Add(len(orders))

	for _, order := range orders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processOrder(order)
		}()
	}

	wg.Wait()
	fmt.Println("Executed Successfully.")
}
