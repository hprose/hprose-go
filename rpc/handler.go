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
 * rpc/handler.go                                         *
 *                                                        *
 * hprose handler manager for Go.                         *
 *                                                        *
 * LastModified: Sep 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import "reflect"

// NextInvokeHandler is the next invoke handler function
type NextInvokeHandler func(
	name string,
	args []reflect.Value,
	context Context) (results []reflect.Value, err error)

// InvokeHandler is the invoke handler function
type InvokeHandler func(
	name string,
	args []reflect.Value,
	context Context,
	next NextInvokeHandler) (results []reflect.Value, err error)

// NextFilterHandler is the next filter handler function
type NextFilterHandler func(
	request []byte,
	context Context) (response []byte, err error)

// FilterHandler is the filter handler function
type FilterHandler func(
	request []byte,
	context Context,
	next NextFilterHandler) (response []byte, err error)

// handlerManager is the hprose handler manager
type handlerManager struct {
	invokeHandlers             []InvokeHandler
	beforeFilterHandlers       []FilterHandler
	afterFilterHandlers        []FilterHandler
	defaultInvokeHandler       NextInvokeHandler
	defaultBeforeFilterHandler NextFilterHandler
	defaultAfterFilterHandler  NextFilterHandler
	invokeHandler              NextInvokeHandler
	beforeFilterHandler        NextFilterHandler
	afterFilterHandler         NextFilterHandler
	override                   struct {
		invokeHandler       NextInvokeHandler
		beforeFilterHandler NextFilterHandler
		afterFilterHandler  NextFilterHandler
	}
}

// newHandlerManager is the constructor of HandlerManager
func newHandlerManager() (hm *handlerManager) {
	hm = new(handlerManager)
	hm.defaultInvokeHandler = func(
		name string,
		args []reflect.Value,
		context Context) (results []reflect.Value, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = NewPanicError(e)
			}
		}()
		return hm.override.invokeHandler(name, args, context)
	}
	hm.defaultBeforeFilterHandler = func(
		request []byte,
		context Context) (response []byte, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = NewPanicError(e)
			}
		}()
		return hm.override.beforeFilterHandler(request, context)
	}
	hm.defaultAfterFilterHandler = func(
		request []byte,
		context Context) (response []byte, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = NewPanicError(e)
			}
		}()
		return hm.override.afterFilterHandler(request, context)
	}
	hm.invokeHandler = hm.defaultInvokeHandler
	hm.beforeFilterHandler = hm.defaultBeforeFilterHandler
	hm.afterFilterHandler = hm.defaultAfterFilterHandler
	return
}

func getNextInvokeHandler(
	next NextInvokeHandler, handler InvokeHandler) NextInvokeHandler {
	return func(name string,
		args []reflect.Value,
		context Context) (results []reflect.Value, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = NewPanicError(e)
			}
		}()
		return handler(name, args, context, next)
	}
}

func getNextFilterHandler(
	next NextFilterHandler, handler FilterHandler) NextFilterHandler {
	return func(request []byte, context Context) (response []byte, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = NewPanicError(e)
			}
		}()
		return handler(request, context, next)
	}
}

// AddInvokeHandler add the invoke handler
func (hm *handlerManager) AddInvokeHandler(handler InvokeHandler) {
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
func (hm *handlerManager) AddBeforeFilterHandler(handler FilterHandler) {
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
func (hm *handlerManager) AddAfterFilterHandler(handler FilterHandler) {
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
