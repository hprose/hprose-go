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
 * rpc/args_fixer.go                                      *
 *                                                        *
 * hprose args fixer for Go.                              *
 *                                                        *
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"reflect"
	"unsafe"
)

type argsFixer interface {
	FixArgs(args []reflect.Value, lastParamType reflect.Type, context Context) []reflect.Value
}

func fixArgs(args []reflect.Value, lastParamType reflect.Type, context Context) []reflect.Value {
	typ := (*emptyInterface)(unsafe.Pointer(&lastParamType)).ptr
	if typ == interfaceType || typ == contextType {
		args = append(args, reflect.ValueOf(context))
	}
	return args
}
