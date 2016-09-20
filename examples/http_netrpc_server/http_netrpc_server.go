package main

import (
	"errors"
	"net/http"
	"runtime"

	"github.com/hprose/hprose-golang/rpc"
)

// Args ...
type Args struct {
	A, B int
}

// Quotient ...
type Quotient struct {
	Quo, Rem int
}

// Arith ...
type Arith int

// Multiply ...
func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

// Divide ...
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	service := rpc.NewHTTPService()
	service.AddNetRPCMethods(new(Arith), rpc.Options{})
	http.ListenAndServe(":8080", service)
}
