package main

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
)

func main() {
	client := rpc.NewTCPClient("tcp4://127.0.0.1:2016/")
	count := 0
	id, _ := client.AutoID()
	done := make(chan bool)
	client.Subscribe("time", id, nil, func(data string) {
		count++
		if count > 10 {
			client.Unsubscribe("time")
			done <- true
		}
		fmt.Println(data)
	})
	<-done
}
