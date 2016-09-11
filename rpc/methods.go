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
 * rpc/methods.go                                         *
 *                                                        *
 * hprose methods for Go.                                 *
 *                                                        *
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"reflect"
	"strings"
)

// MethodOptions is the options of the published service method
type MethodOptions struct {
	Mode      ResultMode
	Simple    bool
	Oneway    bool
	NameSpace string
}

// Method is the published service method
type Method struct {
	Function reflect.Value
	MethodOptions
}

// Methods is the published service methods
type Methods struct {
	MethodNames   []string
	RemoteMethods map[string]Method
}

// NewMethods is the constructor for Methods
func NewMethods() (methods *Methods) {
	methods = new(Methods)
	methods.MethodNames = make([]string, 0, 64)
	methods.RemoteMethods = make(map[string]Method)
	return
}

// AddFunction publish a func or bound method
// name is the method name
// function is a func or bound method
// options includes Mode, Simple, Oneway and NameSpace
func (methods *Methods) AddFunction(
	name string, function interface{}, options MethodOptions) {
	if name == "" {
		panic("name can't be empty")
	}
	if function == nil {
		panic("function can't be nil")
	}
	f, ok := function.(reflect.Value)
	if !ok {
		f = reflect.ValueOf(function)
	}
	if f.Kind() != reflect.Func {
		panic("function must be func or bound method")
	}
	if options.NameSpace != "" && name != "*" {
		name = options.NameSpace + "_" + name
	}
	methods.MethodNames = append(methods.MethodNames, name)
	methods.RemoteMethods[strings.ToLower(name)] = Method{f, options}
}

// AddFunctions is used for batch publishing service method
func (methods *Methods) AddFunctions(
	names []string, functions []interface{}, options MethodOptions) {
	count := len(names)
	if count != len(functions) {
		panic("names and functions must have the same length")
	}
	for i := 0; i < count; i++ {
		methods.AddFunction(names[i], functions[i], options)
	}
}

// AddMethod is used for publishing a method on the obj with an alias
func (methods *Methods) AddMethod(
	name string, obj interface{}, options MethodOptions, alias ...string) {
	if obj == nil {
		panic("obj can't be nil")
	}
	f := reflect.ValueOf(obj).MethodByName(name)
	if len(alias) == 1 && alias[0] != "" {
		name = alias[0]
	}
	if f.CanInterface() {
		methods.AddFunction(name, f, options)
	}
}

// AddMethods is used for batch publishing methods on the obj with aliases
func (methods *Methods) AddMethods(
	names []string, obj interface{}, options MethodOptions, aliases ...[]string) {
	if obj == nil {
		panic("obj can't be nil")
	}
	count := len(names)
	if len(aliases) == 1 {
		if len(aliases[0]) != count {
			panic("names and aliases must have the same length")
		}
		for i := 0; i < count; i++ {
			methods.AddMethod(names[i], obj, options, aliases[0][i])
		}
		return
	}
	for i := 0; i < count; i++ {
		methods.AddMethod(names[i], obj, options)
	}
}

func (methods *Methods) addInstanceMethods(
	v reflect.Value, t reflect.Type, options MethodOptions) {
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		name := t.Method(i).Name
		method := v.Method(i)
		if method.CanInterface() {
			methods.AddFunction(name, method, options)
		}
	}
}

func getPtrTo(v reflect.Value, t reflect.Type) (reflect.Value, reflect.Type) {
	for ; t.Kind() == reflect.Ptr && !v.IsNil(); v = v.Elem() {
		t = t.Elem()
	}
	return v, t
}

// AddInstanceMethods is used for publishing all the public methods and func fields with options.
func (methods *Methods) AddInstanceMethods(
	obj interface{}, options MethodOptions) {
	if obj == nil {
		panic("obj can't be nil")
	}
	v := reflect.ValueOf(obj)
	t := v.Type()
	methods.addInstanceMethods(v, t, options)
	v, t = getPtrTo(v, t)
	if t.Kind() == reflect.Struct {
		n := t.NumField()
		for i := 0; i < n; i++ {
			f := v.Field(i)
			name := t.Field(i).Name
			if f.CanInterface() && f.IsValid() {
				f, _ = getPtrTo(f, f.Type())
				if !f.IsNil() && f.Kind() == reflect.Func {
					methods.AddFunction(name, f, options)
				}
			}
		}
	}
}

// AddAllMethods will publish all methods and non-nil function fields on the
// obj self and on its anonymous or non-anonymous struct fields (or pointer to
// pointer ... to pointer struct fields). This is a recursive operation.
// So it's a pit, if you do not know what you are doing, do not step on.
func (methods *Methods) AddAllMethods(
	obj interface{}, options MethodOptions) {
	if obj == nil {
		panic("obj can't be nil")
	}
	v := reflect.ValueOf(obj)
	t := v.Type()
	methods.addInstanceMethods(v, t, options)
	v, t = getPtrTo(v, t)
	if t.Kind() == reflect.Struct {
		n := t.NumField()
		for i := 0; i < n; i++ {
			f := v.Field(i)
			fs := t.Field(i)
			name := fs.Name
			if f.CanInterface() && f.IsValid() {
				f, _ = getPtrTo(f, f.Type())
				if !f.IsNil() && f.Kind() == reflect.Func {
					methods.AddFunction(name, f, options)
				} else if f.Kind() == reflect.Struct {
					if fs.Anonymous {
						methods.AddAllMethods(f.Interface(), options)
					} else {
						newOptions := options
						if newOptions.NameSpace == "" {
							newOptions.NameSpace = name
						} else {
							newOptions.NameSpace += "_" + name
						}
						methods.AddAllMethods(f.Interface(), newOptions)
					}
				}
			}
		}
	}
}

// MissingMethod is missing method
type MissingMethod func(name string, args []reflect.Value) (result []reflect.Value)

// AddMissingMethod is used for publishing a method,
// all methods not explicitly published will be redirected to this method.
func (methods *Methods) AddMissingMethod(
	method MissingMethod, options MethodOptions) {
	methods.AddFunction("*", method, options)
}
