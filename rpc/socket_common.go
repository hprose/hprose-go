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
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"net"
	"strings"
)

type packet struct {
	fullDuplex bool
	id         [4]byte
	body       []byte
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
	if data.fullDuplex {
		n |= 0x80000000
	}
	buf := make([]byte, l)
	buf[0] = byte((n >> 24) & 0xff)
	buf[1] = byte((n >> 16) & 0xff)
	buf[2] = byte((n >> 8) & 0xff)
	buf[3] = byte(n & 0xff)
	i := 4
	if data.fullDuplex {
		buf[4] = data.id[0]
		buf[5] = data.id[1]
		buf[6] = data.id[2]
		buf[7] = data.id[3]
		i = 8
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

func parseUnixURI(uri string) (scheme, path string) {
	t := strings.SplitN(uri, ":", 2)
	return t[0], t[1]
}
