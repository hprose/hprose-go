package main

import (
	"github.com/hprose/hprose-golang/rpc"
	"github.com/valyala/fasthttp"
)

func hello(name string) string {
	return "Hello " + name + "!"
}

func main() {
	service := rpc.NewFastHTTPService()
	service.AddFunction("hello", hello, rpc.Options{})
	fasthttp.ListenAndServe(":8080", service.ServeFastHTTP)
}
