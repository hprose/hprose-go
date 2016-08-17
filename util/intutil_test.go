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
	data := []int32{
		0, 9, 10, 99, 100, 999, 1000, -1000, 10000, -10000,
		123456789, -123456789, math.MaxInt32, math.MinInt32}
	for _, i := range data {
		b := GetInt32Bytes(i)
		if !reflect.DeepEqual(b, []byte(strconv.Itoa(int(i)))) {
			t.Error("b must be []byte(\"" + strconv.Itoa(int(i)) + "\")")
		}
	}
}

func TestGetInt64Bytes(t *testing.T) {
	data := []int64{
		0, 9, 10, 99, 100, 999, 1000, -1000, 10000, -10000,
		123456789, -123456789, math.MaxInt32, math.MinInt32,
		math.MaxInt64, math.MinInt64}
	for _, i := range data {
		b := GetInt64Bytes(i)
		if !reflect.DeepEqual(b, []byte(strconv.Itoa(int(i)))) {
			t.Error("b must be []byte(\"" + strconv.Itoa(int(i)) + "\")")
		}
	}
}

func TestGetUint32Bytes(t *testing.T) {
	data := []uint32{
		0, 9, 10, 99, 100, 999, 1000, 10000, 123456789,
		math.MaxInt32, math.MaxUint32}
	for _, i := range data {
		b := GetUint32Bytes(i)
		if !reflect.DeepEqual(b, []byte(strconv.Itoa(int(i)))) {
			t.Error("b must be []byte(\"" + strconv.Itoa(int(i)) + "\")")
		}
	}
}

func TestGetUint64Bytes(t *testing.T) {
	data := []uint64{
		0, 9, 10, 99, 100, 999, 1000, 10000, 123456789,
		math.MaxInt32, math.MaxUint32, math.MaxInt64, math.MaxUint64}
	for _, i := range data {
		b := GetUint64Bytes(i)
		if !reflect.DeepEqual(b, []byte(strconv.FormatUint(i, 10))) {
			t.Error("b must be []byte(\"" + strconv.FormatUint(i, 10) + "\")")
		}
	}
}
