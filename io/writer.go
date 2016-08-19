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
 * LastModified: Aug 19, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"io"
	"math"
	"reflect"
	"strconv"

	"github.com/hprose/hprose-golang/util"
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
	reflect.Invalid:       &nilSerializer{},
	reflect.Bool:          &boolSerializer{},
	reflect.Int:           &intSerializer{},
	reflect.Int8:          &int8Serializer{},
	reflect.Int16:         &int16Serializer{},
	reflect.Int32:         &int32Serializer{},
	reflect.Int64:         &int64Serializer{},
	reflect.Uint:          &uintSerializer{},
	reflect.Uint8:         &uint8Serializer{},
	reflect.Uint16:        &uint16Serializer{},
	reflect.Uint32:        &uint32Serializer{},
	reflect.Uint64:        &uint64Serializer{},
	reflect.Uintptr:       &nilSerializer{},
	reflect.Float32:       &float32Serializer{},
	reflect.Float64:       &float64Serializer{},
	reflect.Complex64:     &complex64Serializer{},
	reflect.Complex128:    &complex128Serializer{},
	reflect.Array:         &arraySerializer{},
	reflect.Chan:          &nilSerializer{},
	reflect.Func:          &nilSerializer{},
	reflect.Interface:     &nilSerializer{},
	reflect.Map:           &nilSerializer{},
	reflect.Ptr:           &nilSerializer{},
	reflect.Slice:         &nilSerializer{},
	reflect.String:        &nilSerializer{},
	reflect.Struct:        &nilSerializer{},
	reflect.UnsafePointer: &nilSerializer{},
}

// NewWriter is the constructor for Hprose Writer
func NewWriter(stream io.Writer, simple bool) *Writer {
	return &Writer{stream, simple, nil}
}

// Serialize a data v to stream
func (writer *Writer) Serialize(v interface{}) error {
	if v == nil {
		return writer.WriteNil()
	}
	return SerializerList[reflect.TypeOf(v).Kind()].Serialize(writer, v)
}

// WriteNil to stream
func (writer *Writer) WriteNil() (err error) {
	_, err = writer.Stream.Write([]byte{TagNull})
	return err
}

// WriteBool to stream
func (writer *Writer) WriteBool(b bool) (err error) {
	s := writer.Stream
	if b {
		_, err = s.Write([]byte{TagTrue})
	} else {
		_, err = s.Write([]byte{TagFalse})
	}
	return err
}

// WriteInt32 to stream
func (writer *Writer) WriteInt32(i int32) (err error) {
	s := writer.Stream
	if (i >= 0) && (i <= 9) {
		_, err = s.Write([]byte{byte('0' + i)})
		return err
	}
	if _, err = s.Write([]byte{TagInteger}); err == nil {
		_, err = s.Write(util.GetIntBytes(int64(i)))
	}
	if err == nil {
		_, err = s.Write([]byte{TagSemicolon})
	}
	return err
}

// WriteInt to stream
func (writer *Writer) WriteInt(i int64) (err error) {
	s := writer.Stream
	if (i >= 0) && (i <= 9) {
		_, err = s.Write([]byte{byte('0' + i)})
		return err
	}
	if (i >= math.MinInt32) && (i <= math.MaxInt32) {
		_, err = s.Write([]byte{TagInteger})
	} else {
		_, err = s.Write([]byte{TagLong})
	}
	if err == nil {
		_, err = s.Write(util.GetIntBytes(i))
	}
	if err == nil {
		_, err = s.Write([]byte{TagSemicolon})
	}
	return err
}

// WriteUint to stream
func (writer *Writer) WriteUint(i uint64) (err error) {
	s := writer.Stream
	if (i >= 0) && (i <= 9) {
		_, err = s.Write([]byte{byte('0' + i)})
		return err
	}
	if i <= math.MaxInt32 {
		_, err = s.Write([]byte{TagInteger})
	} else {
		_, err = s.Write([]byte{TagLong})
	}
	if err == nil {
		_, err = s.Write(util.GetUintBytes(i))
	}
	if err == nil {
		_, err = s.Write([]byte{TagSemicolon})
	}
	return err
}

