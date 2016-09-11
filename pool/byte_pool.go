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
 * pool/byte_pool.go                                      *
 *                                                        *
 * byte pool for Go.                                      *
 *                                                        *
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package pool

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

func pow2roundup(x int) int {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

var debruijn = []int{
	0, 1, 28, 2, 29, 14, 24, 3,
	30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7,
	26, 12, 18, 6, 11, 5, 10, 9,
}

func log2(x int) int {
	return debruijn[uint32(x*0x077CB531)>>27]
}

func init() {
	bytePool.d = time.Second << 2
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
		l := len(p.list)
		for j := l - 1; j >= l>>1; j-- {
			p.list[j] = nil
		}
		p.list = p.list[:l>>1]
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
