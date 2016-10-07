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
 * rpc/client.go                                          *
 *                                                        *
 * hprose rpc client for Go.                              *
 *                                                        *
 * LastModified: Oct 2, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"reflect"
	"time"
)

// InvokeSettings is the invoke settings of hprose client
type InvokeSettings struct {
	ByRef          bool
	Simple         bool
	Idempotent     bool
	Failswitch     bool
	Oneway         bool
	JSONCompatible bool
	Retry          int
	Mode           ResultMode
	Timeout        time.Duration
	ResultTypes    []reflect.Type
}

// Callback is the callback function type of Client.Go
type Callback func([]reflect.Value, error)

// Client is hprose client
type Client interface {
	URI() string
	SetURI(uri string)
	URIList() []string
	SetURIList(uriList []string)
	TLSClientConfig() *tls.Config
	SetTLSClientConfig(config *tls.Config)
	Retry() int
	SetRetry(value int)
	Timeout() time.Duration
	SetTimeout(value time.Duration)
	Failround() int
	SetEvent(ClientEvent)
	Filter() Filter
	FilterByIndex(index int) Filter
	SetFilter(filter ...Filter) Client
	AddFilter(filter ...Filter) Client
	RemoveFilterByIndex(index int) Client
	RemoveFilter(filter ...Filter) Client
	AddInvokeHandler(handler ...InvokeHandler) Client
	AddBeforeFilterHandler(handler ...FilterHandler) Client
	AddAfterFilterHandler(handler ...FilterHandler) Client
	UseService(remoteService interface{}, namespace ...string)
	Invoke(string, []reflect.Value, *InvokeSettings) ([]reflect.Value, error)
	Go(string, []reflect.Value, Callback, *InvokeSettings)
	Close()
}

// ClientContext is the hprose client context
type ClientContext struct {
	BaseContext
	InvokeSettings
	Retried int
	Client  Client
}
