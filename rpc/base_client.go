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
 * rpc/base_client.go                                     *
 *                                                        *
 * hprose rpc base client for Go.                         *
 *                                                        *
 * LastModified: Sep 27, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	hio "github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/util"
)

// BaseClient is the hprose base client
type BaseClient struct {
	*handlerManager
	*filterManager
	uri            string
	uriList        []string
	index          int32
	failround      int
	retry          int
	timeout        time.Duration
	event          ClientEvent
	SendAndReceive func([]byte, *ClientContext) ([]byte, error)
}

// NewBaseClient is the constructor for BaseClient
func NewBaseClient() (client *BaseClient) {
	client = new(BaseClient)
	client.handlerManager = newHandlerManager()
	client.filterManager = &filterManager{}
	client.timeout = 30 * 1000 * 1000 * 1000
	client.retry = 10
	client.override.invokeHandler = func(
		name string, args []reflect.Value,
		context Context) (results []reflect.Value, err error) {
		return client.invoke(name, args, context.(*ClientContext))
	}
	client.override.beforeFilterHandler = func(
		request []byte, context Context) (response []byte, err error) {
		return client.beforeFilter(request, context.(*ClientContext))
	}
	client.override.afterFilterHandler = func(
		request []byte, context Context) (response []byte, err error) {
		return client.afterFilter(request, context.(*ClientContext))
	}
	return client
}

func shuffleStringSlice(src []string) []string {
	dest := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

// URI returns the current hprose service address.
func (client *BaseClient) URI() string {
	return client.uri
}

// SetURI set the current hprose service address.
//
// If you want to set more than one service address, please don't use this
// method, use SetURIList instead.
func (client *BaseClient) SetURI(uri string) {
	client.SetURIList([]string{uri})
}

// URIList returns all of the hprose service addresses
func (client *BaseClient) URIList() []string {
	return client.uriList
}

// SetURIList set a list of server addresses
func (client *BaseClient) SetURIList(uriList []string) {
	client.uriList = shuffleStringSlice(uriList)
	client.index = 0
	client.failround = 0
	client.uri = client.uriList[0]
}

// TLSClientConfig returns the tls config of hprose client
func (client *BaseClient) TLSClientConfig() *tls.Config {
	return nil
}

// SetTLSClientConfig set the tls config of hprose client
func (client *BaseClient) SetTLSClientConfig(config *tls.Config) {}

// Retry returns the max retry count
func (client *BaseClient) Retry() int {
	return client.retry
}

// SetRetry set the max retry count
func (client *BaseClient) SetRetry(value int) {
	client.retry = value
}

// Timeout returns the client timeout setting
func (client *BaseClient) Timeout() time.Duration {
	return client.timeout
}

// SetTimeout set the client timeout setting
func (client *BaseClient) SetTimeout(value time.Duration) {
	client.timeout = value
}

// Failround return the fail round
func (client *BaseClient) Failround() int {
	return client.failround
}

// SetEvent set the client event
func (client *BaseClient) SetEvent(event ClientEvent) {
	client.event = event
}

// Filter return the first filter
func (client *BaseClient) Filter() Filter {
	return client.filterManager.Filter()
}

// FilterByIndex return the filter by index
func (client *BaseClient) FilterByIndex(index int) Filter {
	return client.filterManager.FilterByIndex(index)
}

// SetFilter will replace the current filter settings
func (client *BaseClient) SetFilter(filter ...Filter) Client {
	client.filterManager.SetFilter(filter...)
	return client
}

// AddFilter add the filter to this Service
func (client *BaseClient) AddFilter(filter ...Filter) Client {
	client.filterManager.AddFilter(filter...)
	return client
}

// RemoveFilterByIndex remove the filter by the index
func (client *BaseClient) RemoveFilterByIndex(index int) Client {
	client.filterManager.RemoveFilterByIndex(index)
	return client
}

// RemoveFilter remove the filter from this Service
func (client *BaseClient) RemoveFilter(filter ...Filter) Client {
	client.filterManager.RemoveFilter(filter...)
	return client
}

// AddInvokeHandler add the invoke handler to this Service
func (client *BaseClient) AddInvokeHandler(handler ...InvokeHandler) Client {
	client.handlerManager.AddInvokeHandler(handler...)
	return client
}

// AddBeforeFilterHandler add the filter handler before filters
func (client *BaseClient) AddBeforeFilterHandler(handler ...FilterHandler) Client {
	client.handlerManager.AddBeforeFilterHandler(handler...)
	return client
}

// AddAfterFilterHandler add the filter handler after filters
func (client *BaseClient) AddAfterFilterHandler(handler ...FilterHandler) Client {
	client.handlerManager.AddAfterFilterHandler(handler...)
	return client
}

// UseService build a remote service proxy object with namespace
func (client *BaseClient) UseService(remoteService interface{}, namespace ...string) {
	ns := ""
	if len(namespace) == 1 {
		ns = namespace[0]
	}
	v := reflect.ValueOf(remoteService)
	if v.Kind() != reflect.Ptr {
		panic("UseService: remoteService argument must be a pointer")
	}
	buildRemoteService(client, v, ns)
}

// Invoke the remote method synchronous
func (client *BaseClient) Invoke(name string, args []reflect.Value, settings *InvokeSettings) (results []reflect.Value, err error) {
	context := client.getContext(settings)
	results, err = client.handlerManager.invokeHandler(name, args, context)
	if results == nil && len(context.ResultTypes) > 0 {
		n := len(context.ResultTypes)
		results = make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			results[i] = reflect.New(context.ResultTypes[i]).Elem()
		}
	}
	return
}

