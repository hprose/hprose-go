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
 * LastModified: Aug 17, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"math"

	"github.com/hprose/hprose-golang/util"
)

// Serializer is a interface for serializing build-in type
type Serializer interface {
	Serialize(writer *Writer, v interface{}) error
}

type refSerializer struct {
	value Serializer
}

func (s refSerializer) Serialize(writer *Writer, v interface{}) error {
	if ok, err := writer.WriteRef(v); ok || err != nil {
		return err
	}
	return s.value.Serialize(writer, v)
}

type nilSerializer struct{}

func (nilSerializer) Serialize(writer *Writer, v interface{}) (err error) {
	_, err = writer.Stream.Write([]byte{TagNull})
	return err
}

type boolSerializer struct{}

func (boolSerializer) Serialize(writer *Writer, v interface{}) (err error) {
	var tag byte
	if v.(bool) {
		tag = TagTrue
	} else {
		tag = TagFalse
	}
	_, err = writer.Stream.Write([]byte{tag})
	return err
}

func serializeInt32(writer *Writer, i int32) (err error) {
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

func serializeInt64(writer *Writer, i int64) (err error) {
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

func serializeUint(writer *Writer, i uint64) (err error) {
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

type intSerializer struct{}

func (intSerializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeInt64(writer, int64(v.(int)))
}

type int8Serializer struct{}

func (int8Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeInt32(writer, int32(v.(int8)))
}

type int16Serializer struct{}

func (int16Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeInt32(writer, int32(v.(int16)))
}

type int32Serializer struct{}

func (int32Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeInt32(writer, v.(int32))
}

type int64Serializer struct{}

func (int64Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeInt64(writer, v.(int64))
}

type uintSerializer struct{}

func (uintSerializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeUint(writer, uint64(v.(uint)))
}

type uint8Serializer struct{}

func (uint8Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeUint(writer, uint64(v.(uint8)))
}

type uint16Serializer struct{}

func (uint16Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeUint(writer, uint64(v.(uint16)))
}

type uint32Serializer struct{}

func (uint32Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeUint(writer, uint64(v.(uint32)))
}

type uint64Serializer struct{}

func (uint64Serializer) Serialize(writer *Writer, v interface{}) (err error) {
	return serializeUint(writer, v.(uint64))
}
