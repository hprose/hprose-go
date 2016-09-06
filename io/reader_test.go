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
 * io/reader_test.go                                      *
 *                                                        *
 * hprose Reader Test for Go.                             *
 *                                                        *
 * LastModified: Sep 6, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"math"
	"testing"
)

func TestReadBool(t *testing.T) {
	trueValue := "true"
	data := map[interface{}]bool{
		true:            true,
		false:           false,
		nil:             false,
		"":              false,
		0:               false,
		1:               true,
		9:               true,
		100:             true,
		100000000000000: true,
		0.0:             false,
		"t":             true,
		"f":             false,
		&trueValue:      true,
		&trueValue:      true,
		"false":         false,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&trueValue)
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		b := reader.ReadBool()
		if b != data[k] {
			t.Error(k, data[k], b)
		}
	}
	b := reader.ReadBool()
	if b != true {
		t.Error(trueValue, true, b)
	}
	w.Close()
}

func TestUnserializeBool(t *testing.T) {
	trueValue := "true"
	data := map[interface{}]bool{
		true:            true,
		false:           false,
		nil:             false,
		"":              false,
		0:               false,
		1:               true,
		9:               true,
		100:             true,
		100000000000000: true,
		0.0:             false,
		"t":             true,
		"f":             false,
		&trueValue:      true,
		"false":         false,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&trueValue)
	reader := NewReader(w.Bytes(), false)
	var p bool
	for _, k := range keys {
		reader.Unserialize(&p)
		if p != data[k] {
			t.Error(k, data[k], p)
		}
	}
	reader.Unserialize(&p)
	if p != true {
		t.Error(trueValue, true, p)
	}
	w.Close()
}

func BenchmarkReadBool(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(true)
	bytes := w.Bytes()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.ReadBool()
	}
	w.Close()
}

func BenchmarkUnserializeBool(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(true)
	bytes := w.Bytes()
	var p bool
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.Unserialize(&p)
	}
	w.Close()
}

func TestReadInt(t *testing.T) {
	intValue := "1234567"
	u := uint(math.MaxUint64)
	data := map[interface{}]int{
		true:          1,
		false:         0,
		nil:           0,
		"":            0,
		0:             0,
		1:             1,
		9:             9,
		100:           100,
		math.MaxInt64: int(math.MaxInt64),
		u:             int(u),
		0.0:           0,
		"1":           1,
		"9":           9,
		&intValue:     1234567,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&intValue)
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		i := reader.ReadInt()
		if i != data[k] {
			t.Error(k, data[k], i)
		}
	}
	i := reader.ReadInt()
	if i != 1234567 {
		t.Error(intValue, 1234567, i)
	}
	w.Close()
}

func TestUnserializeInt(t *testing.T) {
	intValue := "1234567"
	u := uint(math.MaxUint64)
	data := map[interface{}]int{
		true:          1,
		false:         0,
		nil:           0,
		"":            0,
		0:             0,
		1:             1,
		9:             9,
		100:           100,
		math.MaxInt64: int(math.MaxInt64),
		u:             int(u),
		0.0:           0,
		"1":           1,
		"9":           9,
		&intValue:     1234567,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&intValue)
	reader := NewReader(w.Bytes(), false)
	var p int
	for _, k := range keys {
		reader.Unserialize(&p)
		if p != data[k] {
			t.Error(k, data[k], p)
		}
	}
	reader.Unserialize(&p)
	if p != 1234567 {
		t.Error(intValue, 1234567, p)
	}
	w.Close()
}

func BenchmarkReadInt(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(12345)
	bytes := w.Bytes()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.ReadInt()
	}
	w.Close()
}

func BenchmarkUnserializeInt(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(12345)
	bytes := w.Bytes()
	var p int
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.Unserialize(&p)
	}
	w.Close()
}