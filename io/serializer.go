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
 * io/serializer.go                                       *
 *                                                        *
 * hprose seriaizer for Go.                               *
 *                                                        *
 * LastModified: Aug 20, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

// Serializer is a interface for serializing build-in type
type Serializer interface {
	Serialize(writer *Writer, v interface{})
}

type refSerializer struct {
	value Serializer
}

func (s *refSerializer) Serialize(writer *Writer, v interface{}) {
	if ok := writer.WriteRef(v); !ok {
		s.value.Serialize(writer, v)
	}
}

type nilSerializer struct{}

func (*nilSerializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteNil()
}

// Nil is an implementation of Serializer interface for serializing nil
var Nil = &nilSerializer{}

type boolSerializer struct{}

func (*boolSerializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteBool(v.(bool))
}

// Bool is an implementation of Serializer interface for serializing bool
var Bool = &boolSerializer{}

type intSerializer struct{}

func (*intSerializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteInt(int64(v.(int)))
}

// Int is an implementation of Serializer interface for serializing int
var Int = &intSerializer{}

type int8Serializer struct{}

func (*int8Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteInt32(int32(v.(int8)))
}

// Int8 is an implementation of Serializer interface for serializing int8
var Int8 = &int8Serializer{}

type int16Serializer struct{}

func (*int16Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteInt32(int32(v.(int16)))
}

// Int16 is an implementation of Serializer interface for serializing int16
var Int16 = &int16Serializer{}

type int32Serializer struct{}

func (*int32Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteInt32(v.(int32))
}

// Int32 is an implementation of Serializer interface for serializing int32
var Int32 = &int32Serializer{}

type int64Serializer struct{}

func (*int64Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteInt(v.(int64))
}

// Int64 is an implementation of Serializer interface for serializing int64
var Int64 = &int64Serializer{}

type uintSerializer struct{}

func (*uintSerializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteUint(uint64(v.(uint)))
}

// Uint is an implementation of Serializer interface for serializing uint
var Uint = &uintSerializer{}

type uint8Serializer struct{}

func (*uint8Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteUint(uint64(v.(uint8)))
}

// Uint8 is an implementation of Serializer interface for serializing uint8
var Uint8 = &uint8Serializer{}

type uint16Serializer struct{}

func (*uint16Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteUint(uint64(v.(uint16)))
}

// Uint16 is an implementation of Serializer interface for serializing uint16
var Uint16 = &uint16Serializer{}

type uint32Serializer struct{}

func (*uint32Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteUint(uint64(v.(uint32)))
}

// Uint32 is an implementation of Serializer interface for serializing uint32
var Uint32 = &uint32Serializer{}

type uint64Serializer struct{}

func (*uint64Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteUint(v.(uint64))
}

// Uint64 is an implementation of Serializer interface for serializing uint64
var Uint64 = &uint64Serializer{}

type uintptrSerializer struct{}

func (*uintptrSerializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteUint(uint64(v.(uintptr)))
}

// Uintptr is an implementation of Serializer interface for serializing uintptr
var Uintptr = &uintptrSerializer{}

type float32Serializer struct{}

func (*float32Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteFloat(float64(v.(float32)), 32)
}

// Float32 is an implementation of Serializer interface for serializing float32
var Float32 = &float32Serializer{}

type float64Serializer struct{}

func (*float64Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteFloat(v.(float64), 64)
}

// Float64 is an implementation of Serializer interface for serializing float64
var Float64 = &float64Serializer{}

type complex64Serializer struct{}

func (*complex64Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteComplex64(v.(complex64))
}

// Complex64 is an implementation of Serializer interface for serializing
// complex64
var Complex64 = &complex64Serializer{}

type complex128Serializer struct{}

func (*complex128Serializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteComplex128(v.(complex128))
}

// Complex128 is an implementation of Serializer interface for serializing
// complex128
var Complex128 = &complex128Serializer{}

type arraySerializer struct{}

func (*arraySerializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteArray(v)
}

// Array is an implementation of Serializer interface for serializing array
var Array = &arraySerializer{}

type sliceSerializer struct{}

func (*sliceSerializer) Serialize(writer *Writer, v interface{}) {
	writer.WriteSlice(v)
}

// Slice is an implementation of Serializer interface for serializing slice
var Slice = &sliceSerializer{}
