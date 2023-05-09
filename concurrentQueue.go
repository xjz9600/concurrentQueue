package main

import (
	"context"
	"sync"
)

type MyCircularDeque[T any] struct {
	header   int
	tail     int
	capacity int
	data     []T
	zero     T
	enqueue  *sync.Cond
	dequeue  *sync.Cond
}

func NewConstructor[T any](k int) *MyCircularDeque[T] {
	mutex := &sync.Mutex{}
	return &MyCircularDeque[T]{capacity: k, data: make([]T, k+1), enqueue: sync.NewCond(mutex), dequeue: sync.NewCond(mutex)}
}

func (this *MyCircularDeque[T]) InsertFront(ctx context.Context, value T) bool {
	if ctx.Err() != nil {
		return false
	}
	ch := make(chan struct{})
	go func() {
		this.enqueue.L.Lock()
		for this.IsFull() {
			this.enqueue.Wait()
		}
		select {
		case ch <- struct{}{}:
		default:
			this.enqueue.Signal()
			this.enqueue.L.Unlock()
		}
	}()
	select {
	case <-ctx.Done():
		return false
	case <-ch:
		defer this.enqueue.L.Unlock()
		this.data[this.header] = value
		this.header = (this.header + 1) % (this.capacity + 1)
		this.dequeue.Signal()
		return true
	}
}

func (this *MyCircularDeque[T]) GetFront(ctx context.Context) T {
	if ctx.Err() != nil {
		return this.zero
	}
	ch := make(chan struct{})
	go func() {
		this.dequeue.L.Lock()
		for this.IsEmpty() {
			this.dequeue.Wait()
		}
		select {
		case ch <- struct{}{}:
		default:
			this.dequeue.Signal()
			this.dequeue.L.Unlock()
		}
	}()
	select {
	case <-ctx.Done():
		return this.zero
	case <-ch:
		defer this.dequeue.L.Unlock()
		data := this.data[(this.header+this.capacity)%(this.capacity+1)]
		this.header = (this.header + this.capacity) % (this.capacity + 1)
		this.enqueue.Signal()
		return data
	}
}

func (this *MyCircularDeque[T]) IsEmpty() bool {
	if this.header == this.tail {
		return true
	}
	return false
}

func (this *MyCircularDeque[T]) IsFull() bool {
	if (this.header+1)%(this.capacity+1) == this.tail {
		return true
	}
	return false
}
