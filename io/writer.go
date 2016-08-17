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
	"reflect"
)

// Writer is a fine-grained operation struct for Hprose serialization
type Writer struct {
	Stream   io.Writer
	Simple   bool
	classref map[string]int
}

// Marshaler is a interface for serializing user custum type
type Marshaler interface {
	MarshalHprose(writer *Writer) error
}

// SerializerList stores a list of build-in type serializer
var SerializerList = [...]Serializer{
	reflect.Invalid:       nilSerializer{},
	reflect.Bool:          boolSerializer{},
	reflect.Int:           intSerializer{},
	reflect.Int8:          int8Serializer{},
	reflect.Int16:         int16Serializer{},
	reflect.Int32:         int32Serializer{},
	reflect.Int64:         int64Serializer{},
	reflect.Uint:          nilSerializer{},
	reflect.Uint8:         nilSerializer{},
	reflect.Uint16:        nilSerializer{},
	reflect.Uint32:        nilSerializer{},
	reflect.Uint64:        nilSerializer{},
	reflect.Uintptr:       nilSerializer{},
	reflect.Float32:       nilSerializer{},
	reflect.Float64:       nilSerializer{},
	reflect.Complex64:     nilSerializer{},
	reflect.Complex128:    nilSerializer{},
	reflect.Array:         nilSerializer{},
	reflect.Chan:          nilSerializer{},
	reflect.Func:          nilSerializer{},
	reflect.Interface:     nilSerializer{},
	reflect.Map:           nilSerializer{},
	reflect.Ptr:           nilSerializer{},
	reflect.Slice:         nilSerializer{},
	reflect.String:        nilSerializer{},
	reflect.Struct:        nilSerializer{},
	reflect.UnsafePointer: nilSerializer{},
}

// NewWriter is the constructor for Hprose Writer
func NewWriter(stream io.Writer, simple bool) *Writer {
	return &Writer{stream, simple, nil}
}

// Serialize a data v to stream
func (writer *Writer) Serialize(v interface{}) error {
	if v == nil {
		return SerializerList[0].Serialize(writer, nil)
	}
	return SerializerList[reflect.TypeOf(v).Kind()].Serialize(writer, v)
}

// WriteRef writes reference of an object
func (writer *Writer) WriteRef(v interface{}) (bool, error) {
	return false, nil
}

// SetRef add v to reference list, if WriteRef is call with the same v, it will
// write the reference index instead of v.
func (writer *Writer) SetRef(v interface{}) {

}
