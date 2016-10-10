package main

import (
	"fmt"
	"reflect"

	"github.com/hprose/hprose-golang/rpc"
)

func main() {
	client := rpc.NewTCPClient("tcp4://127.0.0.1:2016/")
	count := 0
	id, _ := client.ID()
	done := make(chan bool)
	client.Subscribe("time", id, func(result []reflect.Value, err error) {
		count++
		if count > 10 {
			client.Unsubscribe("time")
			done <- true
		}
		fmt.Println(result[0].Interface())
	}, nil)
	<-done
}
