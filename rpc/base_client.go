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
 * LastModified: Oct 10, 2016                             *
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
	"sync"
	"sync/atomic"
	"time"

	hio "github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/util"
)

type clientTopic struct {
	callbacks []Callback
	locker    sync.RWMutex
}

func (ct *clientTopic) addCallback(callback Callback) {
	ct.locker.Lock()
	ct.callbacks = append(ct.callbacks, callback)
	ct.locker.Unlock()
}

type topicManager struct {
	allTopics map[string]map[string]*clientTopic
	locker    sync.RWMutex
}

func (tm *topicManager) getTopic(
	name string, id string, create bool) *clientTopic {
	tm.locker.RLock()
	topics := tm.allTopics[name]
	if topics != nil {
		topic := topics[id]
		tm.locker.RUnlock()
		return topic
	}
	tm.locker.RUnlock()
	if create {
		tm.locker.Lock()
		tm.allTopics[name] = make(map[string]*clientTopic)
		tm.locker.Unlock()
	}
	return nil
}

// BaseClient is the hprose base client
type BaseClient struct {
	handlerManager
	filterManager
	topicManager
	uri            string
	uriList        []string
	index          int32
	failround      int
	retry          int
	timeout        time.Duration
	event          ClientEvent
	contextPool    sync.Pool
	readerPool     sync.Pool
	SendAndReceive func([]byte, *ClientContext) ([]byte, error)
	id             string
}

