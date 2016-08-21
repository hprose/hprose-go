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
 * io/writer_test.go                                      *
 *                                                        *
 * hprose writer test for Go.                             *
 *                                                        *
 * LastModified: Aug 20, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func TestSerializeNil(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(nil)
	if b.String() != "n" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeNil(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(nil)
	}
}

func BenchmarkWriteNil(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteNil()
	}
}

func BenchmarkNilSerialize(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		Nil.Serialize(writer, nil)
	}
}

func TestSerializeTrue(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(true)
	if b.String() != "t" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeTrue(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(true)
	}
}

func BenchmarkWriteTrue(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteBool(true)
	}
}

func BenchmarkTrueSerialize(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		Bool.Serialize(writer, true)
	}
}

func TestSerializeFalse(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(false)
	if b.String() != "f" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeFalse(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(false)
	}
}

func BenchmarkWriteFalse(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteBool(false)
	}
}

func BenchmarkFalseSerialize(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		Bool.Serialize(writer, false)
	}
}

func TestSerializeDigit(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		writer.Serialize(i)
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
}

func TestSerializeInt(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt32-10) + 10
		writer.Serialize(x)
		if b.String() != "i"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt64-math.MaxInt32-1) + math.MaxInt32 + 1
		writer.Serialize(x)
		if b.String() != "l"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
}

func TestSerializeInt8(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		writer.Serialize(int8(i))
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
	for i := 10; i <= 127; i++ {
		b.Truncate(0)
		writer.Serialize(int8(i))
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
	for i := -128; i < 0; i++ {
		b.Truncate(0)
		writer.Serialize(int8(i))
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
}

func TestSerializeInt16(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(int16(math.MaxInt16))
	if b.String() != "i"+strconv.Itoa(math.MaxInt16)+";" {
		t.Error(b.String())
	}
}

func TestSerializeInt32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(int32(math.MaxInt32))
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
}

func TestSerializeInt64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(int64(math.MaxInt32))
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(int64(math.MaxInt64))
	if b.String() != "l"+strconv.Itoa(math.MaxInt64)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt32-10) + 10
		writer.Serialize(uint(x))
		if b.String() != "i"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt64-math.MaxInt32-1) + math.MaxInt32 + 1
		writer.Serialize(uint(x))
		if b.String() != "l"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
}

func TestSerializeUint8(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		writer.Serialize(uint8(i))
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
	for i := 10; i <= 255; i++ {
		b.Truncate(0)
		writer.Serialize(uint8(i))
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
}

func TestSerializeUint16(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(uint16(math.MaxUint16))
	if b.String() != "i"+strconv.Itoa(math.MaxUint16)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(uint32(math.MaxUint32))
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(uint64(math.MaxUint32))
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(uint64(math.MaxUint64))
	if b.String() != "l"+strconv.FormatUint(math.MaxUint64, 10)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUintptr(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(uintptr(123))
	if b.String() != "i123;" {
		t.Error(b.String())
	}
}

func TestSerializeFloat32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(float32(math.NaN()))
	if b.String() != "N" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(float32(math.Inf(1)))
	if b.String() != "I+" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(float32(math.Inf(-1)))
	if b.String() != "I-" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(float32(3.14159))
	if b.String() != "d3.14159;" {
		t.Error(b.String())
	}
}

func TestSerializeFloat64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(math.NaN())
	if b.String() != "N" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(math.Inf(1))
	if b.String() != "I+" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(math.Inf(-1))
	if b.String() != "I-" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(3.14159265358979)
	if b.String() != "d3.14159265358979;" {
		t.Error(b.String())
	}
}

func TestSerializeComplex64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(complex(float32(100), 0))
	if b.String() != "d100;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(complex(0, float32(100)))
	if b.String() != "a2{d0;d100;}" {
		t.Error(b.String())
	}
}

func TestSerializeComplex128(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(complex(100, 0))
	if b.String() != "d100;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(complex(0, 100))
	if b.String() != "a2{d0;d100;}" {
		t.Error(b.String())
	}
}

func TestWriteTuple(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.WriteTuple()
	if b.String() != "a{}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.WriteTuple(1, 3.14, true)
	if b.String() != "a3{1d3.14;t}" {
		t.Error(b.String())
	}
}

func TestSerializeArray(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize([...]int{1, 2, 3})
	if b.String() != "a3{123}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([...]float64{1, 2, 3})
	if b.String() != "a3{d1;d2;d3;}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([...]int{})
	if b.String() != "a{}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([...]byte{'h', 'e', 'l', 'l', 'o'})
	if b.String() != "b5\"hello\"" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([...]byte{})
	if b.String() != "e" {
		t.Error(b.String())
	}
}

func TestSerializeSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize([]int{1, 2, 3})
	if b.String() != "a3{123}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([]float64{1, 2, 3})
	if b.String() != "a3{d1;d2;d3;}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([]int{})
	if b.String() != "a{}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([]byte{'h', 'e', 'l', 'l', 'o'})
	if b.String() != "b5\"hello\"" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize([]byte{})
	if b.String() != "e" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeInt(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(i)
	}
}

func BenchmarkSerializeArray(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	array := [...]int{0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.Serialize(array)
	}
}

func BenchmarkWriteArray(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	array := [...]int{0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.WriteArray(array)
	}
}

func BenchmarkWriteSlice(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	slice := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.WriteSlice(slice)
	}
}

func BenchmarkWriteIntSlice(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	slice := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.WriteIntSlice(slice)
	}
}
