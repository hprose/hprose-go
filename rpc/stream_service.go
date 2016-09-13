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
 * rpc/stream_service.go                                  *
 *                                                        *
 * hprose stream service for Go.                          *
 *                                                        *
 * LastModified: Sep 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"io"
	"net"
)

// StreamContext is the hprose stream context for service
type StreamContext struct {
	*ServiceContext
	net.Conn
}

type packet struct {
	fullDuplex bool
	id         [4]byte
	body       []byte
	context    *ServiceContext
}

// StreamService is the base service for TcpService and UnixService
type StreamService struct {
	BaseService
}

func (service *StreamService) initSendQueue(
	sendQueue chan packet, conn net.Conn) {
	var header [4]byte
	var size int
	var data packet
	var err error
	var ok bool
	for {
		data, ok = <-sendQueue
		if !ok {
			break
		}
		size = len(data.body)
		if data.fullDuplex {
			size |= 0x80000000
		}
		header[0] = byte((size >> 24) & 0xFF)
		header[1] = byte((size >> 16) & 0xFF)
		header[2] = byte((size >> 28) & 0xFF)
		header[3] = byte(size & 0xFF)
		if _, err = conn.Write(header[:]); err != nil {
			break
		}
		if data.fullDuplex {
			if _, err = conn.Write(data.id[:]); err != nil {
				break
			}
		}
		if _, err = conn.Write(data.body); err != nil {
			break
		}
	}
	service.fireErrorEvent(err, data.context)
	conn.Close()
}

func (service *StreamService) onReceived(data packet, sendQueue chan packet) {
	resp, err := service.Handle(data.body, data.context).Get()
	if err == nil {
		data.body = resp.([]byte)
	} else {
		data.body = service.endError(err, data.context)
	}
	sendQueue <- data
}

func bytesToInt(b [4]byte) int {
	return int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
}

// ServeConn runs on a single connection. ServeConn blocks, serving the
// connection until the client hangs up. The caller typically invokes ServeConn // in a go statement.
func (service *StreamService) ServeConn(conn net.Conn) {
	var header [4]byte
	var size int
	var data packet
	var err error
	sendQueue := make(chan packet, 10)
	go service.initSendQueue(sendQueue, conn)
	for {
		context := &StreamContext{NewServiceContext(nil), conn}
		context.TransportContext = context
		data.context = context.ServiceContext
		if _, err := io.ReadAtLeast(conn, header[:], 4); err != nil {
			break
		}
		size = bytesToInt(header)
		if size&0x8000000 != 0 {
			size &= 0x7FFFFFF
			data.fullDuplex = true
			if _, err = io.ReadAtLeast(conn, data.id[:], 4); err != nil {
				break
			}
		} else {
			data.fullDuplex = false
		}
		data.body = make([]byte, size)
		if _, err = io.ReadAtLeast(conn, data.body, size); err != nil {
			break
		}
		go service.onReceived(data, sendQueue)
	}
	service.fireErrorEvent(err, data.context)
	close(sendQueue)
	conn.Close()
}
