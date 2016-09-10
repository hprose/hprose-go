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
	"strconv"
	"time"
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
	case timeType:
		v.Set(reflect.ValueOf(time.Unix(int64(tag-'0'), 0)))
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
	case timeType:
		v.Set(reflect.ValueOf(time.Unix(i, 0)))
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
	case timeType:
		if unix, err := strconv.ParseInt(i, 10, 64); err == nil {
			v.Set(reflect.ValueOf(time.Unix(unix, 0)))
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
	case timeType:
		if unix, err := strconv.ParseFloat(f, 10); err == nil {
			sec := int64(unix)
			nsec := int64((unix - float64(sec)) * 1000000000)
			v.Set(reflect.ValueOf(time.Unix(sec, nsec)))
		} else {
			panic(err)
		}
	default:
		castError(tag, v.Type().String())
	}
}

const timeStringFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

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
	case timeType:
		if t, err := time.Parse(timeStringFormat, str); err == nil {
			v.Set(reflect.ValueOf(t))
		} else {
			panic(err)
		}
	default:
		castError(tag, v.Type().String())
	}
}

func readTimeAsStruct(t time.Time, v reflect.Value, tag byte) {
	typ := (*reflectValue)(unsafe.Pointer(&v)).typ
	switch typ {
	case bigIntType:
		v.Set(reflect.ValueOf(*new(big.Int).SetInt64(t.Unix())))
	case bigRatType:
		v.Set(reflect.ValueOf(*new(big.Rat).SetInt64(t.Unix())))
	case bigFloatType:
		ft := float64(t.Unix()) + float64(t.Nanosecond())/1000000000
		v.Set(reflect.ValueOf(*new(big.Float).SetFloat64(ft)))
	case timeType:
		v.Set(reflect.ValueOf(t))
	default:
		castError(tag, v.Type().String())
	}
}
func readDateTimeStruct(r *Reader, v reflect.Value, tag byte) {
	readTimeAsStruct(r.ReadDateTimeWithoutTag(), v, tag)
}

func readTimeStruct(r *Reader, v reflect.Value, tag byte) {
	readTimeAsStruct(r.ReadTimeWithoutTag(), v, tag)
}

func readListAsStruct(r *Reader, v reflect.Value, tag byte) {
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
	TagDate:    readDateTimeStruct,
	TagTime:    readTimeStruct,
	TagList:    readListAsStruct,
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