// Go invoke the remote method asynchronous
func (client *BaseClient) Go(name string, args []reflect.Value, callback Callback, settings *InvokeSettings) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				err := NewPanicError(e)
				if event, ok := client.event.(onErrorEvent); ok {
					event.OnError(name, err)
				}
			}
		}()
		callback(client.Invoke(name, args, settings))
	}()
}

// Close the client
func (client *BaseClient) Close() {}

func (client *BaseClient) beforeFilter(
	request []byte,
	context *ClientContext) (response []byte, err error) {
	request = client.outputFilter(request, context)
	if context.Oneway {
		go client.handlerManager.afterFilterHandler(request, context)
		return nil, nil
	}
	response, err = client.handlerManager.afterFilterHandler(request, context)
	response = client.inputFilter(response, context)
	return
}

func (client *BaseClient) afterFilter(
	request []byte, context Context) (response []byte, err error) {
	return client.SendAndReceive(request, context.(*ClientContext))
}

func (client *BaseClient) sendRequest(
	request []byte,
	context *ClientContext) (response []byte, err error) {
	response, err = client.handlerManager.beforeFilterHandler(request, context)
	if err != nil {
		response, err = client.retrySendReqeust(request, err, context)
	}
	return
}

func (client *BaseClient) retrySendReqeust(
	request []byte,
	err error,
	context *ClientContext) ([]byte, error) {
	if context.Failswitch {
		client.failswitch()
	}
	if context.Idempotent && context.Retried < context.Retry {
		context.Retried++
		interval := context.Retried * 500
		if context.Failswitch {
			interval -= (len(client.uriList) - 1) * 500
		}
		if interval > 5000 {
			interval = 5000
		}
		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
		return client.sendRequest(request, context)
	}
	return nil, err
}

func (client *BaseClient) failswitch() {
	n := int32(len(client.uriList))
	if n > 1 {
		if atomic.CompareAndSwapInt32(&client.index, n-1, 0) {
			client.uri = client.uriList[0]
			client.failround++
		} else {
			client.uri = client.uriList[atomic.AddInt32(&client.index, 1)]
		}
	} else {
		client.failround++
	}
	if event, ok := client.event.(onFailswitchEvent); ok {
		event.OnFailswitch(client)
	}
}

func (client *BaseClient) getContext(settings *InvokeSettings) *ClientContext {
	context := new(ClientContext)
	context.BaseContext = NewBaseContext()
	context.Client = client
	if settings == nil {
		context.Timeout = client.timeout
		context.Retry = client.retry
	} else {
		context.InvokeSettings = *settings
		if settings.Timeout <= 0 {
			context.Timeout = client.timeout
		}
		if settings.Retry <= 0 {
			context.Retry = client.retry
		}
	}
	return context
}

func (client *BaseClient) encode(
	name string,
	args []reflect.Value,
	context *ClientContext) []byte {
	writer := hio.NewWriter(context.Simple)
	writer.WriteByte(hio.TagCall)
	writer.WriteString(name)
	if len(args) > 0 || context.ByRef {
		writer.Reset()
		writer.WriteSlice(args)
		if context.ByRef {
			writer.WriteBool(true)
		}
	}
	writer.WriteByte(hio.TagEnd)
	return writer.Bytes()
}

func (client *BaseClient) readResults(
	reader *hio.Reader,
	context *ClientContext) (results []reflect.Value) {
	length := len(context.ResultTypes)
	switch length {
	case 0:
		var e interface{}
		reader.Unserialize(&e)
	case 1:
		results = make([]reflect.Value, 1)
		results[0] = reflect.New(context.ResultTypes[0]).Elem()
		reader.ReadValue(results[0])
	default:
		reader.CheckTag(hio.TagList)
		count := reader.ReadCount()
		results = make([]reflect.Value, util.Max(length, count))
		for i := 0; i < length; i++ {
			results[i] = reflect.New(context.ResultTypes[i]).Elem()
		}
		if length < count {
			for i := length; i < count; i++ {
				results[i] = reflect.New(interfaceType).Elem()
			}
		}
		reader.ReadSlice(results[:count])
	}
	return
}

