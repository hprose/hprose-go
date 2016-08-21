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

import "reflect"

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
	writer.writeArray(v)
}

func ptrEncoder(writer *Writer, v reflect.Value) {
	if v.IsNil() {
		writer.WriteNil()
	} else {
		writer.WriteValue(v.Elem())
	}
}

func sliceEncoder(writer *Writer, v reflect.Value) {
	writer.writeSlice(v)
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
		reflect.String:        nilEncoder,
		reflect.Struct:        nilEncoder,
		reflect.UnsafePointer: nilEncoder,
	}
}
