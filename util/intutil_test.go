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
 * util/intutil_test.go                                   *
 *                                                        *
 * intutil test for Go.                                   *
 *                                                        *
 * LastModified: Aug 17, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package util

import (
	"math"
	"reflect"
	"strconv"
	"testing"
)

func BenchmarkGetInt32Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetInt32Bytes(int32(i))
		GetInt32Bytes(math.MaxInt32 - int32(i))
		GetInt32Bytes(int32(-i))
		GetInt32Bytes(math.MinInt32 + int32(i))
	}
}

func BenchmarkGetInt64Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetInt64Bytes(int64(i))
		GetInt64Bytes(math.MaxInt64 - int64(i))
		GetInt64Bytes(int64(-i))
		GetInt64Bytes(math.MinInt64 + int64(i))
	}
}

func BenchmarkGetUint32Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetUint32Bytes(uint32(i))
		GetUint32Bytes(math.MaxUint32 - uint32(i))
		GetUint32Bytes(uint32(-i))
		GetUint32Bytes(math.MaxUint32 + uint32(i))
	}
}

func BenchmarkGetUint64Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetUint64Bytes(uint64(i))
		GetUint64Bytes(math.MaxUint64 - uint64(i))
		GetUint64Bytes(uint64(-i))
		GetUint64Bytes(math.MaxUint64 + uint64(i))
	}
}

func BenchmarkItoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
		strconv.Itoa(math.MaxInt64 - i)
		strconv.Itoa(-i)
		strconv.Itoa(math.MaxInt64 + i)
	}
}

func BenchmarkGetInt32BytesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i int32
		for pb.Next() {
			GetInt32Bytes(i)
			GetInt32Bytes(math.MaxInt32 - i)
			GetInt32Bytes(-i)
			GetInt32Bytes(math.MinInt32 + i)
			i++
		}
	})
}

func BenchmarkGetInt64BytesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i int64
		for pb.Next() {
			GetInt64Bytes(i)
			GetInt64Bytes(math.MaxInt64 - i)
			GetInt64Bytes(-i)
			GetInt64Bytes(math.MinInt64 + i)
			i++
		}
	})
}

func BenchmarkGetUint32BytesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i uint32
		for pb.Next() {
			GetUint32Bytes(i)
			GetUint32Bytes(math.MaxUint32 - i)
			GetUint32Bytes(-i)
			GetUint32Bytes(math.MaxUint32 + i)
			i++
		}
	})
}

func BenchmarkGetUint64BytesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i uint64
		for pb.Next() {
			GetUint64Bytes(i)
			GetUint64Bytes(math.MaxUint64 - i)
			GetUint64Bytes(-i)
			GetUint64Bytes(math.MaxUint64 + i)
			i++
		}
	})
}

func BenchmarkItoaParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			strconv.Itoa(i)
			strconv.Itoa(math.MaxInt64 - i)
			strconv.Itoa(-i)
			strconv.Itoa(math.MinInt64 + i)
			i++
		}
	})
}

