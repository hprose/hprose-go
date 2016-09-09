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
 * io/interface_decoder.go                                *
 *                                                        *
 * hprose interface decoder for Go.                       *
 *                                                        *
 * LastModified: Sep 9, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"errors"
	"math"
	"reflect"
)

func readInterfaceSlice(r *Reader) interface{} {
	var slice []interface{}
	readListAsSlice(r, reflect.ValueOf(&slice).Elem())
	return slice
}

func readInterfaceMap(r *Reader) interface{} {
	if r.JSONCompatible {
		var m map[string]interface{}
		readMap(r, reflect.ValueOf(&m).Elem())
		return m
	}
	var m map[interface{}]interface{}
	readMap(r, reflect.ValueOf(&m).Elem())
	return m
}

var interfaceDecoders = [256]func(r *Reader) interface{}{
	'0':         func(r *Reader) interface{} { return 0 },
	'1':         func(r *Reader) interface{} { return 1 },
	'2':         func(r *Reader) interface{} { return 2 },
	'3':         func(r *Reader) interface{} { return 3 },
	'4':         func(r *Reader) interface{} { return 4 },
	'5':         func(r *Reader) interface{} { return 5 },
	'6':         func(r *Reader) interface{} { return 6 },
	'7':         func(r *Reader) interface{} { return 7 },
	'8':         func(r *Reader) interface{} { return 8 },
	'9':         func(r *Reader) interface{} { return 9 },
	TagNull:     func(r *Reader) interface{} { return nil },
	TagEmpty:    func(r *Reader) interface{} { return "" },
	TagFalse:    func(r *Reader) interface{} { return false },
	TagTrue:     func(r *Reader) interface{} { return true },
	TagNaN:      func(r *Reader) interface{} { return math.NaN() },
	TagInfinity: func(r *Reader) interface{} { return readInf(&r.ByteReader) },
	TagInteger:  func(r *Reader) interface{} { return readInt(&r.ByteReader) },
	TagLong:     func(r *Reader) interface{} { return r.ReadBigIntWithoutTag() },
	TagDouble:   func(r *Reader) interface{} { return readFloat64(&r.ByteReader) },
	TagUTF8Char: func(r *Reader) interface{} { return readUTF8CharAsString(r) },
	TagString:   func(r *Reader) interface{} { return r.ReadStringWithoutTag() },
	TagBytes:    func(r *Reader) interface{} { return r.ReadBytesWithoutTag() },
	TagGUID:     func(r *Reader) interface{} { return readGUIDAsString(r) },
	TagDate:     func(r *Reader) interface{} { return r.ReadDateTimeWithoutTag() },
	TagTime:     func(r *Reader) interface{} { return r.ReadTimeWithoutTag() },
	TagList:     readInterfaceSlice,
	TagMap:      readInterfaceMap,
	TagClass:    func(r *Reader) interface{} { panic("TODO") },
	TagObject:   func(r *Reader) interface{} { panic("TODO") },
	TagRef:      func(r *Reader) interface{} { return r.ReadRef() },
}

func interfaceDecoder(r *Reader, v reflect.Value) {
	p := r.ReadInterface()
	t := v.Type()
	rv := reflect.ValueOf(&p).Elem()
	rt := rv.Type()
	if rt.AssignableTo(t) {
		v.Set(rv)
	} else if rt.ConvertibleTo(t) {
		v.Set(rv.Convert(t))
	} else {
		panic(errors.New(rt.String() +
			" cannot be converted to type" +
			t.String()))
	}
}
