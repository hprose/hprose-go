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
 * rpc/topic.go                                           *
 *                                                        *
 * hprose push topic for Go.                              *
 *                                                        *
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"time"

	"github.com/hprose/hprose-golang/promise"
)

type message struct {
	Detector chan bool
	Result   interface{}
}

type topic struct {
	*time.Timer
	Request   promise.Promise
	Messages  chan *message
	Count     int64
	Heartbeat time.Duration
}

func newTopic(heartbeat time.Duration) *topic {
	t := new(topic)
	t.Heartbeat = heartbeat
	return t
}