func (client *BaseClient) readArguments(
	reader *hio.Reader,
	args []reflect.Value,
	context *ClientContext) {
	length := len(args)
	reader.CheckTag(hio.TagList)
	count := reader.ReadCount()
	a := make([]reflect.Value, util.Max(length, count))
	for i := 0; i < length; i++ {
		a[i] = args[i].Elem()
	}
	if length < count {
		for i := length; i < count; i++ {
			a[i] = reflect.New(interfaceType).Elem()
		}
	}
	reader.ReadSlice(a[:count])
	return
}

func (client *BaseClient) decode(
	data []byte,
	args []reflect.Value,
	context *ClientContext) (results []reflect.Value, err error) {
	if context.Oneway {
		return
	}
	n := len(data)
	if n == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	if data[n-1] != hio.TagEnd {
		return nil, fmt.Errorf("Wrong Response: \r\n%s", data)
	}
	if context.Mode == RawWithEndTag {
		results = make([]reflect.Value, 1)
		results[0] = reflect.ValueOf(data)
		return
	}
	if context.Mode == Raw {
		results = make([]reflect.Value, 1)
		results[0] = reflect.ValueOf(data[:n-1])
		return
	}
	reader := hio.NewReader(data, false)
	tag, _ := reader.ReadByte()
	if tag == hio.TagResult {
		switch context.Mode {
		case Normal:
			results = client.readResults(reader, context)
		case Serialized:
			results = make([]reflect.Value, 1)
			results[0] = reflect.ValueOf(reader.ReadRaw())
		}
		tag, _ = reader.ReadByte()
		if tag == hio.TagArgument {
			reader.Reset()
			client.readArguments(reader, args, context)
			tag, _ = reader.ReadByte()
		}
	} else if tag == hio.TagError {
		return nil, errors.New(reader.ReadString())
	}
	if tag != hio.TagEnd {
		return nil, fmt.Errorf("Wrong Response: \r\n%s", data)
	}
	return
}

func (client *BaseClient) invoke(
	name string,
	args []reflect.Value,
	context *ClientContext) (results []reflect.Value, err error) {
	request := client.encode(name, args, context)
	response, err := client.sendRequest(request, context)
	if err != nil {
		return nil, err
	}
	return client.decode(response, args, context)
}

func buildRemoteService(client *BaseClient, v reflect.Value, ns string) {
	v = v.Elem()
	t := v.Type()
	et := t
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}
	ptr := reflect.New(et)
	obj := ptr.Elem()
	count := obj.NumField()
	for i := 0; i < count; i++ {
		f := obj.Field(i)
		ft := f.Type()
		sf := et.Field(i)
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if f.CanSet() {
			switch ft.Kind() {
			case reflect.Struct:
				buildRemoteSubService(client, f, ft, sf, ns)
			case reflect.Func:
				buildRemoteMethod(client, f, ft, sf, ns)
			}
		}
	}
	if t.Kind() == reflect.Ptr {
		v.Set(ptr)
	} else {
		v.Set(obj)
	}
}

func buildRemoteSubService(client *BaseClient, f reflect.Value, ft reflect.Type,
	sf reflect.StructField, ns string) {
	namespace := ns
	if !sf.Anonymous {
		if ns == "" {
			namespace = sf.Name
		} else {
			namespace += "_" + sf.Name
		}
	}
	fp := reflect.New(ft)
	buildRemoteService(client, fp, namespace)
	if f.Kind() == reflect.Ptr {
		f.Set(fp)
	} else {
		f.Set(fp.Elem())
	}
}

func getRemoteMethodName(sf reflect.StructField, ns string) (name string) {
	name = sf.Tag.Get("name")
	if name == "" {
		name = sf.Name
	}
	if ns != "" {
		name = ns + "_" + name
	}
	return
}

func getBoolValue(tag reflect.StructTag, key string) bool {
	value := tag.Get(key)
	if value == "" {
		return false
	}
	result, _ := strconv.ParseBool(value)
	return result
}

func getResultMode(tag reflect.StructTag) ResultMode {
	value := tag.Get("result")
	switch strings.ToLower(value) {
	case "normal":
		return Normal
	case "serialized":
		return Serialized
	case "raw":
		return Raw
	case "rawwithendtag":
		return RawWithEndTag
	}
	return Normal
}

