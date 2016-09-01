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
 * io/byte_pool.go                                        *
 *                                                        *
 * byte pool for Go.                                      *
 *                                                        *
 * LastModified: Sep 1, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"sync"
	"time"
)

const (
	poolNum = 20
	maxSize = 1 << (poolNum + 5)
)

type pool struct {
	list   [][]byte
	locker sync.Mutex
}

var bytePool = struct {
	pools [poolNum]pool
	timer *time.Timer
	d     time.Duration
}{}

func init() {
	bytePool.d = time.Second * 4
	if bytePool.d > 0 {
		bytePool.timer = time.AfterFunc(bytePool.d, func() {
			drain()
			bytePool.timer.Reset(bytePool.d)
		})
	}
}

func drain() {
	n := len(bytePool.pools)
	for i := 0; i < n; i++ {
		p := &bytePool.pools[i]
		p.locker.Lock()
		p.list = p.list[:len(p.list)>>1]
		p.locker.Unlock()
	}
}

// Alloc a []byte from pool.
func Alloc(size int) []byte {
	if size < 1 || size > maxSize {
		return make([]byte, size)
	}
	if bytePool.d > 0 {
		bytePool.timer.Reset(bytePool.d)
	}
	var bytes []byte
	capacity := pow2roundup(size)
	if capacity < 64 {
		capacity = 64
	}
	p := &bytePool.pools[log2(capacity)-6]
	p.locker.Lock()
	if n := len(p.list); n > 0 {
		bytes = p.list[n-1]
		p.list[n-1] = nil
		p.list = p.list[:n-1]
	}
	p.locker.Unlock()
	if bytes == nil {
		return make([]byte, size, capacity)
	}
	return bytes[:size]
}

// Recycle a []byte to pool.
func Recycle(bytes []byte) {
	capacity := cap(bytes)
	if capacity < 64 || capacity > maxSize || capacity != pow2roundup(capacity) {
		return
	}
	p := &bytePool.pools[log2(capacity)-6]
	p.locker.Lock()
	p.list = append(p.list, bytes[:capacity])
	p.locker.Unlock()
}
