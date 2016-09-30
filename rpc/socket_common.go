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
 * rpc/socket_common.go                                   *
 *                                                        *
 * hprose socket common for Go.                           *
 *                                                        *
 * LastModified: Sep 25, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"net"
	"time"
)

type packet struct {
	fullDuplex bool
	id         [4]byte
	body       []byte
}

func toUint32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func fromUint32(b []byte, i uint32) {
	b[0] = byte(i >> 24)
	b[1] = byte(i >> 16)
	b[2] = byte(i >> 8)
	b[3] = byte(i)
}

func sendData(conn net.Conn, data packet) (err error) {
	n := len(data.body)
	var l int
	switch {
	case n > 1020 && n <= 1400:
		l = 2048
	case n > 508:
		l = 1024
	default:
		l = 512
	}
	buf := make([]byte, l)
	i := 4
	if data.fullDuplex {
		fromUint32(buf, uint32(n|0x80000000))
		buf[4] = data.id[0]
		buf[5] = data.id[1]
		buf[6] = data.id[2]
		buf[7] = data.id[3]
		i = 8
	} else {
		fromUint32(buf, uint32(n))
	}
	p := l - i
	if n <= p {
		copy(buf[i:], data.body)
		_, err = conn.Write(buf[:n+i])
	} else {
		copy(buf[i:], data.body[:p])
		_, err = conn.Write(buf)
		if err != nil {
			return err
		}
		_, err = conn.Write(data.body[p:])
	}
	return err
}

func nextTempDelay(
	err error, event ServiceEvent, tempDelay time.Duration) time.Duration {
	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}
		fireErrorEvent(event, err, nil)
		time.Sleep(tempDelay)
		return tempDelay
	}
	return 0
}
