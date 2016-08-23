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
	"reflect"
)

var bigIntType = reflect.TypeOf((*big.Int)(nil))
var bigRatType = reflect.TypeOf((*big.Rat)(nil))
var bigFloatType = reflect.TypeOf((*big.Float)(nil))
