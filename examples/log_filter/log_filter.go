package main

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
)

type LogFilter struct {
	Prompt string
}

func (lf LogFilter) InputFilter(data []byte, context rpc.Context) []byte {
	fmt.Printf("%v: %s\r\n", lf.Prompt, data)
	return data
}

func (lf LogFilter) OutputFilter(data []byte, context rpc.Context) []byte {
	fmt.Printf("%v: %s\r\n", lf.Prompt, data)
	return data
}
