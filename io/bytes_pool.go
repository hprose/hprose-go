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
 * io/bytes_pool.go                                       *
 *                                                        *
 * bytes pool for Go.                                     *
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

type bytesPool struct {
	pools [poolNum]pool
	timer *time.Timer
	d     time.Duration
}

func newBytesPool(d time.Duration) (bp *bytesPool) {
	bp = new(bytesPool)
	bp.d = d
	if d > 0 {
		bp.timer = time.AfterFunc(d, func() {
			bp.Drain()
			bp.timer.Reset(d)
		})
	}
	return bp
}

// BytesPool is a pool of []byte.
var BytesPool = newBytesPool(time.Second * 2)

// Get a []byte from pool.
func (bp *bytesPool) Get(size int) []byte {
	if size < 1 || size > maxSize {
		return make([]byte, size)
	}
	if bp.d > 0 {
		bp.timer.Reset(bp.d)
	}
	var bytes []byte
	capacity := pow2roundup(size)
	if capacity < 64 {
		capacity = 64
	}
	p := &bp.pools[log2(capacity)-6]
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

// Put a []byte to pool.
func (bp *bytesPool) Put(bytes []byte) {
	capacity := cap(bytes)
	if capacity < 64 || capacity > maxSize || capacity != pow2roundup(capacity) {
		return
	}
	p := &bp.pools[log2(capacity)-6]
	p.locker.Lock()
	p.list = append(p.list, bytes[:capacity])
	p.locker.Unlock()

}

// Drain some items from the pool and make them availabe for garbage collection.
func (bp *bytesPool) Drain() {
	n := len(bp.pools)
	for i := 0; i < n; i++ {
		p := &bp.pools[i]
		p.locker.Lock()
		p.list = p.list[:len(p.list)>>1]
		p.locker.Unlock()
	}
}
