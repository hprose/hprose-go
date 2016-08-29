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
 * io/types.go                                            *
 *                                                        *
 * reflect types for Go.                                  *
 *                                                        *
 * LastModified: Aug 27, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"container/list"
	"math/big"
	"time"
	"unsafe"
)

type emptyInterface struct {
	typ uintptr
	ptr uintptr
}

type reflectValue struct {
	typ  uintptr
	ptr  unsafe.Pointer
	flag uintptr
}

func getType(v interface{}) uintptr {
	return *(*uintptr)(unsafe.Pointer(&v))
}

var bigIntType = getType((*big.Int)(nil))
var bigRatType = getType((*big.Rat)(nil))
var bigFloatType = getType((*big.Float)(nil))
var timeType = getType((*time.Time)(nil))
var listType = getType((*list.List)(nil))
var bytesType = getType(([]byte)(nil))
