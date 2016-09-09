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
 * io/complex128_decoder.go                               *
 *                                                        *
 * hprose complex128 decoder for Go.                      *
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
	"strconv"
)

func readLongAsComplex128(r *Reader) complex128 {
	return complex(ReadLongAsFloat64(&r.ByteReader), 0)
}

func stringToComplex128(s string) complex128 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return complex(f, 0)
}

func readInfinityAsComplex128(r *Reader) complex128 {
	return complex(readInf(&r.ByteReader), 0)
}

func readUTF8CharAsComplex128(r *Reader) complex128 {
	return stringToComplex128(byteString(readUTF8Slice(&r.ByteReader, 1)))
}

func readStringAsComplex128(r *Reader) complex128 {
	return stringToComplex128(r.ReadStringWithoutTag())
}

func readListAsComplex128(r *Reader) complex128 {
	var floatPair [2]float64
	readListAsArray(r, reflect.ValueOf(&floatPair).Elem())
	return complex(floatPair[0], floatPair[1])
}

func readRefAsComplex128(r *Reader) complex128 {
	ref := r.ReadRef()
	if str, ok := ref.(string); ok {
		return stringToComplex128(str)
	}
	if v, ok := ref.(reflect.Value); ok {
		if floatPair, ok := v.Interface().([2]float64); ok {
			return complex(floatPair[0], floatPair[1])
		}
	}
	panic(errors.New("value of type " +
		reflect.TypeOf(ref).String() +
		" cannot be converted to type complex128"))
}

var complex128Decoders = [256]func(r *Reader) complex128{
	'0':         func(r *Reader) complex128 { return 0 },
	'1':         func(r *Reader) complex128 { return 1 },
	'2':         func(r *Reader) complex128 { return 2 },
	'3':         func(r *Reader) complex128 { return 3 },
	'4':         func(r *Reader) complex128 { return 4 },
	'5':         func(r *Reader) complex128 { return 5 },
	'6':         func(r *Reader) complex128 { return 6 },
	'7':         func(r *Reader) complex128 { return 7 },
	'8':         func(r *Reader) complex128 { return 8 },
	'9':         func(r *Reader) complex128 { return 9 },
	TagNull:     func(r *Reader) complex128 { return 0 },
	TagEmpty:    func(r *Reader) complex128 { return 0 },
	TagFalse:    func(r *Reader) complex128 { return 0 },
	TagTrue:     func(r *Reader) complex128 { return 1 },
	TagNaN:      func(r *Reader) complex128 { return complex(math.NaN(), 0) },
	TagInfinity: readInfinityAsComplex128,
	TagInteger:  readLongAsComplex128,
	TagLong:     readLongAsComplex128,
	TagDouble:   func(r *Reader) complex128 { return complex(readFloat64(&r.ByteReader), 0) },
	TagUTF8Char: readUTF8CharAsComplex128,
	TagString:   readStringAsComplex128,
	TagList:     readListAsComplex128,
	TagRef:      readRefAsComplex128,
}

func complex128Decoder(r *Reader, v reflect.Value) {
	v.SetComplex(r.ReadComplex128())
}