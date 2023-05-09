package main

import (
	"context"
	"fmt"
)

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(Handler) Handler

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

func TestChain() {
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		fmt.Println("reply start")
		return "reply", nil
	}
	Chain(test1Middleware, test2Middleware, test3Middleware)(next)
	//got, err := Chain(test1Middleware, test2Middleware, test3Middleware)(next)(context.Background(), "hello kratos!")
	//if err != nil {
	//	fmt.Println("err is:", err)
	//}
	//fmt.Println(got)
}

func test1Middleware(handler Handler) Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("test1 before")
		reply, err = handler(ctx, req)
		fmt.Println("test1 after")
		return
	}
}

func test2Middleware(handler Handler) Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("test2 before")
		reply, err = handler(ctx, req)
		fmt.Println("test2 after")
		return
	}
}

func test3Middleware(handler Handler) Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("test3 before")
		reply, err = handler(ctx, req)
		fmt.Println("test3 after")
		return
	}
}
