package main

import (
	"fmt"
	"time"

	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/util"
)

func mySliceEncoder(w *io.Writer, v interface{}) {
	slice := v.([]map[string]interface{})
	var buf [20]byte
	for i := range slice {
		w.WriteByte(io.TagMap)
		w.Write(util.GetIntBytes(buf[:], int64(len(slice[i]))))
		w.WriteByte(io.TagOpenbrace)
		for key, val := range slice[i] {
			w.WriteString(key)
			w.Serialize(val)
		}
		w.WriteByte(io.TagClosebrace)
	}
}

func main() {
	slice := make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		slice[i] = make(map[string]interface{})
		slice[i]["id"] = i
	}

	fmt.Printf("%s\r\n", io.Marshal(slice))
	start := time.Now()
	for i := 0; i < 500000; i++ {
		io.Marshal(slice)
	}
	stop := time.Now()
	fmt.Println((stop.UnixNano() - start.UnixNano()) / 1000000)

	io.RegisterSliceEncoder(([]map[string]interface{})(nil), mySliceEncoder)

	fmt.Printf("%s\r\n", io.Marshal(slice))
	start = time.Now()
	for i := 0; i < 500000; i++ {
		io.Marshal(slice)
	}
	stop = time.Now()
	fmt.Println((stop.UnixNano() - start.UnixNano()) / 1000000)
}
