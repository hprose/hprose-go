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
 * LastModified: Aug 22, 2016                             *
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

type emptyInterface struct {
	typ uintptr
	ptr unsafe.Pointer
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
		v := reflect.ValueOf(v)
		valueEncoders[v.Kind()](writer, v)
	}
}

// WriteValue to stream
func (writer *Writer) WriteValue(v reflect.Value) {
	valueEncoders[v.Kind()](writer, v)
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
	if i >= 0 && i <= 9 {
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
	if i >= 0 && i <= 9 {
		s.WriteByte(byte('0' + i))
		return
	}
	if i >= math.MinInt32 && i <= math.MaxInt32 {
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
	if i >= 0 && i <= 9 {
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

func (writer *Writer) writeListHeader(count int) {
	s := writer.Stream
	s.WriteByte(TagList)
	s.Write(util.GetIntBytes(int64(count)))
	s.WriteByte(TagOpenbrace)
}

func (writer *Writer) writeListFooter() {
	writer.Stream.WriteByte(TagClosebrace)
}

func (writer *Writer) writeEmptyList() {
	writer.Stream.Write([]byte{TagList, TagOpenbrace, TagClosebrace})
}

// WriteComplex64 to stream
func (writer *Writer) WriteComplex64(c complex64) {
	if imag(c) == 0 {
		writer.WriteFloat(float64(real(c)), 32)
		return
	}
	writer.SetRef(nil)
	writer.writeListHeader(2)
	writer.WriteFloat(float64(real(c)), 32)
	writer.WriteFloat(float64(imag(c)), 32)
	writer.writeListFooter()
}

// WriteComplex128 to stream
func (writer *Writer) WriteComplex128(c complex128) {
	if imag(c) == 0 {
		writer.WriteFloat(real(c), 64)
		return
	}
	writer.SetRef(nil)
	writer.writeListHeader(2)
	writer.WriteFloat(real(c), 64)
	writer.WriteFloat(imag(c), 64)
	writer.writeListFooter()
}

// WriteTuple to stream
func (writer *Writer) WriteTuple(tuple ...interface{}) {
	writer.SetRef(nil)
	count := len(tuple)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	for _, v := range tuple {
		writer.Serialize(v)
	}
	writer.writeListFooter()
}

func (writer *Writer) writeListBody(list reflect.Value, count int) {
	for i := 0; i < count; i++ {
		writer.WriteValue(list.Index(i))
	}
}

// WriteArray to stream
func (writer *Writer) writeArray(v reflect.Value) {
	writer.SetRef(nil)
	count := v.Len()
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	kind := v.Type().Elem().Kind()
	if encoder := sliceBodyEncoders[kind]; encoder != nil {
		ptr := (*emptyInterface)(unsafe.Pointer(&v)).ptr
		sliceHeader := reflect.SliceHeader{
			Data: uintptr(ptr),
			Len:  count,
			Cap:  count,
		}
		encoder(writer, unsafe.Pointer(&sliceHeader))
	} else {
		writer.writeListBody(v, count)
	}
	writer.writeListFooter()
}

// WriteArray to stream
func (writer *Writer) WriteArray(v interface{}) {
	array := reflect.ValueOf(v)
	writer.writeArray(array)
}

func (writer *Writer) writeSlice(v reflect.Value) {
	kind := v.Type().Elem().Kind()
	ptr := (*emptyInterface)(unsafe.Pointer(&v)).ptr
	if kind == reflect.Uint8 {
		writer.WriteBytes(*(*[]byte)(ptr))
		return
	}
	writer.SetRef(v)
	count := v.Len()
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	if encoder := sliceBodyEncoders[kind]; encoder != nil {
		encoder(writer, ptr)
	} else {
		writer.writeListBody(v, count)
	}
	writer.writeListFooter()
}

// WriteSlice to stream
func (writer *Writer) WriteSlice(v interface{}) {
	slice := reflect.ValueOf(v)
	writer.writeSlice(slice)
}

// WriteBoolSlice to stream
func (writer *Writer) WriteBoolSlice(slice []bool) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	boolSliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteIntSlice to stream
func (writer *Writer) WriteIntSlice(slice []int) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	intSliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteInt8Slice to stream
func (writer *Writer) WriteInt8Slice(slice []int8) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	int8SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteInt16Slice to stream
func (writer *Writer) WriteInt16Slice(slice []int16) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	int16SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteInt32Slice to stream
func (writer *Writer) WriteInt32Slice(slice []int32) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	int32SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteInt64Slice to stream
func (writer *Writer) WriteInt64Slice(slice []int64) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	int64SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteUintSlice to stream
func (writer *Writer) WriteUintSlice(slice []uint) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	intSliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteUint8Slice to stream
func (writer *Writer) WriteUint8Slice(slice []uint8) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	uint8SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteUint16Slice to stream
func (writer *Writer) WriteUint16Slice(slice []uint16) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	uint16SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteUint32Slice to stream
func (writer *Writer) WriteUint32Slice(slice []uint32) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	uint32SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteUint64Slice to stream
func (writer *Writer) WriteUint64Slice(slice []uint64) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	uint64SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteUintptrSlice to stream
func (writer *Writer) WriteUintptrSlice(slice []uintptr) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	uintptrSliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteFloat32Slice to stream
func (writer *Writer) WriteFloat32Slice(slice []float32) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	float32SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteFloat64Slice to stream
func (writer *Writer) WriteFloat64Slice(slice []float64) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	float64SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteComplex64Slice to stream
func (writer *Writer) WriteComplex64Slice(slice []complex64) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	complex64SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
}

// WriteComplex128Slice to stream
func (writer *Writer) WriteComplex128Slice(slice []complex128) {
	writer.SetRef(slice)
	count := len(slice)
	if count == 0 {
		writer.writeEmptyList()
		return
	}
	writer.writeListHeader(count)
	complex128SliceEncoder(writer, unsafe.Pointer(&slice))
	writer.writeListFooter()
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
