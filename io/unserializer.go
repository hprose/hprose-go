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
 * io/unserializer.go                                     *
 *                                                        *
 * hprose unserializer for Go.                            *
 *                                                        *
 * LastModified: Aug 16, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "reflect"

// Unserializer is a interface for unserializing build-in type
type Unserializer interface {
	Unserialize(reader *Reader, tag byte, typ reflect.Type) (interface{}, error)
}
