package main

import (
	"runtime"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/valyala/fasthttp"
)

// Hello ...
type Hello struct{}

// Hello ...
func (*Hello) Hello(name string) string {
	return "Hello " + name + "!"
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	service := rpc.NewFastHTTPService()
	service.AddInstanceMethods(&Hello{}, rpc.Options{})
	//	service.AddFunction("hello", hello, rpc.Options{})
	fasthttp.ListenAndServe(":8080", service.ServeFastHTTP)
}
