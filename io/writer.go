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
 * LastModified: Aug 20, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/hprose/hprose-golang/util"
)

// Writer is a fine-grained operation struct for Hprose serialization
type Writer struct {
	Stream   *bytes.Buffer
	Simple   bool
	classref map[string]int
}

// Marshaler is a interface for serializing user custum type
type Marshaler interface {
	MarshalHprose(writer *Writer)
}

// SerializerList stores a list of build-in type serializer
var SerializerList = [...]Serializer{
	reflect.Invalid:       Nil,
	reflect.Bool:          Bool,
	reflect.Int:           Int,
	reflect.Int8:          Int8,
	reflect.Int16:         Int16,
	reflect.Int32:         Int32,
	reflect.Int64:         Int64,
	reflect.Uint:          Uint,
	reflect.Uint8:         Uint8,
	reflect.Uint16:        Uint16,
	reflect.Uint32:        Uint32,
	reflect.Uint64:        Uint64,
	reflect.Uintptr:       Uintptr,
	reflect.Float32:       Float32,
	reflect.Float64:       Float64,
	reflect.Complex64:     Complex64,
	reflect.Complex128:    Complex128,
	reflect.Array:         Array,
	reflect.Chan:          Nil,
	reflect.Func:          Nil,
	reflect.Interface:     Nil,
	reflect.Map:           Nil,
	reflect.Ptr:           Nil,
	reflect.Slice:         Nil,
	reflect.String:        Nil,
	reflect.Struct:        Nil,
	reflect.UnsafePointer: Nil,
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ uintptr
	ptr uintptr
}

// NewWriter is the constructor for Hprose Writer
func NewWriter(stream *bytes.Buffer, simple bool) *Writer {
	return &Writer{stream, simple, nil}
}

// Serialize a data v to stream
func (writer *Writer) Serialize(v interface{}) {
	if v == nil {
		writer.WriteNil()
	} else {
		SerializerList[reflect.TypeOf(v).Kind()].Serialize(writer, v)
	}
}

// WriteNil to stream
func (writer *Writer) WriteNil() {
	writer.Stream.WriteByte(TagNull)
}

// WriteBool to stream
func (writer *Writer) WriteBool(b bool) {
	s := writer.Stream
	if b {
		s.WriteByte(TagTrue)
	} else {
		s.WriteByte(TagFalse)
	}
}

// WriteInt32 to stream
func (writer *Writer) WriteInt32(i int32) {
	s := writer.Stream
	if (i >= 0) && (i <= 9) {
		s.WriteByte(byte('0' + i))
		return
	}
	s.WriteByte(TagInteger)
	s.Write(util.GetIntBytes(int64(i)))
	s.WriteByte(TagSemicolon)
}

// WriteInt to stream
func (writer *Writer) WriteInt(i int64) {
	s := writer.Stream
	if (i >= 0) && (i <= 9) {
		s.WriteByte(byte('0' + i))
		return
	}
	if (i >= math.MinInt32) && (i <= math.MaxInt32) {
		s.WriteByte(TagInteger)
	} else {
		s.WriteByte(TagLong)
	}
	s.Write(util.GetIntBytes(i))
	s.WriteByte(TagSemicolon)
}

// WriteUint to stream
func (writer *Writer) WriteUint(i uint64) {
	s := writer.Stream
	if (i >= 0) && (i <= 9) {
		s.WriteByte(byte('0' + i))
		return
	}
	if i <= math.MaxInt32 {
		s.WriteByte(TagInteger)
	} else {
		s.WriteByte(TagLong)
	}
	s.Write(util.GetUintBytes(i))
	s.WriteByte(TagSemicolon)
}

// WriteFloat to stream
func (writer *Writer) WriteFloat(f float64, bitSize int) {
	s := writer.Stream
	if f != f {
		s.WriteByte(TagNaN)
		return
	}
	if f > math.MaxFloat64 {
		s.Write([]byte{TagInfinity, TagPos})
		return
	}
	if f < -math.MaxFloat64 {
		s.Write([]byte{TagInfinity, TagNeg})
		return
	}
	var buf [64]byte
	s.WriteByte(TagDouble)
	s.Write(strconv.AppendFloat(buf[:0], f, 'g', -1, bitSize))
	s.WriteByte(TagSemicolon)
}

// WriteComplex64 to stream
func (writer *Writer) WriteComplex64(c complex64) {
	if imag(c) == 0 {
		writer.WriteFloat(float64(real(c)), 32)
		return
	}
	writer.WriteTuple(real(c), imag(c))
}

// WriteComplex128 to stream
func (writer *Writer) WriteComplex128(c complex128) {
	if imag(c) == 0 {
		writer.WriteFloat(real(c), 64)
		return
	}
	writer.WriteTuple(real(c), imag(c))
}

// WriteTuple to stream
func (writer *Writer) WriteTuple(tuple ...interface{}) {
	writer.SetRef(tuple)
	s := writer.Stream
	count := len(tuple)
	if count == 0 {
		s.Write([]byte{TagList, TagOpenbrace, TagClosebrace})
		return
	}
	s.WriteByte(TagList)
	s.Write(util.GetIntBytes(int64(count)))
	s.WriteByte(TagOpenbrace)
	for _, v := range tuple {
		writer.Serialize(v)
	}
	s.WriteByte(TagClosebrace)
}

// WriteArray to stream
func (writer *Writer) WriteArray(v interface{}) {
	t := reflect.TypeOf(v)
	count := t.Len()
	et := t.Elem()
	kind := et.Kind()
	ptr := (*emptyInterface)(unsafe.Pointer(&v)).ptr
	if kind == reflect.Uint8 {
		var bytes []byte
		byteSlice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
		byteSlice.Data = ptr
		byteSlice.Len = count
		byteSlice.Cap = count
		writer.WriteBytes(bytes)
		return
	}
	writer.SetRef(v)
	s := writer.Stream
	if count == 0 {
		s.Write([]byte{TagList, TagOpenbrace, TagClosebrace})
		return
	}
	s.WriteByte(TagList)
	s.Write(util.GetIntBytes(int64(count)))
	s.WriteByte(TagOpenbrace)
	serializer := SerializerList[kind]
	typ := (*emptyInterface)(unsafe.Pointer(&et)).ptr
	size := et.Size()
	for i := 0; i < count; i++ {
		var e interface{}
		es := (*emptyInterface)(unsafe.Pointer(&e))
		es.typ = typ
		es.ptr = ptr + uintptr(i)*size
		serializer.Serialize(writer, e)
	}
	s.WriteByte(TagClosebrace)
}

// WriteBytes to stream
func (writer *Writer) WriteBytes(bytes []byte) {
	writer.SetRef(bytes)
	s := writer.Stream
	count := len(bytes)
	if count == 0 {
		s.WriteByte(TagEmpty)
		return
	}
	s.WriteByte(TagBytes)
	s.Write(util.GetIntBytes(int64(count)))
	s.WriteByte(TagQuote)
	s.Write(bytes)
	s.WriteByte(TagQuote)
}

// WriteRef writes reference of an object to stream
func (writer *Writer) WriteRef(v interface{}) bool {
	return false
}

// SetRef add v to reference list, if WriteRef is call with the same v, it will
// write the reference index instead of v.
func (writer *Writer) SetRef(v interface{}) {

}
