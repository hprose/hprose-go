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
 * LastModified: Sep 23, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import "time"

// Service interface
type Service interface {
	AddFunction(name string, function interface{}, options Options) Service
	AddFunctions(names []string, functions []interface{}, options Options) Service
	AddMethod(name string, obj interface{}, options Options, alias ...string) Service
	AddMethods(names []string, obj interface{}, options Options, aliases ...[]string) Service
	AddInstanceMethods(obj interface{}, options Options) Service
	AddAllMethods(obj interface{}, options Options) Service
	AddMissingMethod(method MissingMethod, options Options) Service
	AddNetRPCMethods(rcvr interface{}, options Options) Service
	Remove(name string) Service
	Filter() Filter
	FilterByIndex(index int) Filter
	SetFilter(filter ...Filter) Service
	AddFilter(filter ...Filter) Service
	RemoveFilterByIndex(index int) Service
	RemoveFilter(filter ...Filter) Service
	AddInvokeHandler(handler ...InvokeHandler) Service
	AddBeforeFilterHandler(handler ...FilterHandler) Service
	AddAfterFilterHandler(handler ...FilterHandler) Service
	Publish(topic string, timeout time.Duration, heartbeat time.Duration) Service
	Clients
}
