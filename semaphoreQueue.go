package main

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
)

type MyCirculaSemaphore[T any] struct {
	header   int
	tail     int
	capacity int
	data     []T
	zero     T
	enqueue  *semaphore.Weighted
	dequeue  *semaphore.Weighted
	mutex    *sync.RWMutex
}

func NewSemaphoreConstructor[T any](k int) *MyCirculaSemaphore[T] {
	res := &MyCirculaSemaphore[T]{
		capacity: k,
		data:     make([]T, k+1),
		mutex:    &sync.RWMutex{},
		enqueue:  semaphore.NewWeighted(int64(k)),
		dequeue:  semaphore.NewWeighted(int64(k)),
	}
	res.dequeue.Acquire(context.TODO(), int64(k))
	return res
}

func (this *MyCirculaSemaphore[T]) InsertFront(ctx context.Context, value T) bool {
	err := this.enqueue.Acquire(ctx, 1)
	if err != nil {
		return false
	}
	this.mutex.Lock()
	if ctx.Err() != nil {
		return false
	}
	defer this.mutex.Unlock()
	this.data[this.header] = value
	this.header = (this.header + 1) % (this.capacity + 1)
	this.dequeue.Release(1)
	return true
}

func (this *MyCirculaSemaphore[T]) GetFront(ctx context.Context) T {
	err := this.dequeue.Acquire(ctx, 1)
	if err != nil {
		return this.zero
	}
	this.mutex.Lock()
	if ctx.Err() != nil {
		return this.zero
	}
	defer this.mutex.Unlock()
	data := this.data[(this.header+this.capacity)%(this.capacity+1)]
	this.header = (this.header + this.capacity) % (this.capacity + 1)
	this.enqueue.Release(1)
	return data

}

func (this *MyCirculaSemaphore[T]) IsEmpty() bool {
	if this.header == this.tail {
		return true
	}
	return false
}

func (this *MyCirculaSemaphore[T]) IsFull() bool {
	if (this.header+1)%(this.capacity+1) == this.tail {
		return true
	}
	return false
}
