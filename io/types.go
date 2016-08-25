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
 * LastModified: Aug 25, 2016                             *
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

func getType(v interface{}) uintptr {
	return *(*uintptr)(unsafe.Pointer(&v))
}

var bigIntType = getType((*big.Int)(nil))
var bigRatType = getType((*big.Rat)(nil))
var bigFloatType = getType((*big.Float)(nil))
var timeType = getType((*time.Time)(nil))
var listType = getType((*list.List)(nil))

var bytesType = getType(([]byte)(nil))

var stringStringMapType = getType((map[string]string)(nil))
var stringInterfaceMapType = getType((map[string]interface{})(nil))
var stringIntMapType = getType((map[string]int)(nil))
var intIntMapType = getType((map[int]int)(nil))
var intStringMapType = getType((map[int]string)(nil))
var intInterfaceMapType = getType((map[int]interface{})(nil))
var interfaceInterfaceMapType = getType((map[interface{}]interface{})(nil))
