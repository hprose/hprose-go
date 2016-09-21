package main

import (
	"runtime"
	"time"

	"github.com/hprose/hprose-golang/rpc"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	server := rpc.NewTCPServer("tcp4://0.0.0.0:2016/")
	server.Publish("time", 0, 0)
	var timer *time.Timer
	timer = time.AfterFunc(1*1000*1000*1000, func() {
		server.Push("time", time.Now().String())
		timer.Reset(1 * 1000 * 1000 * 1000)
	})
	server.Start()
}
