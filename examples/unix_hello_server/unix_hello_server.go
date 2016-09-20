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
	server := rpc.NewUnixServer("unix:/tmp/my.sock")
	server.AddFunction("hello", hello, rpc.Options{})
	server.Start()
}
