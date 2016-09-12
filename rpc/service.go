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
 * rpc/service.go                                         *
 *                                                        *
 * hprose service for Go.                                 *
 *                                                        *
 * LastModified: Sep 12, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/promise"
)

// Service interface
type Service interface {
	AddFunction(name string, function interface{}, options MethodOptions) Service
	AddFunctions(names []string, functions []interface{}, options MethodOptions) Service
	AddMethod(name string, obj interface{}, options MethodOptions, alias ...string) Service
	AddMethods(names []string, obj interface{}, options MethodOptions, aliases ...[]string) Service
	AddInstanceMethods(obj interface{}, options MethodOptions) Service
	AddAllMethods(obj interface{}, options MethodOptions) Service
	AddMissingMethod(method MissingMethod, options MethodOptions) Service
	Remove(name string) Service
	Filter() Filter
	FilterByIndex(index int) Filter
	SetFilter(filter ...Filter) Service
	AddFilter(filter ...Filter) Service
	RemoveFilterByIndex(index int) Service
	RemoveFilter(filter ...Filter) Service
	AddInvokeHandler(handler InvokeHandler) Service
	AddBeforeFilterHandler(handler FilterHandler) Service
	AddAfterFilterHandler(handler FilterHandler) Service
}

type fixer interface {
	FixArguments(args []reflect.Value, lastParamType reflect.Type, context Context) []reflect.Value
}

func fixArguments(args []reflect.Value, lastParamType reflect.Type, context Context) []reflect.Value {
	typ := (*emptyInterface)(unsafe.Pointer(&lastParamType)).ptr
	if typ == interfaceType || typ == contextType {
		args = append(args, reflect.ValueOf(context))
	}
	return args
}

// BaseService is the hprose base service
type BaseService struct {
	*methodManager
	*handlerManager
	*filterManager
	fixer
	Event      ServiceEvent
	Debug      bool
	Simple     bool
	Timeout    time.Duration
	Heartbeat  time.Duration
	ErrorDelay time.Duration
	allTopics  map[string]map[string]*topic
}

