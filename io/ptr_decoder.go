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
 * io/ptr_decoder.go                                      *
 *                                                        *
 * hprose ptr decoder for Go.                             *
 *                                                        *
 * LastModified: Sep 10, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "reflect"

func ptrDecoder(r *Reader, v reflect.Value) {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	r.ReadValue(v.Elem())
}
