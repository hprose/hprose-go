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
 * LastModified: Aug 18, 2016                             *
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
*/

func testSerializeNil(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(nil)
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "n" {
		t.Error(b.String())
	}
}
func testSerializeTrue(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(true)
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "t" {
		t.Error(b.String())
	}
}
func testSerializeFalse(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(false)
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "f" {
		t.Error(b.String())
	}
}
func testSerializeDigit(t *testing.T, writer *Writer, b *bytes.Buffer) {
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		err := writer.Serialize(i)
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
}

func testSerializeInt(t *testing.T, writer *Writer, b *bytes.Buffer) {
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt32-10) + 10
		err := writer.Serialize(x)
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != "i"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt64-math.MaxInt32-1) + math.MaxInt32 + 1
		err := writer.Serialize(x)
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != "l"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
}

func testSerializeInt8(t *testing.T, writer *Writer, b *bytes.Buffer) {
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		err := writer.Serialize(int8(i))
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
	for i := 10; i <= 127; i++ {
		b.Truncate(0)
		err := writer.Serialize(int8(i))
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
	for i := -128; i < 0; i++ {
		b.Truncate(0)
		err := writer.Serialize(int8(i))
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
}

func testSerializeInt16(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(int16(math.MaxInt16))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "i"+strconv.Itoa(math.MaxInt16)+";" {
		t.Error(b.String())
	}
}

func testSerializeInt32(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(int32(math.MaxInt32))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
}

func testSerializeInt64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(int64(math.MaxInt32))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(int64(math.MaxInt64))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "l"+strconv.Itoa(math.MaxInt64)+";" {
		t.Error(b.String())
	}
}

func testSerializeUint(t *testing.T, writer *Writer, b *bytes.Buffer) {
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt32-10) + 10
		err := writer.Serialize(uint(x))
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != "i"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt64-math.MaxInt32-1) + math.MaxInt32 + 1
		err := writer.Serialize(uint(x))
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != "l"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
}

func testSerializeUint8(t *testing.T, writer *Writer, b *bytes.Buffer) {
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		err := writer.Serialize(uint8(i))
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
	for i := 10; i <= 255; i++ {
		b.Truncate(0)
		err := writer.Serialize(uint8(i))
		if err != nil {
			t.Error(err.Error())
		}
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
}

func testSerializeUint16(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(uint16(math.MaxUint16))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "i"+strconv.Itoa(math.MaxUint16)+";" {
		t.Error(b.String())
	}
}

func testSerializeUint32(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(uint32(math.MaxUint32))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
}

func testSerializeUint64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(uint64(math.MaxUint32))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(uint64(math.MaxUint64))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "l"+strconv.FormatUint(math.MaxUint64, 10)+";" {
		t.Error(b.String())
	}
}

func testSerializeFloat32(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(float32(math.NaN()))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "N" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(float32(math.Inf(1)))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "I+" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(float32(math.Inf(-1)))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "I-" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(float32(3.14159))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "d3.14159;" {
		t.Error(b.String())
	}
}

func testSerializeFloat64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(math.NaN())
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "N" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(math.Inf(1))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "I+" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(math.Inf(-1))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "I-" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(3.14159265358979)
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "d3.14159265358979;" {
		t.Error(b.String())
	}
}

func testSerializeComplex64(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(complex(float32(100), 0))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "d100;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(complex(0, float32(100)))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "a2{d0;d100;}" {
		t.Error(b.String())
	}
}

func testSerializeComplex128(t *testing.T, writer *Writer, b *bytes.Buffer) {
	b.Truncate(0)
	err := writer.Serialize(complex(100, 0))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "d100;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	err = writer.Serialize(complex(0, 100))
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "a2{d0;d100;}" {
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
	testSerializeFloat32(t, writer, b)
	testSerializeFloat64(t, writer, b)
	testSerializeComplex64(t, writer, b)
	testSerializeComplex128(t, writer, b)
}
