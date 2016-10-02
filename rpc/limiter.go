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
 * rpc/limiter.go                                         *
 *                                                        *
 * hprose client requests limiter for Go.                 *
 *                                                        *
 * LastModified: Oct 2, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import "sync"

type limiter struct {
	cond                  sync.Cond
	reqcount              int
	MaxConcurrentRequests int
}

func (limiter *limiter) initLimiter() {
	limiter.MaxConcurrentRequests = 10
	limiter.cond.L = &sync.Mutex{}
}

func (limiter *limiter) limit() {
	for {
		if limiter.reqcount < limiter.MaxConcurrentRequests {
			break
		}
		limiter.cond.Wait()
	}
	limiter.reqcount++
}

func (limiter *limiter) unlimit() {
	limiter.reqcount--
	limiter.cond.Signal()
}

func (limiter *limiter) reset() {
	limiter.reqcount = 0
	for i := 0; i < limiter.MaxConcurrentRequests; i++ {
		limiter.cond.Signal()
	}
}