func TestGetInt32Bytes(t *testing.T) {
	b := GetInt32Bytes(0)
	if !reflect.DeepEqual(b, []byte{'0'}) {
		t.Error("b must be []byte{'0'}")
	}
	b = GetInt32Bytes(9)
	if !reflect.DeepEqual(b, []byte{'9'}) {
		t.Error("b must be []byte{'9'}")
	}
	b = GetInt32Bytes(10)
	if !reflect.DeepEqual(b, []byte{'1', '0'}) {
		t.Error("b must be []byte{'1', '0'}")
	}
	b = GetInt32Bytes(99)
	if !reflect.DeepEqual(b, []byte{'9', '9'}) {
		t.Error("b must be []byte{'9', '9'}")
	}
	b = GetInt32Bytes(100)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0'}")
	}
	b = GetInt32Bytes(999)
	if !reflect.DeepEqual(b, []byte{'9', '9', '9'}) {
		t.Error("b must be []byte{'9', '9', '9'}")
	}
	b = GetInt32Bytes(1000)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0', '0'}")
	}
	b = GetInt32Bytes(-1000)
	if !reflect.DeepEqual(b, []byte{'-', '1', '0', '0', '0'}) {
		t.Error("b must be []byte{'-', '1', '0', '0', '0'}")
	}
	b = GetInt32Bytes(10000)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0', '0', '0'}")
	}
	b = GetInt32Bytes(-10000)
	if !reflect.DeepEqual(b, []byte{'-', '1', '0', '0', '0', '0'}) {
		t.Error("b must be []byte{'-', '1', '0', '0', '0', '0'}")
	}
	b = GetInt32Bytes(123456789)
	if !reflect.DeepEqual(b, []byte("123456789")) {
		t.Error("b must be []byte(\"123456789\")")
	}
	b = GetInt32Bytes(-123456789)
	if !reflect.DeepEqual(b, []byte("-123456789")) {
		t.Error("b must be []byte(\"-123456789\")")
	}
	b = GetInt32Bytes(math.MaxInt32)
	if !reflect.DeepEqual(b, []byte("2147483647")) {
		t.Error("b must be []byte(\"2147483647\")")
	}
	b = GetInt32Bytes(math.MinInt32)
	if !reflect.DeepEqual(b, []byte("-2147483648")) {
		t.Error("b must be []byte(\"-2147483648\")")
	}
}

func TestGetInt64Bytes(t *testing.T) {
	b := GetInt64Bytes(0)
	if !reflect.DeepEqual(b, []byte{'0'}) {
		t.Error("b must be []byte{'0'}")
	}
	b = GetInt64Bytes(9)
	if !reflect.DeepEqual(b, []byte{'9'}) {
		t.Error("b must be []byte{'9'}")
	}
	b = GetInt64Bytes(10)
	if !reflect.DeepEqual(b, []byte{'1', '0'}) {
		t.Error("b must be []byte{'1', '0'}")
	}
	b = GetInt64Bytes(99)
	if !reflect.DeepEqual(b, []byte{'9', '9'}) {
		t.Error("b must be []byte{'9', '9'}")
	}
	b = GetInt64Bytes(100)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0'}")
	}
	b = GetInt64Bytes(999)
	if !reflect.DeepEqual(b, []byte{'9', '9', '9'}) {
		t.Error("b must be []byte{'9', '9', '9'}")
	}
	b = GetInt64Bytes(1000)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0', '0'}")
	}
	b = GetInt64Bytes(-1000)
	if !reflect.DeepEqual(b, []byte{'-', '1', '0', '0', '0'}) {
		t.Error("b must be []byte{'-', '1', '0', '0', '0'}")
	}
	b = GetInt64Bytes(10000)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0', '0', '0'}")
	}
	b = GetInt64Bytes(-10000)
	if !reflect.DeepEqual(b, []byte{'-', '1', '0', '0', '0', '0'}) {
		t.Error("b must be []byte{'-', '1', '0', '0', '0', '0'}")
	}
	b = GetInt64Bytes(123456789)
	if !reflect.DeepEqual(b, []byte("123456789")) {
		t.Error("b must be []byte(\"123456789\")")
	}
	b = GetInt64Bytes(-123456789)
	if !reflect.DeepEqual(b, []byte("-123456789")) {
		t.Error("b must be []byte(\"-123456789\")")
	}
	b = GetInt64Bytes(math.MaxInt32)
	if !reflect.DeepEqual(b, []byte("2147483647")) {
		t.Error("b must be []byte(\"2147483647\")")
	}
	b = GetInt64Bytes(math.MinInt32)
	if !reflect.DeepEqual(b, []byte("-2147483648")) {
		t.Error("b must be []byte(\"-2147483648\")")
	}
	b = GetInt64Bytes(math.MaxInt64)
	if !reflect.DeepEqual(b, []byte(strconv.Itoa(math.MaxInt64))) {
		t.Error("b must be []byte(\"" + strconv.Itoa(math.MaxInt64) + "\")")
	}
	b = GetInt64Bytes(math.MinInt64)
	if !reflect.DeepEqual(b, []byte(strconv.Itoa(math.MinInt64))) {
		t.Error("b must be []byte(\"" + strconv.Itoa(math.MinInt64) + "\")")
	}
}

