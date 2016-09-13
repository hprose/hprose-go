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
	"bufio"
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
	sendQueue <-chan packet, conn net.Conn) {
	for data := range sendQueue {
		sendDataOverStream(conn, data.body, data.id, data.fullDuplex)
	}
	conn.Close()
}

func (service *StreamService) onReceived(
	conn net.Conn, data packet, sendQueue chan<- packet) {
	if resp, err := service.Handle(data.body, data.context); err == nil {
		data.body = resp
	} else {
		data.body = service.endError(err, data.context)
	}
	if data.fullDuplex {
		sendQueue <- data
	} else {
		sendDataOverStream(conn, data.body, data.id, data.fullDuplex)
	}
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
	reader := bufio.NewReader(conn)
	sendQueue := make(chan packet, 10)
	go service.initSendQueue(sendQueue, conn)
	for {
		context := &StreamContext{NewServiceContext(nil), conn}
		context.TransportContext = context
		data.context = context.ServiceContext
		if _, err = reader.Read(header[:]); err != nil {
			break
		}
		size = bytesToInt(header)
		data.fullDuplex = (size&0x8000000 != 0)
		if data.fullDuplex {
			size &= 0x7FFFFFF
			data.fullDuplex = true
			if _, err = reader.Read(data.id[:]); err != nil {
				break
			}
		}
		data.body = make([]byte, size)
		if _, err = reader.Read(data.body); err != nil {
			break
		}
		if data.fullDuplex {
			go service.onReceived(conn, data, sendQueue)
		} else {
			service.onReceived(conn, data, sendQueue)
		}
	}
	close(sendQueue)
	conn.Close()
}
