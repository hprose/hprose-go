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
 * LastModified: Oct 5, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"io"
	"net"
	"runtime"
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

func recvData(reader io.Reader, data *packet) (err error) {
	header := data.id[:]
	if _, err = reader.Read(header); err != nil {
		return
	}
	size := toUint32(header)
	data.fullDuplex = (size&0x80000000 != 0)
	if data.fullDuplex {
		size &= 0x7FFFFFFF
		data.fullDuplex = true
		data.body = nil
		if _, err = reader.Read(data.id[:]); err != nil {
			return
		}
	}
	if cap(data.body) >= int(size) {
		data.body = data.body[:size]
	} else {
		data.body = make([]byte, size)
	}
	_, err = reader.Read(data.body)
	return
}

var bufferPool = make(chan []byte, runtime.NumCPU()*2)

func acquireBuffer() (buf []byte) {
	select {
	case buf = <-bufferPool:
		return
	default:
		return make([]byte, 2048)
	}
}

func releaseBuffer(buf []byte) {
	select {
	case bufferPool <- buf:
	default:
	}
}

func sendData(writer io.Writer, data packet) (err error) {
	n := len(data.body)
	i := 4
	buf := acquireBuffer()
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
	p := 2048 - i
	if n <= p {
		copy(buf[i:], data.body)
		_, err = writer.Write(buf[:n+i])
		releaseBuffer(buf)
	} else {
		copy(buf[i:], data.body[:p])
		_, err = writer.Write(buf)
		releaseBuffer(buf)
		if err != nil {
			return err
		}
		_, err = writer.Write(data.body[p:])
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