func (client *BaseClient) initBaseClient() {
	client.initHandlerManager()
	client.timeout = 30 * 1000 * 1000 * 1000
	client.retry = 10
	client.contextPool = sync.Pool{
		New: func() interface{} { return new(ClientContext) },
	}
	client.readerPool = sync.Pool{
		New: func() interface{} { return new(hio.Reader) },
	}
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
	client.allTopics = make(map[string]map[string]*clientTopic)
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

func shuffleStringSlice(src []string) []string {
	dest := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
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

func (client *BaseClient) acquireContext() (context *ClientContext) {
	return client.contextPool.Get().(*ClientContext)
}

func (client *BaseClient) releaseContext(context *ClientContext) {
	client.contextPool.Put(context)
}

func (client *BaseClient) initClientContext(
	context *ClientContext, settings *InvokeSettings) {
	context.initBaseContext()
	context.Client = client
	context.Retried = 0
	if settings == nil {
		context.InvokeSettings = InvokeSettings{
			Timeout: client.timeout,
			Retry:   client.retry,
		}
	} else {
		context.InvokeSettings = *settings
		if settings.Timeout <= 0 {
			context.Timeout = client.timeout
		}
		if settings.Retry <= 0 {
			context.Retry = client.retry
		}
	}
}

// Invoke the remote method synchronous
func (client *BaseClient) Invoke(name string, args []reflect.Value, settings *InvokeSettings) (results []reflect.Value, err error) {
	context := client.acquireContext()
	client.initClientContext(context, settings)
	results, err = client.handlerManager.invokeHandler(name, args, context)
	if results == nil && len(context.ResultTypes) > 0 {
		n := len(context.ResultTypes)
		results = make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			results[i] = reflect.New(context.ResultTypes[i]).Elem()
		}
	}
	client.releaseContext(context)
	return
}

// Go invoke the remote method asynchronous
func (client *BaseClient) Go(name string, args []reflect.Value, callback Callback, settings *InvokeSettings) {
	go func() {
		defer client.fireErrorEvent(name, nil)
		callback(client.Invoke(name, args, settings))
	}()
}

// Close the client
func (client *BaseClient) Close() {}

func (client *BaseClient) fireErrorEvent(name string, err error) {
	if e := recover(); e != nil {
		err = NewPanicError(e)
	}
	if err != nil {
		if event, ok := client.event.(onErrorEvent); ok {
			event.OnError(name, err)
		}
	}
}

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

func readMultiResults(
	reader *hio.Reader, resultTypes []reflect.Type) (results []reflect.Value) {
	length := len(resultTypes)
	reader.CheckTag(hio.TagList)
	count := reader.ReadCount()
	results = make([]reflect.Value, util.Max(length, count))
	for i := 0; i < length; i++ {
		results[i] = reflect.New(resultTypes[i]).Elem()
	}
	if length < count {
		for i := length; i < count; i++ {
			results[i] = reflect.New(interfaceType).Elem()
		}
	}
	reader.ReadSlice(results[:count])
	return
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
		results = readMultiResults(reader, context.ResultTypes)
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

func (client *BaseClient) acquireReader(buf []byte) (reader *hio.Reader) {
	reader = client.readerPool.Get().(*hio.Reader)
	reader.Init(buf)
	return
}

func (client *BaseClient) releaseReader(reader *hio.Reader) {
	reader.Init(nil)
	reader.Reset()
	client.readerPool.Put(reader)
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
	reader := client.acquireReader(data)
	defer client.releaseReader(reader)
	reader.JSONCompatible = context.JSONCompatible
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

func getIn(in []reflect.Value) []reflect.Value {
	inlen := len(in)
	varlen := 0
	argc := inlen - 1
	varlen = in[argc].Len()
	argc += varlen
	args := make([]reflect.Value, argc)
	if argc > 0 {
		for i := 0; i < inlen-1; i++ {
			args[i] = in[i]
		}
		v := in[inlen-1]
		for i := 0; i < varlen; i++ {
			args[inlen-1+i] = v.Index(i)
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
		if isVariadic {
			in = getIn(in)
		}
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
			if isVariadic {
				in = getIn(in)
			}
			callback := in[0]
			in = in[1:]
			out, err := client.Invoke(name, in, settings)
			if hasError {
				out = append(out, reflect.ValueOf(&err).Elem())
			}
			defer client.fireErrorEvent(name, err)
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
		ByRef:          getBoolValue(sf.Tag, "byref"),
		Simple:         getBoolValue(sf.Tag, "simple"),
		Idempotent:     getBoolValue(sf.Tag, "idempotent"),
		Failswitch:     getBoolValue(sf.Tag, "failswitch"),
		Oneway:         getBoolValue(sf.Tag, "oneway"),
		JSONCompatible: getBoolValue(sf.Tag, "jsoncompat"),
		Retry:          int(getInt64Value(sf.Tag, "retry")),
		Mode:           getResultMode(sf.Tag),
		Timeout:        time.Duration(getInt64Value(sf.Tag, "timeout")),
		ResultTypes:    outTypes,
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

var autoIDSettings = InvokeSettings{
	Simple:      true,
	Idempotent:  true,
	Failswitch:  true,
	ResultTypes: []reflect.Type{stringType},
}

// ID returns the auto id of this hprose client
func (client *BaseClient) ID() (string, error) {
	client.topicManager.locker.RLock()
	if client.id != "" {
		client.topicManager.locker.RUnlock()
		return client.id, nil
	}
	client.topicManager.locker.RUnlock()
	client.topicManager.locker.Lock()
	defer client.topicManager.locker.Unlock()
	if client.id != "" {
		return client.id, nil
	}
	results, err := client.Invoke("#", nil, &autoIDSettings)
	if err != nil {
		return "", err
	}
	client.id = results[0].String()
	return client.id, nil
}

func (client *BaseClient) subscribe(
	name string, id string, settings *InvokeSettings) {
	resultTypes := settings.ResultTypes
	settings.ResultTypes = []reflect.Type{interfaceType}
	args := []reflect.Value{reflect.ValueOf(id)}
	for {
		topic := client.getTopic(name, id, false)
		if topic == nil {
			return
		}
		results, err := client.Invoke(name, args, settings)
		if !results[0].IsNil() {
			func() {
				defer client.fireErrorEvent(name, nil)
				if resultTypes != nil && len(resultTypes) > 0 {
					writer := hio.NewWriter(false)
					writer.WriteValue(results[0])
					reader := client.acquireReader(writer.Bytes())
					if len(resultTypes) == 1 {
						results = make([]reflect.Value, 1)
						results[0] = reflect.New(resultTypes[0]).Elem()
						reader.ReadValue(results[0])
					} else {
						results = readMultiResults(reader, resultTypes)
					}
					client.releaseReader(reader)
				}
				topic.locker.RLock()
				callbacks := topic.callbacks
				topic.locker.RUnlock()
				for _, callback := range callbacks {
					callback(results, err)
				}
			}()
		}
	}
}

// Subscribe a push topic
func (client *BaseClient) Subscribe(
	name string, id string,
	settings *InvokeSettings, callback interface{}) (err error) {
	if id == "" {
		id, err = client.ID()
		if err != nil {
			return err
		}
	}
	f := reflect.ValueOf(callback)
	if f.Kind() != reflect.Func {
		return errors.New("Subscribe: callback must be a function")
	}
	resultTypes, hasError := getCallbackResultTypes(f.Type())
	cb := func(results []reflect.Value, err error) {
		if hasError {
			results = append(results, reflect.ValueOf(&err).Elem())
		}
		f.Call(results)
	}
	if settings == nil {
		settings = new(InvokeSettings)
	}
	settings.ByRef = false
	settings.Idempotent = true
	settings.Mode = Normal
	settings.Oneway = false
	settings.Simple = true
	settings.ResultTypes = resultTypes
	topic := client.getTopic(name, id, true)
	if topic == nil {
		topic = new(clientTopic)
		topic.addCallback(cb)
		client.topicManager.locker.Lock()
		client.allTopics[name][id] = topic
		client.topicManager.locker.Unlock()
		go client.subscribe(name, id, settings)
	} else {
		topic.addCallback(cb)
	}
	return nil
}

// Unsubscribe a push topic
func (client *BaseClient) Unsubscribe(name string, id ...string) {
	client.topicManager.locker.Lock()
	if client.allTopics[name] != nil {
		if len(id) == 0 {
			if client.id == "" {
				delete(client.allTopics, name)
			} else {
				delete(client.allTopics[name], client.id)
			}
		} else {
			for i := range id {
				delete(client.allTopics[name], id[i])
			}
		}
	}
	client.topicManager.locker.Unlock()
}
