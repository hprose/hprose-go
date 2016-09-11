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
 * rpc/service_event.go                                   *
 *                                                        *
 * hprose service event for Go.                           *
 *                                                        *
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import "reflect"

// ServiceEvent is the service event
type ServiceEvent interface{}

type beforeInvokeEvent interface {
	OnBeforeInvoke(name string, args []reflect.Value, byref bool, context Context)
}

type beforeInvoke2Event interface {
	OnBeforeInvoke(name string, args []reflect.Value, byref bool, context Context) error
}

type afterInvokeEvent interface {
	OnAfterInvoke(name string, args []reflect.Value, byref bool, result []reflect.Value, context Context)
}

type afterInvoke2Event interface {
	OnAfterInvoke(name string, args []reflect.Value, byref bool, result []reflect.Value, context Context) error
}

type sendErrorEvent interface {
	OnSendError(err error, context Context)
}

type sendError2Event interface {
	OnSendError(err error, context Context) error
}
