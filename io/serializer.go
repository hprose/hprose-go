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
 * LastModified: Aug 15, 2016                             *
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

type invalidSerializer struct{}

func (invalidSerializer) Serialize(writer *Writer, v interface{}) (err error) {
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

type intSerializer struct{}

func (intSerializer) Serializer(writer *Writer, v interface{}) (err error) {
	s := writer.Stream
	i := v.(int)
	if (i >= 0) && (i <= 9) {
		_, err = s.Write([]byte{byte('0' + i)})
		return err
	}
	if (i >= math.MinInt32) && (i <= math.MaxInt32) {
		_, err = s.Write([]byte{TagInteger})
	} else {
		_, err = s.Write([]byte{TagLong})
	}
	if err != nil {
		return err
	}
	if _, err = s.Write(util.GetInt64Bytes(int64(i))); err != nil {
		return err
	}
	_, err = s.Write([]byte{TagSemicolon})
	return err
}
