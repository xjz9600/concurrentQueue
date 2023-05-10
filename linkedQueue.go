package main

import (
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
	"unsafe"
)

type concurrentLinkedQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node[T any] struct {
	val  T
	next unsafe.Pointer
}

func NewConcurrentLinkedQueue[T any]() *concurrentLinkedQueue[T] {
	node := &node[T]{}
	ptr := unsafe.Pointer(node)
	return &concurrentLinkedQueue[T]{
		head: ptr,
		tail: ptr,
	}
}

func (c *concurrentLinkedQueue[T]) Dequeue() (T, error) {
	for {
		headPtr := atomic.LoadPointer(&c.head)
		head := (*node[T])(headPtr)
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		if head == tail {
			// 不需要做更多检测，在当下这一刻，我们就认为没有元素，即便这时候正好有人入队
			// 但是并不妨碍我们在它彻底入队完成——即所有的指针都调整好——之前，
			// 认为其实还是没有元素
			var t T
			return t, errors.New("队列为空")
		}
		headNextPtr := atomic.LoadPointer(&head.next)
		if atomic.CompareAndSwapPointer(&c.head, headPtr, headNextPtr) {
			headNext := (*node[T])(headNextPtr)
			return headNext.val, nil
		}
	}
}

func (c *concurrentLinkedQueue[T]) EnSequenceQueue(t T) error {
	newNode := &node[T]{val: t}
	newPtr := unsafe.Pointer(newNode)
	for {
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		tailNext := atomic.LoadPointer(&tail.next)
		if tailNext != nil {
			continue
		}
		if atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr) {
			time.Sleep(2 * time.Second)
			atomic.CompareAndSwapPointer(&tail.next, tailNext, newPtr)
			return nil
		}
	}
}

func (c *concurrentLinkedQueue[T]) Enqueue(t T) error {
	newNode := &node[T]{val: t}
	newPtr := unsafe.Pointer(newNode)
	for {
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		tailNext := atomic.LoadPointer(&tail.next)
		if tailNext != nil {
			continue
		}
		if atomic.CompareAndSwapPointer(&tail.next, tailNext, newPtr) {
			atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr)
			return nil
		}
	}
}
