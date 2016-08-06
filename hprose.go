package hprose

import (
	"github.com/hprose/hprose-golang/io"
)

// Reader is a fine-grained operation struct for Hprose unserialization
// when JSONCompatible is true, the Map data will unserialize to map[string]interface as the default type
type Reader io.Reader
