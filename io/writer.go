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
 * LastModified: Aug 29, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"container/list"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

// Writer is a fine-grained operation struct for Hprose serialization
type Writer struct {
	Stream   *bytes.Buffer
	Simple   bool
	classref map[uintptr]int
	ref      map[uintptr]int
	refcount int
}

// NewWriter is the constructor for Hprose Writer
func NewWriter(stream *bytes.Buffer, simple bool) (writer *Writer) {
	writer = new(Writer)
	writer.Stream = stream
	writer.Simple = simple
	writer.classref = map[uintptr]int{}
	if !simple {
		writer.ref = map[uintptr]int{}
	}
	return writer
}

// Serialize a data v to stream
func (writer *Writer) Serialize(v interface{}) {
	if v == nil {
		writer.WriteNil()
	} else {
		writer.WriteValue(reflect.ValueOf(v))
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
	var buf [20]byte
	s.Write(getIntBytes(buf[:], i))
	s.WriteByte(TagSemicolon)
}

// WriteUint to stream
func (writer *Writer) WriteUint(i uint64) {
	s := writer.Stream
	if i <= 9 {
		s.WriteByte(byte('0' + i))
		return
	}
	if i <= math.MaxInt32 {
		s.WriteByte(TagInteger)
	} else {
		s.WriteByte(TagLong)
	}
	var buf [20]byte
	s.Write(getUintBytes(buf[:], i))
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
	s.WriteByte(TagDouble)
	var buf [64]byte
	s.Write(strconv.AppendFloat(buf[:0], f, 'g', -1, bitSize))
	s.WriteByte(TagSemicolon)
}

// WriteComplex64 to stream
func (writer *Writer) WriteComplex64(c complex64) {
	if imag(c) == 0 {
		writer.WriteFloat(float64(real(c)), 32)
		return
	}
	writer.SetRef(nil)
	writeListHeader(writer, 2)
	writer.WriteFloat(float64(real(c)), 32)
	writer.WriteFloat(float64(imag(c)), 32)
	writeListFooter(writer)
}

// WriteComplex128 to stream
func (writer *Writer) WriteComplex128(c complex128) {
	if imag(c) == 0 {
		writer.WriteFloat(real(c), 64)
		return
	}
	writer.SetRef(nil)
	writeListHeader(writer, 2)
	writer.WriteFloat(real(c), 64)
	writer.WriteFloat(imag(c), 64)
	writeListFooter(writer)
}

// WriteString to stream
func (writer *Writer) WriteString(str string) {
	length := utf16Length(str)
	switch {
	case length == 0:
		writer.Stream.WriteByte(TagEmpty)
	case length < 0:
		writer.WriteBytes(*(*[]byte)(unsafe.Pointer(&str)))
	case length == 1:
		writer.Stream.WriteByte(TagUTF8Char)
		writer.Stream.WriteString(str)
	default:
		writer.SetRef(nil)
		writeString(writer, str, length)
	}
}

// WriteBytes to stream
func (writer *Writer) WriteBytes(bytes []byte) {
	writer.SetRef(nil)
	writeBytes(writer, bytes)
}

// WriteBigInt to stream
func (writer *Writer) WriteBigInt(bi *big.Int) {
	s := writer.Stream
	s.WriteByte(TagLong)
	s.WriteString(bi.String())
	s.WriteByte(TagSemicolon)
}

// WriteBigRat to stream
func (writer *Writer) WriteBigRat(br *big.Rat) {
	if br.IsInt() {
		writer.WriteBigInt(br.Num())
	} else {
		str := br.String()
		writer.SetRef(nil)
		writeString(writer, str, len(str))
	}
}

// WriteBigFloat to stream
func (writer *Writer) WriteBigFloat(bf *big.Float) {
	s := writer.Stream
	s.WriteByte(TagDouble)
	var buf [64]byte
	s.Write(bf.Append(buf[:0], 'g', -1))
	s.WriteByte(TagSemicolon)
}

// WriteTime to stream
func (writer *Writer) WriteTime(t *time.Time) {
	ptr := unsafe.Pointer(t)
	if writer.WriteRef(ptr) {
		return
	}
	writer.SetRef(ptr)
	s := writer.Stream
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	nsec := t.Nanosecond()
	tag := TagSemicolon
	if t.Location() == time.UTC {
		tag = TagUTC
	}
	var buf [27]byte
	if hour == 0 && min == 0 && sec == 0 && nsec == 0 {
		datelen := formatDate(buf[:], year, int(month), day)
		buf[datelen] = tag
		s.Write(buf[:datelen+1])
	} else if year == 1970 && month == 1 && day == 1 {
		timelen := formatTime(buf[:], hour, min, sec, nsec)
		buf[timelen] = tag
		s.Write(buf[:timelen+1])
	} else {
		datelen := formatDate(buf[:], year, int(month), day)
		timelen := formatTime(buf[datelen:], hour, min, sec, nsec)
		datetimelen := datelen + timelen
		buf[datetimelen] = tag
		s.Write(buf[:datetimelen+1])
	}
}

// WriteList to stream
func (writer *Writer) WriteList(lst *list.List) {
	ptr := unsafe.Pointer(lst)
	if writer.WriteRef(ptr) {
		return
	}
	writer.SetRef(ptr)
	count := lst.Len()
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	for e := lst.Front(); e != nil; e = e.Next() {
		writer.Serialize(e.Value)
	}
	writeListFooter(writer)
}

// WriteTuple to stream
func (writer *Writer) WriteTuple(tuple ...interface{}) {
	writer.SetRef(nil)
	count := len(tuple)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	for _, v := range tuple {
		writer.Serialize(v)
	}
	writeListFooter(writer)
}

// WriteBoolSlice to stream
func (writer *Writer) WriteBoolSlice(slice []bool) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	boolSliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteIntSlice to stream
func (writer *Writer) WriteIntSlice(slice []int) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	intSliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteInt8Slice to stream
func (writer *Writer) WriteInt8Slice(slice []int8) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	int8SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteInt16Slice to stream
func (writer *Writer) WriteInt16Slice(slice []int16) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	int16SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteInt32Slice to stream
func (writer *Writer) WriteInt32Slice(slice []int32) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	int32SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteInt64Slice to stream
func (writer *Writer) WriteInt64Slice(slice []int64) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	int64SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteUintSlice to stream
func (writer *Writer) WriteUintSlice(slice []uint) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	uintSliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteUint8Slice to stream
func (writer *Writer) WriteUint8Slice(slice []uint8) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	uint8SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteUint16Slice to stream
func (writer *Writer) WriteUint16Slice(slice []uint16) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	uint16SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteUint32Slice to stream
func (writer *Writer) WriteUint32Slice(slice []uint32) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	uint32SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteUint64Slice to stream
func (writer *Writer) WriteUint64Slice(slice []uint64) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	uint64SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteUintptrSlice to stream
func (writer *Writer) WriteUintptrSlice(slice []uintptr) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	uintptrSliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteFloat32Slice to stream
func (writer *Writer) WriteFloat32Slice(slice []float32) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	float32SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteFloat64Slice to stream
func (writer *Writer) WriteFloat64Slice(slice []float64) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	float64SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteComplex64Slice to stream
func (writer *Writer) WriteComplex64Slice(slice []complex64) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	complex64SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteComplex128Slice to stream
func (writer *Writer) WriteComplex128Slice(slice []complex128) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	complex128SliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteStringSlice to stream
func (writer *Writer) WriteStringSlice(slice []string) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	stringSliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteBytesSlice to stream
func (writer *Writer) WriteBytesSlice(slice [][]byte) {
	writer.SetRef(nil)
	count := len(slice)
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	bytesSliceEncoder(writer, unsafe.Pointer(&slice))
	writeListFooter(writer)
}

// WriteRef writes reference of an object to stream
func (writer *Writer) WriteRef(ref unsafe.Pointer) bool {
	if writer.Simple {
		return false
	}
	return writeRef(writer, ref)
}

// SetRef add v to reference list, if WriteRef is call with the same v, it will
// write the reference index instead of v.
func (writer *Writer) SetRef(ref unsafe.Pointer) {
	if writer.Simple {
		return
	}
	if ref != nil {
		writer.ref[uintptr(ref)] = writer.refcount
	}
	writer.refcount++
}

// Reset the reference counter
func (writer *Writer) Reset() {
	for k := range writer.classref {
		delete(writer.classref, k)
	}
	if writer.Simple {
		return
	}
	writer.refcount = 0
	for k := range writer.ref {
		delete(writer.ref, k)
	}
}

// private functions

func writeRef(writer *Writer, ref unsafe.Pointer) bool {
	n, found := writer.ref[uintptr(ref)]
	if found {
		s := writer.Stream
		s.WriteByte(TagRef)
		var buf [20]byte
		s.Write(getIntBytes(buf[:], int64(n)))
		s.WriteByte(TagSemicolon)
	}
	return found
}

func writeString(writer *Writer, str string, length int) {
	s := writer.Stream
	s.WriteByte(TagString)
	var buf [20]byte
	s.Write(getIntBytes(buf[:], int64(length)))
	s.WriteByte(TagQuote)
	s.WriteString(str)
	s.WriteByte(TagQuote)
}

func writeBytes(writer *Writer, bytes []byte) {
	s := writer.Stream
	count := len(bytes)
	if count == 0 {
		s.Write([]byte{TagBytes, TagQuote, TagQuote})
		return
	}
	s.WriteByte(TagBytes)
	var buf [20]byte
	s.Write(getIntBytes(buf[:], int64(count)))
	s.WriteByte(TagQuote)
	s.Write(bytes)
	s.WriteByte(TagQuote)
}

func writeListHeader(writer *Writer, count int) {
	s := writer.Stream
	s.WriteByte(TagList)
	var buf [20]byte
	s.Write(getIntBytes(buf[:], int64(count)))
	s.WriteByte(TagOpenbrace)
}

func writeListBody(writer *Writer, list reflect.Value, count int) {
	for i := 0; i < count; i++ {
		e := list.Index(i)
		valueEncoders[e.Kind()](writer, e)
	}
}

func writeListFooter(writer *Writer) {
	writer.Stream.WriteByte(TagClosebrace)
}

func writeEmptyList(writer *Writer) {
	writer.Stream.Write([]byte{TagList, TagOpenbrace, TagClosebrace})
}

func writeArray(writer *Writer, v reflect.Value) {
	st := reflect.SliceOf(v.Type().Elem())
	sliceType := (*emptyInterface)(unsafe.Pointer(&st)).ptr
	count := v.Len()
	if sliceType == bytesType {
		sliceHeader := reflect.SliceHeader{
			Data: (*emptyInterface)(unsafe.Pointer(&v)).ptr,
			Len:  count,
			Cap:  count,
		}
		writeBytes(writer, *(*[]byte)(unsafe.Pointer(&sliceHeader)))
		return
	}
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	encoder := sliceBodyEncoders[sliceType]
	if encoder != nil {
		sliceHeader := reflect.SliceHeader{
			Data: (*emptyInterface)(unsafe.Pointer(&v)).ptr,
			Len:  count,
			Cap:  count,
		}
		encoder(writer, unsafe.Pointer(&sliceHeader))
	} else {
		writeListBody(writer, v, count)
	}
	writeListFooter(writer)
}

func writeSlice(writer *Writer, v reflect.Value) {
	val := (*reflectValue)(unsafe.Pointer(&v))
	if val.typ == bytesType {
		writeBytes(writer, v.Bytes())
		return
	}
	count := v.Len()
	if count == 0 {
		writeEmptyList(writer)
		return
	}
	writeListHeader(writer, count)
	encoder := sliceBodyEncoders[val.typ]
	if encoder != nil {
		encoder(writer, val.ptr)
	} else {
		writeListBody(writer, v, count)
	}
	writeListFooter(writer)
}

func writeEmptyMap(writer *Writer) {
	writer.Stream.Write([]byte{TagMap, TagOpenbrace, TagClosebrace})
}

func writeMapHeader(writer *Writer, count int) {
	s := writer.Stream
	s.WriteByte(TagMap)
	var buf [20]byte
	s.Write(getIntBytes(buf[:], int64(count)))
	s.WriteByte(TagOpenbrace)
}

func writeMapBody(writer *Writer, v reflect.Value) {
	mapType := v.Type()
	keyEncoder := valueEncoders[mapType.Key().Kind()]
	valueEncoder := valueEncoders[mapType.Elem().Kind()]
	keys := v.MapKeys()
	for _, key := range keys {
		keyEncoder(writer, key)
		valueEncoder(writer, v.MapIndex(key))
	}
}

func writeMapFooter(writer *Writer) {
	writer.Stream.WriteByte(TagClosebrace)
}

func writeMap(writer *Writer, v reflect.Value) {
	count := v.Len()
	if count == 0 {
		writeEmptyMap(writer)
		return
	}
	writeMapHeader(writer, count)
	val := (*reflectValue)(unsafe.Pointer(&v))
	mapEncoder := mapBodyEncoders[val.typ]
	if mapEncoder != nil {
		mapEncoder(writer, unsafe.Pointer(&val.ptr))
	} else {
		writeMapBody(writer, v)
	}
	writeMapFooter(writer)
}

func writeMapPtr(writer *Writer, v reflect.Value) {
	count := v.Len()
	if count == 0 {
		writeEmptyMap(writer)
		return
	}
	writeMapHeader(writer, count)
	val := (*reflectValue)(unsafe.Pointer(&v))
	mapEncoder := mapBodyEncoders[val.typ]
	if mapEncoder != nil {
		mapEncoder(writer, unsafe.Pointer(val.ptr))
	} else {
		writeMapBody(writer, v)
	}
	writeMapFooter(writer)
}

func writeStruct(writer *Writer, v reflect.Value) {
	val := (*reflectValue)(unsafe.Pointer(&v))
	cache := getStructCache(v.Type().Elem())
	index, found := writer.classref[val.typ]
	if !found {
		writer.Stream.Write(cache.Data)
		if !writer.Simple {
			writer.refcount += len(cache.Fields)
		}
		index = len(writer.classref)
		writer.classref[val.typ] = index
	}
	writer.SetRef(val.ptr)
	s := writer.Stream
	s.WriteByte(TagObject)
	var buf [20]byte
	s.Write(getIntBytes(buf[:], int64(index)))
	s.WriteByte(TagOpenbrace)
	fields := cache.Fields
	for _, field := range fields {
		var f interface{}
		fp := (*emptyInterface)(unsafe.Pointer(&f))
		fp.typ = field.Type
		fp.ptr = uintptr(val.ptr) + field.Offset
		if field.Kind == reflect.Ptr {
			fp.ptr = **(**uintptr)(unsafe.Pointer(&fp.ptr))
		}
		writer.Serialize(f)
	}
	s.WriteByte(TagClosebrace)
}