func TestGetUint32Bytes(t *testing.T) {
	b := GetUint32Bytes(0)
	if !reflect.DeepEqual(b, []byte{'0'}) {
		t.Error("b must be []byte{'0'}")
	}
	b = GetUint32Bytes(9)
	if !reflect.DeepEqual(b, []byte{'9'}) {
		t.Error("b must be []byte{'9'}")
	}
	b = GetUint32Bytes(10)
	if !reflect.DeepEqual(b, []byte{'1', '0'}) {
		t.Error("b must be []byte{'1', '0'}")
	}
	b = GetUint32Bytes(99)
	if !reflect.DeepEqual(b, []byte{'9', '9'}) {
		t.Error("b must be []byte{'9', '9'}")
	}
	b = GetUint32Bytes(100)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0'}")
	}
	b = GetUint32Bytes(999)
	if !reflect.DeepEqual(b, []byte{'9', '9', '9'}) {
		t.Error("b must be []byte{'9', '9', '9'}")
	}
	b = GetUint32Bytes(1000)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0', '0'}")
	}
	b = GetUint32Bytes(10000)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0', '0', '0'}")
	}
	b = GetUint32Bytes(123456789)
	if !reflect.DeepEqual(b, []byte("123456789")) {
		t.Error("b must be []byte(\"123456789\")")
	}
	b = GetUint32Bytes(math.MaxInt32)
	if !reflect.DeepEqual(b, []byte("2147483647")) {
		t.Error("b must be []byte(\"2147483647\")")
	}
	b = GetUint32Bytes(math.MaxUint32)
	if !reflect.DeepEqual(b, []byte("4294967295")) {
		t.Error("b must be []byte(\"4294967295\")")
	}
}

func TestGetUint64Bytes(t *testing.T) {
	b := GetUint64Bytes(0)
	if !reflect.DeepEqual(b, []byte{'0'}) {
		t.Error("b must be []byte{'0'}")
	}
	b = GetUint64Bytes(9)
	if !reflect.DeepEqual(b, []byte{'9'}) {
		t.Error("b must be []byte{'9'}")
	}
	b = GetUint64Bytes(10)
	if !reflect.DeepEqual(b, []byte{'1', '0'}) {
		t.Error("b must be []byte{'1', '0'}")
	}
	b = GetUint64Bytes(99)
	if !reflect.DeepEqual(b, []byte{'9', '9'}) {
		t.Error("b must be []byte{'9', '9'}")
	}
	b = GetUint64Bytes(100)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0'}")
	}
	b = GetUint64Bytes(999)
	if !reflect.DeepEqual(b, []byte{'9', '9', '9'}) {
		t.Error("b must be []byte{'9', '9', '9'}")
	}
	b = GetUint64Bytes(1000)
	if !reflect.DeepEqual(b, []byte{'1', '0', '0', '0'}) {
		t.Error("b must be []byte{'1', '0', '0', '0'}")
	}
	b = GetUint64Bytes(123456789)
	if !reflect.DeepEqual(b, []byte("123456789")) {
		t.Error("b must be []byte(\"123456789\")")
	}
	b = GetUint64Bytes(math.MaxInt32)
	if !reflect.DeepEqual(b, []byte("2147483647")) {
		t.Error("b must be []byte(\"2147483647\")")
	}
	b = GetUint64Bytes(math.MaxUint32)
	if !reflect.DeepEqual(b, []byte("4294967295")) {
		t.Error("b must be []byte(\"4294967295\")")
	}
	b = GetUint64Bytes(math.MaxInt64)
	if !reflect.DeepEqual(b, []byte(strconv.Itoa(math.MaxInt64))) {
		t.Error("b must be []byte(\"" + strconv.Itoa(math.MaxInt64) + "\")")
	}
	b = GetUint64Bytes(math.MaxUint64)
	if !reflect.DeepEqual(b, []byte(strconv.FormatUint(math.MaxUint64, 10))) {
		t.Error("b must be []byte(\"" + strconv.FormatUint(math.MaxUint64, 10) + "\")")
	}
}
