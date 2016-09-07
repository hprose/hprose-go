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
 * io/int_decoder.go                                      *
 *                                                        *
 * hprose int decoder for Go.                             *
 *                                                        *
 * LastModified: Sep 6, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"errors"
	"reflect"
	"strconv"
)

func readInt64(r *Reader) int64 {
	return ReadInt64(&r.ByteReader, TagSemicolon)
}

func readFloat64AsInt64(r *Reader) int64 {
	return int64(readFloat64(&r.ByteReader))
}

func stringToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func readUTF8CharAsInt(r *Reader) int64 {
	return stringToInt(byteString(readUTF8Slice(&r.ByteReader, 1)))
}

func readStringAsInt(r *Reader) int64 {
	return stringToInt(r.ReadStringWithoutTag())
}

func readRefAsInt(r *Reader) int64 {
	ref := r.ReadRef()
	if str, ok := ref.(string); ok {
		return stringToInt(str)
	}
	panic(errors.New("value of type " +
		reflect.TypeOf(ref).String() +
		" cannot be converted to type int"))
}

var intDecoders = [256]func(r *Reader) int64{
	'0':         func(r *Reader) int64 { return 0 },
	'1':         func(r *Reader) int64 { return 1 },
	'2':         func(r *Reader) int64 { return 2 },
	'3':         func(r *Reader) int64 { return 3 },
	'4':         func(r *Reader) int64 { return 4 },
	'5':         func(r *Reader) int64 { return 5 },
	'6':         func(r *Reader) int64 { return 6 },
	'7':         func(r *Reader) int64 { return 7 },
	'8':         func(r *Reader) int64 { return 8 },
	'9':         func(r *Reader) int64 { return 9 },
	TagNull:     func(r *Reader) int64 { return 0 },
	TagEmpty:    func(r *Reader) int64 { return 0 },
	TagFalse:    func(r *Reader) int64 { return 0 },
	TagTrue:     func(r *Reader) int64 { return 1 },
	TagInteger:  readInt64,
	TagLong:     readInt64,
	TagDouble:   readFloat64AsInt64,
	TagUTF8Char: readUTF8CharAsInt,
	TagString:   readStringAsInt,
	TagRef:      readRefAsInt,
}

func intDecoder(r *Reader, v reflect.Value) {
	v.SetInt(r.ReadInt64())
}
