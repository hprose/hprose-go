package main

import (
	"net/http"
	"runtime"

	"github.com/hprose/hprose-golang/rpc"
)

func hello(name string) string {
	return "Hello " + name + "!"
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	service := rpc.NewHTTPService()
	service.AddFunction("hello", hello, rpc.Options{})
	http.ListenAndServe(":8080", service)
}
