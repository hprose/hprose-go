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
 * LastModified: Sep 14, 2016                             *
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

type acceptEvent interface {
	OnAccept(context *StreamContext)
}

type acceptEvent2 interface {
	OnAccept(context *StreamContext) error
}

type closeEvent interface {
	OnClose(context *StreamContext)
}

type closeEvent2 interface {
	OnClose(context *StreamContext) error
}

type packet struct {
	fullDuplex bool
	id         [4]byte
	body       []byte
	context    *ServiceContext
}

type serviceHandler struct {
	sendQueue chan packet
	conn      net.Conn
	service   *BaseService
}

func bytesToInt(b [4]byte) int {
	return int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
}

func fireAcceptEvent(service *BaseService, context *StreamContext) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = NewPanicError(e)
		}
	}()
	switch event := service.Event.(type) {
	case acceptEvent:
		event.OnAccept(context)
	case acceptEvent2:
		err = event.OnAccept(context)
	}
	return err
}

func fireCloseEvent(service *BaseService, context *StreamContext) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = NewPanicError(e)
		}
	}()
	switch event := service.Event.(type) {
	case closeEvent:
		event.OnClose(context)
	case closeEvent2:
		err = event.OnClose(context)
	}
	return err
}

func serveConn(conn net.Conn, service *BaseService) {
	context := &StreamContext{NewServiceContext(nil), conn}
	context.TransportContext = context
	if err := fireAcceptEvent(service, context); err != nil {
		service.fireErrorEvent(err, context)
		return
	}
	handler := new(serviceHandler)
	handler.sendQueue = make(chan packet, 10)
	handler.conn = conn
	handler.service = service
	go handler.init()
	handler.serve()
	if err := fireCloseEvent(service, context); err != nil {
		service.fireErrorEvent(err, context)
	}
}

func (handler *serviceHandler) init() {
	for data := range handler.sendQueue {
		sendDataOverStream(handler.conn, data.body, data.id, data.fullDuplex)
	}
}

func (handler *serviceHandler) serve() {
	var header [4]byte
	var size int
	var data packet
	var err error
	reader := bufio.NewReader(handler.conn)
	for {
		context := &StreamContext{NewServiceContext(nil), handler.conn}
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
			go handler.handle(data)
		} else {
			handler.handle(data)
		}
	}
	close(handler.sendQueue)
	handler.conn.Close()
}

func (handler *serviceHandler) handle(data packet) {
	if resp, err := handler.service.Handle(data.body, data.context); err == nil {
		data.body = resp
	} else {
		data.body = handler.service.endError(err, data.context)
	}
	if data.fullDuplex {
		handler.sendQueue <- data
	} else {
		sendDataOverStream(handler.conn, data.body, data.id, data.fullDuplex)
	}
}
