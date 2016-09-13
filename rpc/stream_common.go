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
 * rpc/stream_common.go                                   *
 *                                                        *
 * hprose stream common for Go.                           *
 *                                                        *
 * LastModified: Sep 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import "net"

func sendDataOverStream(
	conn net.Conn,
	data []byte,
	id [4]byte,
	fullDuplex bool) (err error) {
	n := len(data)
	var l int
	switch {
	case n > 1020 && n <= 1400:
		l = 2048
	case n > 508:
		l = 1024
	default:
		l = 512
	}
	if fullDuplex {
		n |= 0x80000000
	}
	buf := make([]byte, l)
	buf[0] = byte((n >> 24) & 0xff)
	buf[1] = byte((n >> 16) & 0xff)
	buf[2] = byte((n >> 8) & 0xff)
	buf[3] = byte(n & 0xff)
	i := 4
	if fullDuplex {
		buf[4] = id[0]
		buf[5] = id[1]
		buf[6] = id[2]
		buf[7] = id[3]
		i = 8
	}
	p := l - i
	if n <= p {
		copy(buf[i:], data)
		_, err = conn.Write(buf[:n+i])
	} else {
		copy(buf[i:], data[:p])
		_, err = conn.Write(buf)
		if err != nil {
			return err
		}
		_, err = conn.Write(data[p:])
	}
	return err
}
