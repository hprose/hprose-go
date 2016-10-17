package test_interface

import "testing"

func BenchmarkStruct(b *testing.B) {
	m := make(map[int]struct{})
	for i := 0; i < b.N; i++ {
		m[i] = struct{}{}
	}
}

func BenchmarkBool(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < b.N; i++ {
		m[i] = true
	}
}
