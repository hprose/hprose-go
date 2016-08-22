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

func arrayEncoder(writer *Writer, v reflect.Value) {
	writer.SetRef(nil)
	writer.writeArray(v)
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
	case reflect.Array:
		addr := v.Pointer()
		if !writer.WriteRef(addr) {
			writer.SetRef(addr)
			writer.writeArray(e)
		}
	case reflect.Ptr, reflect.Interface:
		ptrEncoder(writer, e)
	case reflect.Map:
		addr := v.Pointer()
		if !writer.WriteRef(addr) {
			writer.SetRef(addr)
			//writer.writeMap(e)
		}
	case reflect.Slice:
		addr := v.Pointer()
		if !writer.WriteRef(addr) {
			writer.SetRef(addr)
			writer.writeSlice(e)
		}
	case reflect.String:
		str := e.String()
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
			addr := v.Pointer()
			if !writer.WriteRef(addr) {
				writer.SetRef(addr)
				writer.writeString(str, length)
			}
		}
	case reflect.Struct:
		if e.Type().PkgPath() == "big" {
			e.Interface().(Marshaler).MarshalHprose(writer)
			return
		}
		addr := v.Pointer()
		if !writer.WriteRef(addr) {
			writer.SetRef(addr)
			//writer.writeStruct(e)
		}
	default:
		writer.WriteNil()
	}
}

func sliceEncoder(writer *Writer, v reflect.Value) {
	writer.SetRef(nil)
	writer.writeSlice(v)
}

func stringEncoder(writer *Writer, v reflect.Value) {
	writer.WriteString(v.String())
}

type valueEncoder func(writer *Writer, v reflect.Value)

var valueEncoders []valueEncoder

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
		reflect.Interface:     ptrEncoder,
		reflect.Map:           nilEncoder,
		reflect.Ptr:           ptrEncoder,
		reflect.Slice:         sliceEncoder,
		reflect.String:        stringEncoder,
		reflect.Struct:        nilEncoder,
		reflect.UnsafePointer: nilEncoder,
	}
}
