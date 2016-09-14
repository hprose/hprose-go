package main

import (
	"fmt"
	"reflect"

	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/valyala/fasthttp"
)

func hello(name string, context *rpc.FastHTTPContext) string {
	return "Hello " + name + "!  -  " + context.RequestCtx.RemoteAddr().String()
}

type A struct {
	S string `json:"str"`
}

func getEmptySlice() interface{} {
	s := make([]A, 100)
	return s
}

type ServerEvent struct{}

func (e *ServerEvent) OnBeforeInvoke(name string, args []reflect.Value, byref bool, context rpc.Context) {
	fmt.Println("Before OK")
}

func (e *ServerEvent) OnAfterInvoke(name string, args []reflect.Value, byref bool, result []reflect.Value, context rpc.Context) {
	fmt.Println("After OK")
}
func (e *ServerEvent) OnSendError(err error, context rpc.Context) {
	fmt.Println(err)
}

func main() {
	io.Register(reflect.TypeOf(A{}), "A", "json")
	service := rpc.NewFastHTTPService()
	service.Event = &ServerEvent{}
	service.Debug = true
	service.AddFunction("hello", hello, rpc.Options{})
	service.AddFunction("getEmptySlice", getEmptySlice, rpc.Options{})
	fasthttp.ListenAndServe(":8080", service.ServeFastHTTP)
}
