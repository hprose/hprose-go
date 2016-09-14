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
 * rpc/socket_service.go                                  *
 *                                                        *
 * hprose socket service for Go.                          *
 *                                                        *
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"bufio"
	"net"
	"reflect"
)

// SocketContext is the hprose socket context for service
type SocketContext struct {
	*ServiceContext
	net.Conn
}

type acceptEvent interface {
	OnAccept(context *SocketContext)
}

type acceptEvent2 interface {
	OnAccept(context *SocketContext) error
}

type closeEvent interface {
	OnClose(context *SocketContext)
}

type closeEvent2 interface {
	OnClose(context *SocketContext) error
}

type serviceHandler struct {
	sendQueue chan packet
	conn      net.Conn
	service   *BaseService
}

type socketFixer struct{}

func (socketFixer) FixArguments(args []reflect.Value, context *ServiceContext) {
	i := len(args) - 1
	typ := args[i].Type()
	if typ == socketContextType {
		if c, ok := context.TransportContext.(*SocketContext); ok {
			args[i] = reflect.ValueOf(c)
		}
		return
	}
	if typ == netConnType {
		if c, ok := context.TransportContext.(*SocketContext); ok {
			args[i] = reflect.ValueOf(c.Conn)
		}
		return
	}
	fixArguments(args, context)
}

func bytesToInt(b [4]byte) int {
	return int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
}

func fireAcceptEvent(service *BaseService, context *SocketContext) (err error) {
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

func fireCloseEvent(service *BaseService, context *SocketContext) (err error) {
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
	context := &SocketContext{NewServiceContext(nil), conn}
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
		sendData(handler.conn, data)
	}
}

func (handler *serviceHandler) serve() {
	var header [4]byte
	var size int
	var data packet
	var err error
	reader := bufio.NewReader(handler.conn)
	for {
		context := &SocketContext{NewServiceContext(nil), handler.conn}
		context.TransportContext = context
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
			go handler.handle(data, context.ServiceContext)
		} else {
			handler.handle(data, context.ServiceContext)
		}
	}
	close(handler.sendQueue)
	handler.conn.Close()
}

func (handler *serviceHandler) handle(data packet, context *ServiceContext) {
	if resp, err := handler.service.Handle(data.body, context); err == nil {
		data.body = resp
	} else {
		data.body = handler.service.endError(err, context)
	}
	if data.fullDuplex {
		handler.sendQueue <- data
	} else {
		sendData(handler.conn, data)
	}
}
