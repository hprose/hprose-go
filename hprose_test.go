package hprose

import "testing"

func TestReader(t *testing.T) {
	reader := new(Reader)
	reader.JSONCompatible = true
}
