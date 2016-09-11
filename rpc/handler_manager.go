/**********************************************************\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: http://www.hprose.com/                 |
|                   http://www.hprose.org/                 |
|                                                          |
\**********************************************************/
/**********************************************************\
 *                                                        *
 * rpc/handler_manager.go                                 *
 *                                                        *
 * hprose handler manager for Go.                         *
 *                                                        *
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"reflect"

	"github.com/hprose/hprose-golang/promise"
)

// NextInvokeHandler is the next invoke handler function
type NextInvokeHandler func(name string, args []reflect.Value, context Context) promise.Promise

// InvokeHandler is the invoke handler function
type InvokeHandler func(name string, args []reflect.Value, context Context, next NextInvokeHandler) promise.Promise

// NextFilterHandler is the next filter handler function
// The result type is promise.Promise<[]byte>
type NextFilterHandler func(request []byte, context Context) promise.Promise

// FilterHandler is the filter handler function
// The result type is promise.Promise<[]byte>
type FilterHandler func(request []byte, context Context, next NextFilterHandler) promise.Promise

// HandlerManager is the hprose handler manager
type HandlerManager struct {
	invokeHandlers             []InvokeHandler
	beforeFilterHandlers       []FilterHandler
	afterFilterHandlers        []FilterHandler
	defaultInvokeHandler       NextInvokeHandler
	defaultBeforeFilterHandler NextFilterHandler
	defaultAfterFilterHandler  NextFilterHandler
	invokeHandler              NextInvokeHandler
	beforeFilterHandler        NextFilterHandler
	afterFilterHandler         NextFilterHandler
	BeforeFilter               *FilterHandlerMangager
	AfterFilter                *FilterHandlerMangager
	override                   struct {
		invokeHandler       NextInvokeHandler
		beforeFilterHandler NextFilterHandler
		afterFilterHandler  NextFilterHandler
	}
}

// FilterHandlerMangager is the hprose filter handler manager
type FilterHandlerMangager struct {
	addFilterHandler func(handler FilterHandler)
}

// Use is the same as AddBeforeFilterHandler/AddAfterFilterHandler
// but it can be called by a chain
func (fhm *FilterHandlerMangager) Use(
	handler FilterHandler) *FilterHandlerMangager {
	fhm.addFilterHandler(handler)
	return fhm
}

// NewHandlerManager is the constructor of HandlerManager
func NewHandlerManager() (hm *HandlerManager) {
	hm = new(HandlerManager)
	hm.defaultInvokeHandler = func(name string, args []reflect.Value, context Context) promise.Promise {
		return hm.override.invokeHandler(name, args, context)
	}
	hm.defaultBeforeFilterHandler = func(request []byte, context Context) promise.Promise {
		return hm.override.beforeFilterHandler(request, context)
	}
	hm.defaultAfterFilterHandler = func(request []byte, context Context) promise.Promise {
		return hm.override.afterFilterHandler(request, context)
	}
	hm.invokeHandler = hm.defaultInvokeHandler
	hm.beforeFilterHandler = hm.defaultBeforeFilterHandler
	hm.afterFilterHandler = hm.defaultAfterFilterHandler
	hm.BeforeFilter = &FilterHandlerMangager{hm.AddBeforeFilterHandler}
	hm.AfterFilter = &FilterHandlerMangager{hm.AddAfterFilterHandler}
	return
}

func getNextInvokeHandler(
	next NextInvokeHandler, handler InvokeHandler) NextInvokeHandler {
	return func(name string, args []reflect.Value, context Context) (result promise.Promise) {
		defer func() {
			if e := recover(); e != nil {
				result = promise.Reject(promise.NewPanicError(e))
			}
		}()
		return handler(name, args, context, next)
	}
}

func getNextFilterHandler(
	next NextFilterHandler, handler FilterHandler) NextFilterHandler {
	return func(request []byte, context Context) (result promise.Promise) {
		defer func() {
			if e := recover(); e != nil {
				result = promise.Reject(promise.NewPanicError(e))
			}
		}()
		return handler(request, context, next)
	}
}

// AddInvokeHandler add the invoke handler
func (hm *HandlerManager) AddInvokeHandler(handler InvokeHandler) {
	if handler == nil {
		return
	}
	hm.invokeHandlers = append(hm.invokeHandlers, handler)
	next := hm.defaultInvokeHandler
	for i := len(hm.invokeHandlers) - 1; i >= 0; i-- {
		next = getNextInvokeHandler(next, hm.invokeHandlers[i])
	}
	hm.invokeHandler = next
}

// AddBeforeFilterHandler add the filter handler before filters
func (hm *HandlerManager) AddBeforeFilterHandler(handler FilterHandler) {
	if handler == nil {
		return
	}
	hm.beforeFilterHandlers = append(hm.beforeFilterHandlers, handler)
	next := hm.defaultBeforeFilterHandler
	for i := len(hm.beforeFilterHandlers) - 1; i >= 0; i-- {
		next = getNextFilterHandler(next, hm.beforeFilterHandlers[i])
	}
	hm.beforeFilterHandler = next
}

// AddAfterFilterHandler add the filter handler after filters
func (hm *HandlerManager) AddAfterFilterHandler(handler FilterHandler) {
	if handler == nil {
		return
	}
	hm.afterFilterHandlers = append(hm.afterFilterHandlers, handler)
	next := hm.defaultAfterFilterHandler
	for i := len(hm.afterFilterHandlers) - 1; i >= 0; i-- {
		next = getNextFilterHandler(next, hm.afterFilterHandlers[i])
	}
	hm.afterFilterHandler = next
}

// Use is the same as AddInvokeHandler, but it can be called by a chain
func (hm *HandlerManager) Use(handler InvokeHandler) *HandlerManager {
	hm.AddInvokeHandler(handler)
	return hm
}
