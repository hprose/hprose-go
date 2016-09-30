package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hprose/hprose-golang/rpc"
)

// Stub is ...
type Stub struct {
	Hello      func(string) (string, error)
	AsyncHello func(func(string), string) `name:"hello"`
}

func main() {
	client := rpc.NewWebSocketClient("ws://127.0.0.1:8080/")
	client.SetTimeout(1 * time.Second)
	var stub *Stub
	client.UseService(&stub)
	stub.AsyncHello(func(result string) {
		fmt.Println(result)
	}, "async world")
	fmt.Println(stub.Hello("world"))
	start := time.Now()
	var n int32 = 50000
	done := make(chan bool)
	for i := 0; i < 50000; i++ {
		go func() {
			fmt.Println(stub.Hello("world"))
			if atomic.AddInt32(&n, -1) == 0 {
				done <- true
			}
		}()
	}
	<-done
	stop := time.Now()
	fmt.Println((stop.UnixNano() - start.UnixNano()) / 1000000)
}
