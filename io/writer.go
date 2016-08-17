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
 * io/writer.go                                           *
 *                                                        *
 * hprose writer for Go.                                  *
 *                                                        *
 * LastModified: Aug 17, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"io"
)

// Writer is a fine-grained operation struct for Hprose serialization
type Writer struct {
	Stream io.Writer
}

// Marshaler is a interface for serializing user custum type
type Marshaler interface {
	MarshalHprose(writer *Writer) error
}

// SerializerList stores a list of build-in type serializer
var SerializerList = [...]Serializer{
	invalidSerializer{},
	boolSerializer{},
	intSerializer{},
}

// WriteRef writes reference of an object
func (writer *Writer) WriteRef(v interface{}) (bool, error) {
	return false, nil
}

// SetRef add v to reference list, if WriteRef is call with the same v, it will
// write the reference index instead of v.
func (writer *Writer) SetRef(v interface{}) {

}
