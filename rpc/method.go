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
 * rpc/method.go                                          *
 *                                                        *
 * hprose method manager for Go.                          *
 *                                                        *
 * LastModified: Sep 20, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"reflect"
	"strings"
	"sync"
)

// Options is the options of the published service method
type Options struct {
	Mode      ResultMode
	Simple    bool
	Oneway    bool
	NameSpace string
}

// Method is the published service method
type Method struct {
	Function reflect.Value
	Options
}

// methodManager manages published service methods
type methodManager struct {
	MethodNames   []string
	RemoteMethods map[string]*Method
	sync.RWMutex
}

func initMethodManager(methods *methodManager) {
	methods.MethodNames = make([]string, 0, 64)
	methods.RemoteMethods = make(map[string]*Method)
}

// AddFunction publish a func or bound method
// name is the method name
// function is a func or bound method
// options includes Mode, Simple, Oneway and NameSpace
func (mm *methodManager) AddFunction(
	name string, function interface{}, options Options) {
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
	mm.Lock()
	mm.MethodNames = append(mm.MethodNames, name)
	mm.RemoteMethods[strings.ToLower(name)] = &Method{f, options}
	mm.Unlock()
}

// AddFunctions is used for batch publishing service method
func (mm *methodManager) AddFunctions(
	names []string, functions []interface{}, options Options) {
	count := len(names)
	if count != len(functions) {
		panic("names and functions must have the same length")
	}
	for i := 0; i < count; i++ {
		mm.AddFunction(names[i], functions[i], options)
	}
}

// AddMethod is used for publishing a method on the obj with an alias
func (mm *methodManager) AddMethod(
	name string, obj interface{}, options Options, alias ...string) {
	if obj == nil {
		panic("obj can't be nil")
	}
	f := reflect.ValueOf(obj).MethodByName(name)
	if len(alias) == 1 && alias[0] != "" {
		name = alias[0]
	}
	if f.CanInterface() {
		mm.AddFunction(name, f, options)
	}
}

// AddMethods is used for batch publishing methods on the obj with aliases
func (mm *methodManager) AddMethods(
	names []string, obj interface{}, options Options, aliases ...[]string) {
	if obj == nil {
		panic("obj can't be nil")
	}
	count := len(names)
	if len(aliases) == 1 {
		if len(aliases[0]) != count {
			panic("names and aliases must have the same length")
		}
		for i := 0; i < count; i++ {
			mm.AddMethod(names[i], obj, options, aliases[0][i])
		}
		return
	}
	for i := 0; i < count; i++ {
		mm.AddMethod(names[i], obj, options)
	}
}

func (mm *methodManager) addMethods(
	v reflect.Value, t reflect.Type, options Options) {
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		name := t.Method(i).Name
		method := v.Method(i)
		if method.CanInterface() {
			mm.AddFunction(name, method, options)
		}
	}
}

func getPtrTo(v reflect.Value, t reflect.Type) (reflect.Value, reflect.Type) {
	for ; t.Kind() == reflect.Ptr && !v.IsNil(); v = v.Elem() {
		t = t.Elem()
	}
	return v, t
}

func (mm *methodManager) addFuncField(
	v reflect.Value, t reflect.Type, i int, options Options) {
	f := v.Field(i)
	name := t.Field(i).Name
	if !f.CanInterface() || !f.IsValid() {
		return
	}
	f, _ = getPtrTo(f, f.Type())
	if !f.IsNil() && f.Kind() == reflect.Func {
		mm.AddFunction(name, f, options)
	}
}

func (mm *methodManager) recursiveAddFuncFields(
	v reflect.Value, t reflect.Type, i int, options Options) {
	f := v.Field(i)
	fs := t.Field(i)
	name := fs.Name
	if !f.CanInterface() || !f.IsValid() {
		return
	}
	f, _ = getPtrTo(f, f.Type())
	if !f.IsNil() && f.Kind() == reflect.Func {
		mm.AddFunction(name, f, options)
		return
	}
	if f.Kind() != reflect.Struct {
		return
	}
	if fs.Anonymous {
		mm.AddAllMethods(f.Interface(), options)
	} else {
		newOptions := options
		if newOptions.NameSpace == "" {
			newOptions.NameSpace = name
		} else {
			newOptions.NameSpace += "_" + name
		}
		mm.AddAllMethods(f.Interface(), newOptions)
	}
}

type addFuncFunc func(
	v reflect.Value,
	t reflect.Type,
	i int,
	options Options)

func (mm *methodManager) addInstanceMethods(
	obj interface{}, options Options, addFunc addFuncFunc) {
	if obj == nil {
		panic("obj can't be nil")
	}
	v := reflect.ValueOf(obj)
	t := v.Type()
	mm.addMethods(v, t, options)
	v, t = getPtrTo(v, t)
	if t.Kind() == reflect.Struct {
		n := t.NumField()
		for i := 0; i < n; i++ {
			addFunc(v, t, i, options)
		}
	}
}

