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
 * io/struct_decoder.go                                   *
 *                                                        *
 * hprose struct decoder for Go.                          *
 *                                                        *
 * LastModified: Sep 10, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"math/big"
	"reflect"
	"unsafe"
)

func readDigitAsStruct(r *Reader, v reflect.Value, tag byte) {
	typ := (*reflectValue)(unsafe.Pointer(&v)).typ
	switch typ {
	case bigIntType:
		v.Set(reflect.ValueOf(*big.NewInt(int64(tag - '0'))))
	case bigRatType:
		v.Set(reflect.ValueOf(*big.NewRat(int64(tag-'0'), 1)))
	case bigFloatType:
		v.Set(reflect.ValueOf(*big.NewFloat(float64(tag - '0'))))
	default:
		castError(tag, v.Type().String())
	}
}

func readIntAsStruct(r *Reader, v reflect.Value, tag byte) {
	i := ReadInt64(&r.ByteReader, TagSemicolon)
	typ := (*reflectValue)(unsafe.Pointer(&v)).typ
	switch typ {
	case bigIntType:
		v.Set(reflect.ValueOf(*big.NewInt(i)))
	case bigRatType:
		v.Set(reflect.ValueOf(*big.NewRat(i, 1)))
	case bigFloatType:
		v.Set(reflect.ValueOf(*big.NewFloat(float64(i))))
	default:
		castError(tag, v.Type().String())
	}
}

func readLongAsStruct(r *Reader, v reflect.Value, tag byte) {
	i := byteString(readUntil(&r.ByteReader, TagSemicolon))
	typ := (*reflectValue)(unsafe.Pointer(&v)).typ
	switch typ {
	case bigIntType:
		if bi, ok := new(big.Int).SetString(i, 10); ok {
			v.Set(reflect.ValueOf(*bi))
		}
	case bigRatType:
		if br, ok := new(big.Rat).SetString(i); ok {
			v.Set(reflect.ValueOf(*br))
		}
	case bigFloatType:
		if bf, _, err := new(big.Float).Parse(i, 10); err == nil {
			v.Set(reflect.ValueOf(*bf))
		} else {
			panic(err)
		}
	default:
		castError(tag, v.Type().String())
	}
}

func readDoubleAsStruct(r *Reader, v reflect.Value, tag byte) {
	f := byteString(readUntil(&r.ByteReader, TagSemicolon))
	typ := (*reflectValue)(unsafe.Pointer(&v)).typ
	switch typ {
	case bigFloatType:
		if bf, _, err := new(big.Float).Parse(f, 10); err == nil {
			v.Set(reflect.ValueOf(*bf))
		} else {
			panic(err)
		}
	default:
		castError(tag, v.Type().String())
	}
}

func readStringAsStruct(r *Reader, v reflect.Value, tag byte) {
	str := r.ReadStringWithoutTag()
	typ := (*reflectValue)(unsafe.Pointer(&v)).typ
	switch typ {
	case bigIntType:
		if bi, ok := new(big.Int).SetString(str, 10); ok {
			v.Set(reflect.ValueOf(*bi))
		}
	case bigRatType:
		if br, ok := new(big.Rat).SetString(str); ok {
			v.Set(reflect.ValueOf(*br))
		}
	case bigFloatType:
		if bf, _, err := new(big.Float).Parse(str, 10); err == nil {
			v.Set(reflect.ValueOf(*bf))
		} else {
			panic(err)
		}
	default:
		castError(tag, v.Type().String())
	}
}

func readMapAsStruct(r *Reader, v reflect.Value, tag byte) {
}

func readStructMeta(r *Reader, v reflect.Value, tag byte) {
}

func readStructData(r *Reader, v reflect.Value, tag byte) {
}

func readRefAsStruct(r *Reader, v reflect.Value, tag byte) {
}

var structDecoders = [256]func(r *Reader, v reflect.Value, tag byte){
	TagNull:    nilDecoder,
	'0':        readDigitAsStruct,
	'1':        readDigitAsStruct,
	'2':        readDigitAsStruct,
	'3':        readDigitAsStruct,
	'4':        readDigitAsStruct,
	'5':        readDigitAsStruct,
	'6':        readDigitAsStruct,
	'7':        readDigitAsStruct,
	'8':        readDigitAsStruct,
	'9':        readDigitAsStruct,
	TagInteger: readIntAsStruct,
	TagLong:    readLongAsStruct,
	TagDouble:  readDoubleAsStruct,
	TagString:  readStringAsStruct,
	TagMap:     readMapAsStruct,
	TagClass:   readStructMeta,
	TagObject:  readStructData,
	TagRef:     readRefAsStruct,
}

func structDecoder(r *Reader, v reflect.Value, tag byte) {
	decoder := structDecoders[tag]
	if decoder != nil {
		decoder(r, v, tag)
		return
	}
	castError(tag, v.Type().String())
}
