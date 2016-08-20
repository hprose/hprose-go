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

/*
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
*/

func testSerializeNil(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(nil)
	if b.String() != "n" {
		t.Error(b.String())
	}
}
func testSerializeTrue(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(true)
	if b.String() != "t" {
		t.Error(b.String())
	}
}
func testSerializeFalse(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(false)
	if b.String() != "f" {
		t.Error(b.String())
	}
}
func testSerializeDigit(t *testing.T, writer *Writer, b *bytes.Buffer) {
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		writer.Serialize(i)
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
}

func testSerializeInt(t *testing.T, writer *Writer, b *bytes.Buffer) {
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

func testSerializeInt8(t *testing.T, writer *Writer, b *bytes.Buffer) {
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

func testSerializeInt16(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(int16(math.MaxInt16))
	if b.String() != "i"+strconv.Itoa(math.MaxInt16)+";" {
		t.Error(b.String())
	}
}

func testSerializeInt32(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(int32(math.MaxInt32))
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
}

func testSerializeInt64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testSerializeUint(t *testing.T, writer *Writer, b *bytes.Buffer) {
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

func testSerializeUint8(t *testing.T, writer *Writer, b *bytes.Buffer) {
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

func testSerializeUint16(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(uint16(math.MaxUint16))
	if b.String() != "i"+strconv.Itoa(math.MaxUint16)+";" {
		t.Error(b.String())
	}
}

func testSerializeUint32(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(uint32(math.MaxUint32))
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
}

func testSerializeUint64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testSerializeUintptr(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	writer.Serialize(uintptr(123))
	if b.String() != "i123;" {
		t.Error(b.String())
	}
}

func testSerializeFloat32(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testSerializeFloat64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testSerializeComplex64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testSerializeComplex128(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testWriteTuple(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testSerializeArray(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func testSerializeSlice(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
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

func TestSerialize(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	testSerializeNil(t, writer, b)
	testSerializeTrue(t, writer, b)
	testSerializeFalse(t, writer, b)
	testSerializeDigit(t, writer, b)
	testSerializeInt(t, writer, b)
	testSerializeInt8(t, writer, b)
	testSerializeInt16(t, writer, b)
	testSerializeInt32(t, writer, b)
	testSerializeInt64(t, writer, b)
	testSerializeUint(t, writer, b)
	testSerializeUint8(t, writer, b)
	testSerializeUint16(t, writer, b)
	testSerializeUint32(t, writer, b)
	testSerializeUint64(t, writer, b)
	testSerializeUintptr(t, writer, b)
	testSerializeFloat32(t, writer, b)
	testSerializeFloat64(t, writer, b)
	testSerializeComplex64(t, writer, b)
	testSerializeComplex128(t, writer, b)
	testWriteTuple(t, writer, b)
	testSerializeArray(t, writer, b)
	testSerializeSlice(t, writer, b)
}
