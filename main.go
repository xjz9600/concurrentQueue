package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	TestNewSemaphoreConstructor()
}

func TestNewConstructor() {
	queue := NewConstructor[int](5)
	timeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	for i := 0; i < 20; i++ {
		go func(value int) {
			queue.InsertFront(timeout, value)
		}(i)
		go func() {
			fmt.Println(queue.GetFront(timeout))
		}()
	}
	go func() {
		cancel()
	}()
	time.Sleep(3 * time.Second)
}

func TestNewSemaphoreConstructor() {
	queue := NewSemaphoreConstructor[int](5)
	timeout, _ := context.WithTimeout(context.Background(), 1*time.Second)
	for i := 0; i < 20; i++ {
		go func(value int) {
			queue.InsertFront(timeout, value)
		}(i)
		go func() {
			fmt.Println(queue.GetFront(timeout))
		}()
	}
	//go func() {
	//	cancel()
	//}()
	time.Sleep(3 * time.Second)
}
