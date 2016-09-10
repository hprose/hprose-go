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
 * io/map_decoder.go                                      *
 *                                                        *
 * hprose map decoder for Go.                             *
 *                                                        *
 * LastModified: Sep 9, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"errors"
	"reflect"
	"strconv"
)

func setIntKey(kind reflect.Kind, k reflect.Value, i int) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		k.SetInt(int64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		k.SetUint(uint64(i))
	case reflect.Float32, reflect.Float64:
		k.SetFloat(float64(i))
	case reflect.String:
		k.SetString(strconv.Itoa(i))
	case reflect.Interface:
		k.Set(reflect.ValueOf(i))
	default:
		panic(errors.New("cannot convert int to type " + k.Type().String()))
	}
}

func readListAsMap(r *Reader, v reflect.Value, tag byte) {
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	l := readCount(&r.ByteReader)
	if !r.Simple {
		setReaderRef(r, v)
	}
	t := v.Type()
	kt := t.Key()
	vt := t.Elem()
	for i := 0; i < l; i++ {
		key := reflect.New(kt).Elem()
		setIntKey(kt.Kind(), key, i)
		val := reflect.New(vt).Elem()
		r.ReadValue(val)
		v.SetMapIndex(key, val)
	}
	r.readByte()
}

func readMap(r *Reader, v reflect.Value, tag byte) {
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	l := readCount(&r.ByteReader)
	if !r.Simple {
		setReaderRef(r, v)
	}
	t := v.Type()
	kt := t.Key()
	vt := t.Elem()
	for i := 0; i < l; i++ {
		key := reflect.New(kt).Elem()
		r.ReadValue(key)
		val := reflect.New(vt).Elem()
		r.ReadValue(val)
		v.SetMapIndex(key, val)
	}
	r.readByte()
}

func readRefAsMap(r *Reader, v reflect.Value, tag byte) {
	ref := r.ReadRef()
	if m, ok := ref.(reflect.Value); ok {
		if m.Kind() == reflect.Map {
			v.Set(m)
			return
		}
	}
	panic(errors.New("value of type " +
		reflect.TypeOf(ref).String() +
		" cannot be converted to type map"))
}

var mapDecoders = [256]func(r *Reader, v reflect.Value, tag byte){
	TagNull:  nilDecoder,
	TagEmpty: nilDecoder,
	TagList:  readListAsMap,
	TagMap:   readMap,
	TagRef:   readRefAsMap,
}

func mapDecoder(r *Reader, v reflect.Value, tag byte) {
	decoder := mapDecoders[tag]
	if decoder != nil {
		decoder(r, v, tag)
		return
	}
	castError(tag, v.Type().String())
}
