package main

import (
	"runtime"

	"github.com/hprose/hprose-golang/rpc"
)

func hello(name string) string {
	return "Hello " + name + "!"
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	server := rpc.NewTCPServer("tcp4://0.0.0.0:4321/")
	server.AddFunction("hello", hello, rpc.Options{})
	server.Start()
}
