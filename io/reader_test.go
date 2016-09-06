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

import "testing"

func TestReadBool(t *testing.T) {
	w := NewWriter(false)
	w.Serialize(true)
	w.Serialize(false)
	w.Serialize(nil)
	w.Serialize("")
	w.Serialize(0)
	w.Serialize(1)
	w.Serialize(9)
	w.Serialize(100)
	w.Serialize(10000000000000)
	w.Serialize(0.0)
	w.Serialize("t")
	w.Serialize("f")
	trueValue := "true"
	w.Serialize(&trueValue)
	w.Serialize(&trueValue)
	w.Serialize("false")
	reader := NewReader(w.Bytes(), false)
	b := reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != false {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != false {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != false {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != false {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != false {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != false {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != true {
		t.Error(b)
	}
	b = reader.ReadBool()
	if b != false {
		t.Error(b)
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
