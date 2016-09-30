package main

import (
	"fmt"
	"time"

	"github.com/hprose/hprose-golang/rpc"
)

// Stub is ...
type Stub struct {
	Hello      func(string) string
	AsyncHello func(func(string), string) `name:"hello"`
}

func main() {
	client := rpc.NewFastHTTPClient("http://127.0.0.1:8080/")
	var stub *Stub
	client.UseService(&stub)
	stub.AsyncHello(func(result string) {
		fmt.Println(result)
	}, "async world")
	fmt.Println(stub.Hello("world"))
	start := time.Now()
	for i := 0; i < 50000; i++ {
		stub.Hello("world")
	}
	stop := time.Now()
	fmt.Println((stop.UnixNano() - start.UnixNano()) / 1000000)
}
