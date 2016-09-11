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
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import "time"

// Service interface
type Service interface {
	AddFunction(name string, function interface{}, options MethodOptions) Service
	AddFunctions(names []string, functions []interface{}, options MethodOptions) Service
	AddMethod(name string, obj interface{}, options MethodOptions, alias ...string) Service
	AddMethods(names []string, obj interface{}, options MethodOptions, aliases ...[]string) Service
	AddInstanceMethods(obj interface{}, options MethodOptions) Service
	AddAllMethods(obj interface{}, options MethodOptions) Service
	AddMissingMethod(method MissingMethod, options MethodOptions) Service
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

// BaseService is the hprose base service
type BaseService struct {
	ServiceEvent
	Debug          bool
	Simple         bool
	Timeout        time.Duration
	Heartbeat      time.Duration
	ErrorDelay     time.Duration
	methods        *serviceMethods
	handlerManager *handlerManager
	filterManager  *filterManager
	allTopics      map[string]map[string]*topic
	argsFixer
}

// NewBaseService is the constructor for BaseService
func NewBaseService() (service *BaseService) {
	service = new(BaseService)
	service.methods = newServiceMethods()
	service.handlerManager = newHandlerManager()
	service.filterManager = &filterManager{}
	service.Timeout = 120 * 1000 * 1000
	service.Heartbeat = 3 * 1000 * 1000
	service.ErrorDelay = 10 * 1000 * 1000
	service.allTopics = make(map[string]map[string]*topic)
	return service
}

// AddFunction publish a func or bound method
// name is the method name
// function is a func or bound method
// options includes Mode, Simple, Oneway and NameSpace
func (service *BaseService) AddFunction(name string, function interface{}, options MethodOptions) Service {
	service.methods.AddFunction(name, function, options)
	return service
}

// AddFunctions is used for batch publishing service method
func (service *BaseService) AddFunctions(names []string, functions []interface{}, options MethodOptions) Service {
	service.methods.AddFunctions(names, functions, options)
	return service
}

// AddMethod is used for publishing a method on the obj with an alias
func (service *BaseService) AddMethod(name string, obj interface{}, options MethodOptions, alias ...string) Service {
	service.methods.AddMethod(name, obj, options, alias...)
	return service
}

// AddMethods is used for batch publishing methods on the obj with aliases
func (service *BaseService) AddMethods(names []string, obj interface{}, options MethodOptions, aliases ...[]string) Service {
	service.methods.AddMethods(names, obj, options, aliases...)
	return service
}

// AddInstanceMethods is used for publishing all the public methods and func fields with options.
func (service *BaseService) AddInstanceMethods(obj interface{}, options MethodOptions) Service {
	service.methods.AddInstanceMethods(obj, options)
	return service
}

// AddAllMethods will publish all methods and non-nil function fields on the
// obj self and on its anonymous or non-anonymous struct fields (or pointer to
// pointer ... to pointer struct fields). This is a recursive operation.
// So it's a pit, if you do not know what you are doing, do not step on.
func (service *BaseService) AddAllMethods(obj interface{}, options MethodOptions) Service {
	service.methods.AddAllMethods(obj, options)
	return service
}

// AddMissingMethod is used for publishing a method,
// all methods not explicitly published will be redirected to this method.
func (service *BaseService) AddMissingMethod(method MissingMethod, options MethodOptions) Service {
	service.methods.AddMissingMethod(method, options)
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