// WriteFloat to stream
func (writer *Writer) WriteFloat(f float64, bitSize int) (err error) {
	s := writer.Stream
	if f != f {
		_, err = s.Write([]byte{TagNaN})
		return err
	}
	if f > math.MaxFloat64 {
		_, err = s.Write([]byte{TagInfinity, TagPos})
		return err
	}
	if f < -math.MaxFloat64 {
		_, err = s.Write([]byte{TagInfinity, TagNeg})
		return err
	}
	if _, err = s.Write([]byte{TagDouble}); err == nil {
		var buf [32]byte
		_, err = s.Write(strconv.AppendFloat(buf[:0], f, 'g', -1, bitSize))
	}
	if err == nil {
		_, err = s.Write([]byte{TagSemicolon})
	}
	return err
}

// WriteComplex64 to stream
func (writer *Writer) WriteComplex64(c complex64) error {
	if imag(c) == 0 {
		return writer.WriteFloat(float64(real(c)), 32)
	}
	return writer.WriteTuple(real(c), imag(c))
}

// WriteComplex128 to stream
func (writer *Writer) WriteComplex128(c complex128) error {
	if imag(c) == 0 {
		return writer.WriteFloat(real(c), 64)
	}
	return writer.WriteTuple(real(c), imag(c))
}

// WriteTuple to stream
func (writer *Writer) WriteTuple(tuple ...interface{}) (err error) {
	writer.SetRef(tuple)
	s := writer.Stream
	count := len(tuple)
	if count == 0 {
		_, err = s.Write([]byte{TagList, TagOpenbrace, TagClosebrace})
		return err
	}
	if _, err = s.Write([]byte{TagList}); err == nil {
		_, err = s.Write(util.GetIntBytes(int64(count)))
	}
	if err == nil {
		_, err = s.Write([]byte{TagOpenbrace})
	}
	for _, v := range tuple {
		if err == nil {
			err = writer.Serialize(v)
		}
	}
	if err == nil {
		_, err = s.Write([]byte{TagClosebrace})
	}
	return err
}

// WriterArray to stream
func (writer *Writer) WriterArray(v interface{}) (err error) {
	writer.SetRef(v)
	s := writer.Stream
	array := reflect.ValueOf(v)
	count := array.Len()
	if count == 0 {
		_, err = s.Write([]byte{TagList, TagOpenbrace, TagClosebrace})
		return err
	}
	if _, err = s.Write([]byte{TagList}); err == nil {
		_, err = s.Write(util.GetIntBytes(int64(count)))
	}
	if err == nil {
		_, err = s.Write([]byte{TagOpenbrace})
	}
	serializer := SerializerList[array.Type().Elem().Kind()]
	for i := 0; i < count; i++ {
		if err == nil {
			err = serializer.Serialize(writer, array.Index(i).Interface())
		}
	}
	if err == nil {
		_, err = s.Write([]byte{TagClosebrace})
	}
	return err
}

// WriterBytes to stream
func (writer *Writer) WriterBytes(bytes []byte) (err error) {
	writer.SetRef(bytes)
	s := writer.Stream
	count := len(bytes)
	if count == 0 {
		_, err = s.Write([]byte{TagEmpty})
		return err
	}
	if _, err = s.Write([]byte{TagBytes}); err == nil {
		_, err = s.Write(util.GetIntBytes(int64(count)))
	}
	if err == nil {
		_, err = s.Write([]byte{TagQuote})
	}
	if err == nil {
		_, err = s.Write(bytes)
	}
	if err == nil {
		_, err = s.Write([]byte{TagQuote})
	}
	return err
}

// WriteRef writes reference of an object to stream
func (writer *Writer) WriteRef(v interface{}) (bool, error) {
	return false, nil
}

// SetRef add v to reference list, if WriteRef is call with the same v, it will
// write the reference index instead of v.
func (writer *Writer) SetRef(v interface{}) {

}