// GetNextID is the default method for client uid
func GetNextID() (uid string) {
	u := make([]byte, 16)
	rand.Read(u)
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	uid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// NewBaseService is the constructor for BaseService
func NewBaseService() (service *BaseService) {
	service = new(BaseService)
	service.methodManager = newMethodManager()
	service.handlerManager = newHandlerManager()
	service.filterManager = &filterManager{}
	service.Timeout = 120 * 1000 * 1000
	service.Heartbeat = 3 * 1000 * 1000
	service.ErrorDelay = 10 * 1000 * 1000
	service.allTopics = make(map[string]map[string]*topic)
	service.AddFunction("#", GetNextID, MethodOptions{Simple: true})
	return service
}

// AddFunction publish a func or bound method
// name is the method name
// function is a func or bound method
// options includes Mode, Simple, Oneway and NameSpace
func (service *BaseService) AddFunction(name string, function interface{}, options MethodOptions) Service {
	service.methodManager.AddFunction(name, function, options)
	return service
}

// AddFunctions is used for batch publishing service method
func (service *BaseService) AddFunctions(names []string, functions []interface{}, options MethodOptions) Service {
	service.methodManager.AddFunctions(names, functions, options)
	return service
}

// AddMethod is used for publishing a method on the obj with an alias
func (service *BaseService) AddMethod(name string, obj interface{}, options MethodOptions, alias ...string) Service {
	service.methodManager.AddMethod(name, obj, options, alias...)
	return service
}

// AddMethods is used for batch publishing methods on the obj with aliases
func (service *BaseService) AddMethods(names []string, obj interface{}, options MethodOptions, aliases ...[]string) Service {
	service.methodManager.AddMethods(names, obj, options, aliases...)
	return service
}

// AddInstanceMethods is used for publishing all the public methods and func fields with options.
func (service *BaseService) AddInstanceMethods(obj interface{}, options MethodOptions) Service {
	service.methodManager.AddInstanceMethods(obj, options)
	return service
}

// AddAllMethods will publish all methods and non-nil function fields on the
// obj self and on its anonymous or non-anonymous struct fields (or pointer to
// pointer ... to pointer struct fields). This is a recursive operation.
// So it's a pit, if you do not know what you are doing, do not step on.
func (service *BaseService) AddAllMethods(obj interface{}, options MethodOptions) Service {
	service.methodManager.AddAllMethods(obj, options)
	return service
}

// AddMissingMethod is used for publishing a method,
// all methods not explicitly published will be redirected to this method.
func (service *BaseService) AddMissingMethod(method MissingMethod, options MethodOptions) Service {
	service.methodManager.AddMissingMethod(method, options)
	return service
}

// Remove the published func or method by name
func (service *BaseService) Remove(name string) Service {
	service.methodManager.Remove(name)
	return service
}

// Filter return the first filter
func (service *BaseService) Filter() Filter {
	return service.filterManager.Filter()
}

// FilterByIndex return the filter by index
func (service *BaseService) FilterByIndex(index int) Filter {
	return service.filterManager.FilterByIndex(index)
}

// SetFilter will replace the current filter settings
func (service *BaseService) SetFilter(filter ...Filter) Service {
	service.filterManager.SetFilter(filter...)
	return service
}

// AddFilter add the filter to this Service
func (service *BaseService) AddFilter(filter ...Filter) Service {
	service.filterManager.AddFilter(filter...)
	return service
}

// RemoveFilterByIndex remove the filter by the index
func (service *BaseService) RemoveFilterByIndex(index int) Service {
	service.filterManager.RemoveFilterByIndex(index)
	return service
}

// RemoveFilter remove the filter from this Service
func (service *BaseService) RemoveFilter(filter ...Filter) Service {
	service.filterManager.RemoveFilter(filter...)
	return service
}

// AddInvokeHandler add the invoke handler to this Service
func (service *BaseService) AddInvokeHandler(handler InvokeHandler) Service {
	service.handlerManager.AddInvokeHandler(handler)
	return service
}

// AddBeforeFilterHandler add the filter handler before filters
func (service *BaseService) AddBeforeFilterHandler(handler FilterHandler) Service {
	service.handlerManager.AddBeforeFilterHandler(handler)
	return service
}

// AddAfterFilterHandler add the filter handler after filters
func (service *BaseService) AddAfterFilterHandler(handler FilterHandler) Service {
	service.handlerManager.AddAfterFilterHandler(handler)
	return service
}

func (service *BaseService) callService(
	name string, args []reflect.Value, context ServiceContext) []reflect.Value {
	remoteMethod := context.Method
	function := remoteMethod.Function
	if context.IsMissingMethod {
		missingMethod := function.Interface().(MissingMethod)
		return missingMethod(name, args, context)
	}
	ft := function.Type()
	n := len(args)
	if ft.IsVariadic() {
		return function.CallSlice(args)
	}
	if ft.NumIn() == n+1 {
		args = service.FixArguments(args, ft.In(n), context)
	}
	return function.Call(args)
}

func (service *BaseService) getErrorMessage(err error) string {
	if panicError, ok := err.(*promise.PanicError); ok {
		if service.Debug {
			return fmt.Sprintf("%v\r\n%s", panicError.Panic, panicError.Stack)
		}
		return panicError.Error()
	}
	return err.Error()
}

func (service *BaseService) fireErrorEvent(err error, context Context) error {
	defer func() {
		if e := recover(); e != nil {
			err = promise.NewPanicError(e)
		}
	}()
	switch event := service.Event.(type) {
	case sendErrorEvent:
		event.OnSendError(err, context)
	case sendErrorEvent2:
		err = event.OnSendError(err, context)
	}
	return err
}

func (service *BaseService) sendError(err error, context Context) []byte {
	err = service.fireErrorEvent(err, context)
	w := io.NewWriter(true)
	w.WriteByte(io.TagError)
	w.WriteString(service.getErrorMessage(err))
	return w.Bytes()
}

func (service *BaseService) endError(err error, context Context) []byte {
	w := io.NewByteWriter(service.sendError(err, context))
	w.WriteByte(io.TagEnd)
	return w.Bytes()
}

func doOutput(
	args []reflect.Value,
	results []reflect.Value,
	context ServiceContext) []byte {
	method := context.Method
	w := io.NewWriter(method.Simple)
	switch method.Mode {
	case RawWithEndTag:
		return results[0].Bytes()
	case Raw:
		w.Write(results[0].Bytes())
	default:
		w.WriteByte(io.TagResult)
		if method.Mode == Serialized {
			w.Write(results[0].Bytes())
		} else {
			switch len(results) {
			case 0:
				w.WriteNil()
			case 1:
				w.WriteValue(results[0])
			default:
				w.WriteSlice(results)
			}
		}
		if context.ByRef {
			w.WriteByte(io.TagArgument)
			w.Reset()
			w.WriteSlice(args)
		}
	}
	return w.Bytes()
}

func (service *BaseService) fireBeforeInvokeEvent(
	name string, args []reflect.Value, context ServiceContext) error {
	switch event := service.Event.(type) {
	case beforeInvokeEvent:
		event.OnBeforeInvoke(name, args, context.ByRef, context)
	case beforeInvokeEvent2:
		return event.OnBeforeInvoke(name, args, context.ByRef, context)
	}
	return nil
}

func (service *BaseService) fireAfterInvokeEvent(
	name string, args []reflect.Value, results []reflect.Value, context ServiceContext) error {
	switch event := service.Event.(type) {
	case afterInvokeEvent:
		event.OnAfterInvoke(name, args, context.ByRef, results, context)
	case afterInvokeEvent2:
		return event.OnAfterInvoke(name, args, context.ByRef, results, context)
	}
	return nil
}

func (service *BaseService) beforeInvoke(
	name string, args []reflect.Value, context ServiceContext) promise.Promise {
	return promise.Sync(func() error {
		return service.fireBeforeInvokeEvent(name, args, context)
	}).Then(func() promise.Promise {
		return service.handlerManager.invokeHandler(name, args, context)
	}).Then(func(results []reflect.Value) ([]byte, error) {
		err := service.fireAfterInvokeEvent(name, args, results, context)
		if err != nil {
			return nil, err
		}
		return doOutput(args, results, context), nil
	})
}

func (service *BaseService) invokeHandler(
	name string, args []reflect.Value, context ServiceContext) promise.Promise {
	if context.Oneway {
		go func() {
			defer recover()
			service.callService(name, args, context)
		}()
		return promise.Resolve(nil)
	}
	return promise.Sync(func() ([]reflect.Value, error) {
		results := service.callService(name, args, context)
		n := len(results)
		if n == 0 {
			return results, nil
		}
		err, ok := results[n-1].Interface().(error)
		if ok {
			results = results[:n-1]
		}
		return results, err
	})
}