func getInt64Value(tag reflect.StructTag, key string) int64 {
	value := tag.Get(key)
	if value == "" {
		return 0
	}
	result, _ := strconv.ParseInt(value, 10, 64)
	return result
}

func getResultTypes(ft reflect.Type) ([]reflect.Type, bool) {
	n := ft.NumOut()
	if n == 0 {
		return nil, false
	}
	hasError := (ft.Out(n-1) == errorType)
	if hasError {
		n--
	}
	results := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		results[i] = ft.Out(i)
	}
	return results, hasError
}

func getCallbackResultTypes(ft reflect.Type) ([]reflect.Type, bool) {
	n := ft.NumIn()
	if n == 0 {
		return nil, false
	}
	hasError := (ft.In(n-1) == errorType)
	if hasError {
		n--
	}
	results := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		results[i] = ft.In(i)
	}
	return results, hasError
}

func getIn(in []reflect.Value, isVariadic bool) []reflect.Value {
	inlen := len(in)
	varlen := 0
	argc := inlen
	if isVariadic {
		argc--
		varlen = in[argc].Len()
		argc += varlen
	}
	args := make([]reflect.Value, argc)
	if argc > 0 {
		for i := 0; i < inlen-1; i++ {
			args[i] = in[i]
		}
		if isVariadic {
			v := in[inlen-1]
			for i := 0; i < varlen; i++ {
				args[inlen-1+i] = v.Index(i)
			}
		} else {
			args[inlen-1] = in[inlen-1]
		}
	}
	return args
}

func getSyncRemoteMethod(
	client *BaseClient,
	name string,
	settings *InvokeSettings,
	isVariadic, hasError bool) func(in []reflect.Value) (out []reflect.Value) {
	return func(in []reflect.Value) (out []reflect.Value) {
		in = getIn(in, isVariadic)
		var err error
		out, err = client.Invoke(name, in, settings)
		if hasError {
			out = append(out, reflect.ValueOf(&err).Elem())
		} else if err != nil {
			if e, ok := err.(*PanicError); ok {
				panic(fmt.Sprintf("%v\r\n%s", e.Panic, e.Stack))
			} else {
				panic(err)
			}
		}
		return
	}
}

func getAsyncRemoteMethod(
	client *BaseClient,
	name string,
	settings *InvokeSettings,
	isVariadic, hasError bool) func(in []reflect.Value) (out []reflect.Value) {
	return func(in []reflect.Value) (out []reflect.Value) {
		go func() {
			in = getIn(in, isVariadic)
			callback := in[0]
			in = in[1:]
			out, err := client.Invoke(name, in, settings)
			if hasError {
				out = append(out, reflect.ValueOf(&err).Elem())
			}
			defer func() {
				if e := recover(); e != nil {
					err = NewPanicError(e)
				}
				if err != nil {
					if event, ok := client.event.(onErrorEvent); ok {
						event.OnError(name, err)
					}
				}
			}()
			callback.Call(out)
		}()
		return nil
	}
}

func buildRemoteMethod(client *BaseClient, f reflect.Value, ft reflect.Type, sf reflect.StructField, ns string) {
	name := getRemoteMethodName(sf, ns)
	outTypes, hasError := getResultTypes(ft)
	async := false
	if outTypes == nil && hasError == false {
		if ft.NumIn() > 0 && ft.In(0).Kind() == reflect.Func {
			cbft := ft.In(0)
			if cbft.IsVariadic() {
				panic("callback can't be variadic function")
			}
			outTypes, hasError = getCallbackResultTypes(cbft)
			async = true
		}
	}
	settings := &InvokeSettings{
		ByRef:       getBoolValue(sf.Tag, "byref"),
		Simple:      getBoolValue(sf.Tag, "simple"),
		Mode:        getResultMode(sf.Tag),
		Idempotent:  getBoolValue(sf.Tag, "idempotent"),
		Failswitch:  getBoolValue(sf.Tag, "failswitch"),
		Oneway:      getBoolValue(sf.Tag, "oneway"),
		Retry:       int(getInt64Value(sf.Tag, "retry")),
		Timeout:     time.Duration(getInt64Value(sf.Tag, "timeout")),
		ResultTypes: outTypes,
	}
	var fn func(in []reflect.Value) (out []reflect.Value)
	if async {
		fn = getAsyncRemoteMethod(client, name, settings, ft.IsVariadic(), hasError)
	} else {
		fn = getSyncRemoteMethod(client, name, settings, ft.IsVariadic(), hasError)
	}
	if f.Kind() == reflect.Ptr {
		fp := reflect.New(ft)
		fp.Elem().Set(reflect.MakeFunc(ft, fn))
		f.Set(fp)
	} else {
		f.Set(reflect.MakeFunc(ft, fn))
	}
}
