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
 * LastModified: Oct 15, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"errors"
	"math"
	"reflect"
)

func readDigitAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(int(tag - '0')))
}

func readNilAsInterface(r *Reader, v reflect.Value, tag byte) {
	if v.IsNil() {
		return
	}
	v.Set(reflect.Zero(v.Type()))
}

func readEmptyAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(""))
}

func readFalseAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(false))
}

func readTrueAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(true))
}

func readNaNAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(math.NaN()))
}

func readInfAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.readInf()))
}

func readIntAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.readInt()))
}

func readLongAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.ReadBigIntWithoutTag()))
}

func readFloatAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.readFloat64()))
}

func readUTF8CharAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(readUTF8CharAsString(r)))
}

func readStringAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.ReadStringWithoutTag()))
}

func readBytesAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.ReadBytesWithoutTag()))
}

func readGUIDAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(readGUIDAsString(r)))
}

func readDateTimeAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.ReadDateTimeWithoutTag()))
}

func readTimeAsInterface(r *Reader, v reflect.Value, tag byte) {
	v.Set(reflect.ValueOf(r.ReadTimeWithoutTag()))
}

func readListAsInterface(r *Reader, v reflect.Value, tag byte) {
	var slice []interface{}
	sv := reflect.ValueOf(&slice).Elem()
	readListAsSlice(r, sv, TagList)
	v.Set(sv)
}

func readMapAsInterface(r *Reader, v reflect.Value, tag byte) {
	var mv reflect.Value
	if r.JSONCompatible {
		var m map[string]interface{}
		mv = reflect.ValueOf(&m).Elem()
	} else {
		var m map[interface{}]interface{}
		mv = reflect.ValueOf(&m).Elem()
	}
	readMap(r, mv, TagMap)
	v.Set(mv)
}

func readRefAsInterface(r *Reader, v reflect.Value, tag byte) {
	iv := reflect.ValueOf(r.readRef())
	t := v.Type()
	it := iv.Type()
	if it.AssignableTo(t) {
		v.Set(iv)
	} else if it.ConvertibleTo(t) {
		v.Set(iv.Convert(t))
	} else {
		panic(errors.New(it.String() +
			" cannot be converted to type" +
			t.String()))
	}
}

var interfaceDecoders = [256]func(r *Reader, v reflect.Value, tag byte){
	'0':         readDigitAsInterface,
	'1':         readDigitAsInterface,
	'2':         readDigitAsInterface,
	'3':         readDigitAsInterface,
	'4':         readDigitAsInterface,
	'5':         readDigitAsInterface,
	'6':         readDigitAsInterface,
	'7':         readDigitAsInterface,
	'8':         readDigitAsInterface,
	'9':         readDigitAsInterface,
	TagNull:     readNilAsInterface,
	TagEmpty:    readEmptyAsInterface,
	TagFalse:    readFalseAsInterface,
	TagTrue:     readTrueAsInterface,
	TagNaN:      readNaNAsInterface,
	TagInfinity: readInfAsInterface,
	TagInteger:  readIntAsInterface,
	TagLong:     readLongAsInterface,
	TagDouble:   readFloatAsInterface,
	TagUTF8Char: readUTF8CharAsInterface,
	TagString:   readStringAsInterface,
	TagBytes:    readBytesAsInterface,
	TagGUID:     readGUIDAsInterface,
	TagDate:     readDateTimeAsInterface,
	TagTime:     readTimeAsInterface,
	TagList:     readListAsInterface,
	TagMap:      readMapAsInterface,
	TagClass:    readStructMeta,
	TagObject:   readStructData,
	TagRef:      readRefAsInterface,
}

func interfaceDecoder(r *Reader, v reflect.Value, tag byte) {
	decoder := interfaceDecoders[tag]
	if decoder != nil {
		decoder(r, v, tag)
		return
	}
	castError(tag, "interface{}")
}
