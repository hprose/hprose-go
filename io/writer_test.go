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
 * LastModified: Aug 17, 2016                             *
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
}
