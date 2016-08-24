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
 * LastModified: Aug 23, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
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
