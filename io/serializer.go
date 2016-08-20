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
	Serialize(writer *Writer, v interface{}) error
}

type refSerializer struct {
	value Serializer
}

func (s *refSerializer) Serialize(writer *Writer, v interface{}) error {
	if ok, err := writer.WriteRef(v); ok || err != nil {
		return err
	}
	return s.value.Serialize(writer, v)
}

type nilSerializer struct{}

func (*nilSerializer) Serialize(writer *Writer, v interface{}) (err error) {
	return writer.WriteNil()
}

type boolSerializer struct{}

func (*boolSerializer) Serialize(writer *Writer, v interface{}) (err error) {
	return writer.WriteBool(v.(bool))
}

type intSerializer struct{}

func (*intSerializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteInt(int64(v.(int)))
}

type int8Serializer struct{}

func (*int8Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteInt32(int32(v.(int8)))
}

type int16Serializer struct{}

func (*int16Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteInt32(int32(v.(int16)))
}

type int32Serializer struct{}

func (*int32Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteInt32(v.(int32))
}

type int64Serializer struct{}

func (*int64Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteInt(v.(int64))
}

type uintSerializer struct{}

func (*uintSerializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteUint(uint64(v.(uint)))
}

type uint8Serializer struct{}

func (*uint8Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteUint(uint64(v.(uint8)))
}

type uint16Serializer struct{}

func (*uint16Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteUint(uint64(v.(uint16)))
}

type uint32Serializer struct{}

func (*uint32Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteUint(uint64(v.(uint32)))
}

type uint64Serializer struct{}

func (*uint64Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteUint(v.(uint64))
}

type float32Serializer struct{}

func (*float32Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteFloat(float64(v.(float32)), 32)
}

type float64Serializer struct{}

func (*float64Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteFloat(v.(float64), 64)
}

type complex64Serializer struct{}

func (*complex64Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteComplex64(v.(complex64))
}

type complex128Serializer struct{}

func (*complex128Serializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteComplex128(v.(complex128))
}

type arraySerializer struct{}

func (*arraySerializer) Serialize(writer *Writer, v interface{}) error {
	return writer.WriteArray(v)
}