// AddInstanceMethods is used for publishing all the public methods and func fields with options.
func (mm *methodManager) AddInstanceMethods(
	obj interface{}, options Options) {
	mm.addInstanceMethods(obj, options, mm.addFuncField)
}

// AddAllMethods will publish all methods and non-nil function fields on the
// obj self and on its anonymous or non-anonymous struct fields (or pointer to
// pointer ... to pointer struct fields). This is a recursive operation.
// So it's a pit, if you do not know what you are doing, do not step on.
func (mm *methodManager) AddAllMethods(
	obj interface{}, options Options) {
	mm.addInstanceMethods(obj, options, mm.recursiveAddFuncFields)
}

// MissingMethod is missing method
type MissingMethod func(name string, args []reflect.Value, context Context) (result []reflect.Value)

// AddMissingMethod is used for publishing a method,
// all methods not explicitly published will be redirected to this method.
func (mm *methodManager) AddMissingMethod(
	method MissingMethod, options Options) {
	mm.AddFunction("*", method, options)
}

// AddNetRPCMethods is used for publishing methods defined for net/rpc.
//
// For example:
//
//		package server
//
//		type Args struct {
//			A, B int
//		}
//
//		type Quotient struct {
//			Quo, Rem int
//		}
//
//		type Arith int
//
//		func (t *Arith) Multiply(args *Args, reply *int) error {
//			*reply = args.A * args.B
//			return nil
//		}
//
//		func (t *Arith) Divide(args *Args, quo *Quotient) error {
//			if args.B == 0 {
//				return errors.New("divide by zero")
//			}
//			quo.Quo = args.A / args.B
//			quo.Rem = args.A % args.B
//			return nil
//		}
//
// You can publish it in hprose like this:
//
//		service := rpc.NewHTTPService()
//		service.AddNetRPCMethods(new(Arith), rpc.Options{})
//		http.ListenAndServe(":8080", service)
//
//
// Then define stub struct in client:
//
//		type Stub struct {
//			// Synchronous call
//			Multiply func(args *Args) int
//			// Asynchronous call
//			Divide func(args *Args) (<-chan *Quotient, <-chan error)
//		}
//
// Then it can make a remote call:
//
//		client := rpc.NewClient("http://127.0.0.0:8080")
//		var stub *Stub
//		client.UseService(&stub)
//		fmt.Println(stub.Multiply(&Args{8, 7}))
//		quo, e := stub.Device(&Args{8, 7})
//		if err := <-e; err != nil {
//			log.Fatal("arith error:", err)
//		}
//		result := <-quo
//		fmt.Println(result.Quo, result.Rem)
//
// You can also call it in other languages.
//
// Here is a simple example in JavaScript:
//
//		// The method name is case-insensitive
//		var methods = ['multiply', 'divide'];
//		var client = hprose.Client.create('http://127.0.0.1:8080/', methods);
//		// The first letter of the field name is convert to lowercase.
//      client.multiply({a:8, b:7}).then(function(product) {
//			return client.divide({a:product, b:6});
//		}).then(function(quotient) {
//			// The first letter of the field name is convert to lowercase.
//			console.log(quotient.quo, quotient.rem);
//		}).catch(function(err) {
//			// The returned err is here
//			console.error(err);
//		});
//
func (mm *methodManager) AddNetRPCMethods(rcvr interface{}, options Options) {
	if rcvr == nil {
		panic("rcvr can't be nil")
	}
	v := reflect.ValueOf(rcvr)
	t := v.Type()
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		name := t.Method(i).Name
		method := v.Method(i)
		if method.CanInterface() {
			mm.addNetRPCMethod(name, method, options)
		}
	}
}

func (mm *methodManager) addNetRPCMethod(
	name string, method reflect.Value, options Options) {
	ft := method.Type()
	if ft.NumIn() != 2 || ft.IsVariadic() {
		panic("the method " + name + " must has two arguments")
	}
	if ft.In(1).Kind() != reflect.Ptr {
		panic("the second argument of method " + name + " must be a pointer")
	}
	if ft.NumOut() != 1 || ft.Out(0) != errorType {
		panic("the method " + name + " must return type error.")
	}
	argsType := ft.In(0)
	resultType := ft.In(1)
	in := []reflect.Type{argsType}
	out := []reflect.Type{resultType, errorType}
	newft := reflect.FuncOf(in, out, false)
	newMethod := reflect.MakeFunc(newft, func(
		args []reflect.Value) (results []reflect.Value) {
		result := reflect.New(resultType.Elem())
		in := []reflect.Value{args[0], result}
		err := method.Call(in)[0]
		results = []reflect.Value{result, err}
		return
	})
	mm.AddFunction(name, newMethod, options)
}

// Remove the published func or method by name
func (mm *methodManager) Remove(name string) {
	name = strings.ToLower(name)
	mm.Lock()
	n := len(mm.MethodNames)
	for i := 0; i < n; i++ {
		if strings.ToLower(mm.MethodNames[i]) == name {
			mm.MethodNames = append(mm.MethodNames[:i], mm.MethodNames[i+1:]...)
			break
		}
	}
	delete(mm.RemoteMethods, name)
	mm.Unlock()
}
