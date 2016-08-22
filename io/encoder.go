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
 * io/encoder.go                                          *
 *                                                        *
 * hprose encoder for Go.                                 *
 *                                                        *
 * LastModified: Aug 22, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"reflect"
	"unsafe"

	"github.com/hprose/hprose-golang/util"
)

type valueEncoder func(writer *Writer, v reflect.Value)

var valueEncoders []valueEncoder

func nilEncoder(writer *Writer, v reflect.Value) {
	writer.WriteNil()
}

func boolEncoder(writer *Writer, v reflect.Value) {
	writer.WriteBool(v.Bool())
}

func intEncoder(writer *Writer, v reflect.Value) {
	writer.WriteInt(v.Int())
}

func uintEncoder(writer *Writer, v reflect.Value) {
	writer.WriteUint(v.Uint())
}

func float32Encoder(writer *Writer, v reflect.Value) {
	writer.WriteFloat(v.Float(), 32)
}

func float64Encoder(writer *Writer, v reflect.Value) {
	writer.WriteFloat(v.Float(), 64)
}

func complex64Encoder(writer *Writer, v reflect.Value) {
	writer.WriteComplex64(complex64(v.Complex()))
}

func complex128Encoder(writer *Writer, v reflect.Value) {
	writer.WriteComplex128(v.Complex())
}

func interfaceEncoder(writer *Writer, v reflect.Value) {
	if v.IsNil() {
		writer.WriteNil()
		return
	}
	e := v.Elem()
	valueEncoders[e.Kind()](writer, e)
}

func arrayEncoder(writer *Writer, v reflect.Value) {
	writer.SetRef(nil)
	writer.writeArray(v)
}

func sliceEncoder(writer *Writer, v reflect.Value) {
	writer.SetRef(nil)
	writer.writeSlice(v)
}

func stringEncoder(writer *Writer, v reflect.Value) {
	writer.WriteString(v.String())
}

func arrayPtrEncoder(writer *Writer, v reflect.Value, addr uintptr) {
	if !writer.WriteRef(addr) {
		writer.SetRef(addr)
		writer.writeArray(v)
	}
}

func mapPtrEncoder(writer *Writer, v reflect.Value, addr uintptr) {
	if !writer.WriteRef(addr) {
		writer.SetRef(addr)
		//writer.writeMap(v)
	}
}

func slicePtrEncoder(writer *Writer, v reflect.Value, addr uintptr) {
	if !writer.WriteRef(addr) {
		writer.SetRef(addr)
		writer.writeSlice(v)
	}
}

func stringPtrEncoder(writer *Writer, v reflect.Value, addr uintptr) {
	str := v.String()
	length := util.UTF16Length(str)
	switch {
	case length == 0:
		writer.Stream.WriteByte(TagEmpty)
	case length < 0:
		writer.WriteBytes(*(*[]byte)(unsafe.Pointer(&str)))
	case length == 1:
		writer.Stream.WriteByte(TagUTF8Char)
		writer.Stream.WriteString(str)
	default:
		if !writer.WriteRef(addr) {
			writer.SetRef(addr)
			writer.writeString(str, length)
		}
	}
}

func structPtrEncoder(writer *Writer, v reflect.Value, addr uintptr) {
	if v.Type().PkgPath() == "big" {
		v.Interface().(Marshaler).MarshalHprose(writer)
		return
	}
	if !writer.WriteRef(addr) {
		writer.SetRef(addr)
		//writer.writeStruct(v)
	}
}

func ptrEncoder(writer *Writer, v reflect.Value) {
	if v.IsNil() {
		writer.WriteNil()
		return
	}
	e := v.Elem()
	switch e.Kind() {
	case reflect.Bool:
		writer.WriteBool(e.Bool())
	case reflect.Int, reflect.Int64:
		writer.WriteInt(e.Int())
	case reflect.Int8, reflect.Int16, reflect.Int32:
		writer.WriteInt32(int32(e.Int()))
	case reflect.Uint8, reflect.Uint16:
		writer.WriteInt32(int32(e.Uint()))
	case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		writer.WriteUint(e.Uint())
	case reflect.Float32:
		writer.WriteFloat(e.Float(), 32)
	case reflect.Float64:
		writer.WriteFloat(e.Float(), 64)
	case reflect.Complex64:
		writer.WriteComplex64(complex64(v.Complex()))
	case reflect.Complex128:
		writer.WriteComplex128(v.Complex())
	case reflect.Ptr:
		ptrEncoder(writer, e)
	case reflect.Interface:
		interfaceEncoder(writer, e)
	case reflect.Array:
		arrayPtrEncoder(writer, e, v.Pointer())
	case reflect.Map:
		mapPtrEncoder(writer, e, v.Pointer())
	case reflect.Slice:
		slicePtrEncoder(writer, e, v.Pointer())
	case reflect.String:
		stringPtrEncoder(writer, e, v.Pointer())
	case reflect.Struct:
		structPtrEncoder(writer, e, v.Pointer())
	default:
		writer.WriteNil()
	}
}

func init() {
	valueEncoders = []valueEncoder{
		reflect.Invalid:       nilEncoder,
		reflect.Bool:          boolEncoder,
		reflect.Int:           intEncoder,
		reflect.Int8:          intEncoder,
		reflect.Int16:         intEncoder,
		reflect.Int32:         intEncoder,
		reflect.Int64:         intEncoder,
		reflect.Uint:          uintEncoder,
		reflect.Uint8:         uintEncoder,
		reflect.Uint16:        uintEncoder,
		reflect.Uint32:        uintEncoder,
		reflect.Uint64:        uintEncoder,
		reflect.Uintptr:       uintEncoder,
		reflect.Float32:       float32Encoder,
		reflect.Float64:       float64Encoder,
		reflect.Complex64:     complex64Encoder,
		reflect.Complex128:    complex128Encoder,
		reflect.Array:         arrayEncoder,
		reflect.Chan:          nilEncoder,
		reflect.Func:          nilEncoder,
		reflect.Interface:     interfaceEncoder,
		reflect.Map:           nilEncoder,
		reflect.Ptr:           ptrEncoder,
		reflect.Slice:         sliceEncoder,
		reflect.String:        stringEncoder,
		reflect.Struct:        nilEncoder,
		reflect.UnsafePointer: nilEncoder,
	}
}
