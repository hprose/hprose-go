package io

// Reader is a fine-grained operation struct for Hprose unserialization
// when JSONCompatible is true, the Map data will unserialize to map[string]interface as the default type
type Reader struct {
	classref       []interface{}
	fieldsref      [][]string
	JSONCompatible bool
}
