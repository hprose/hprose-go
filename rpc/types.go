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
 * rpc/types.go                                           *
 *                                                        *
 * reflect types for Go.                                  *
 *                                                        *
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"net/http"
	"unsafe"
)

type emptyInterface struct {
	typ uintptr
	ptr uintptr
}

func getType(v interface{}) uintptr {
	return *(*uintptr)(unsafe.Pointer(&v))
}

var interfaceType = getType((interface{})(nil))
var contextType = getType((Context)(nil))
var serviceContextType = getType((*ServiceContext)(nil))
var httpContextType = getType((*HTTPContext)(nil))
var httpRequestType = getType((*http.Request)(nil))
