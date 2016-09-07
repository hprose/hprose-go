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
 * LastModified: Sep 7, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"errors"
	"reflect"
	"strconv"
)

func readUint64(r *Reader) uint64 {
	return ReadUint64(&r.ByteReader, TagSemicolon)
}

func readFloat64AsUint64(r *Reader) uint64 {
	return uint64(readFloat64(&r.ByteReader))
}

func stringToUint(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func readUTF8CharAsUint(r *Reader) uint64 {
	return stringToUint(byteString(readUTF8Slice(&r.ByteReader, 1)))
}

func readStringAsUint(r *Reader) uint64 {
	return stringToUint(r.ReadStringWithoutTag())
}

func readRefAsUint(r *Reader) uint64 {
	ref := r.ReadRef()
	if str, ok := ref.(string); ok {
		return stringToUint(str)
	}
	panic(errors.New("value of type " +
		reflect.TypeOf(ref).String() +
		" cannot be converted to type uint"))
}

var uintDecoders = [256]func(r *Reader) uint64{
	'0':         func(r *Reader) uint64 { return 0 },
	'1':         func(r *Reader) uint64 { return 1 },
	'2':         func(r *Reader) uint64 { return 2 },
	'3':         func(r *Reader) uint64 { return 3 },
	'4':         func(r *Reader) uint64 { return 4 },
	'5':         func(r *Reader) uint64 { return 5 },
	'6':         func(r *Reader) uint64 { return 6 },
	'7':         func(r *Reader) uint64 { return 7 },
	'8':         func(r *Reader) uint64 { return 8 },
	'9':         func(r *Reader) uint64 { return 9 },
	TagNull:     func(r *Reader) uint64 { return 0 },
	TagEmpty:    func(r *Reader) uint64 { return 0 },
	TagFalse:    func(r *Reader) uint64 { return 0 },
	TagTrue:     func(r *Reader) uint64 { return 1 },
	TagInteger:  readUint64,
	TagLong:     readUint64,
	TagDouble:   readFloat64AsUint64,
	TagUTF8Char: readUTF8CharAsUint,
	TagString:   readStringAsUint,
	TagRef:      readRefAsUint,
}

func uintDecoder(r *Reader, v reflect.Value) {
	v.SetUint(r.ReadUint64())
}
