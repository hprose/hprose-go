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
 * LastModified: Oct 2, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"bufio"
	"net"
	"reflect"
	"sync"
)

// SocketContext is the hprose socket context for service
type SocketContext struct {
	ServiceContext
	net.Conn
}

// NewSocketContext is the constructor for SocketContext
func NewSocketContext(service Service, conn net.Conn) (context *SocketContext) {
	context = new(SocketContext)
	context.initServiceContext(service)
	context.TransportContext = context
	context.Conn = conn
	return
}

func socketFixArguments(args []reflect.Value, context *ServiceContext) {
	i := len(args) - 1
	switch args[i].Type() {
	case socketContextType:
		if c, ok := context.TransportContext.(*SocketContext); ok {
			args[i] = reflect.ValueOf(c)
		}
	case netConnType:
		if c, ok := context.TransportContext.(*SocketContext); ok {
			args[i] = reflect.ValueOf(c.Conn)
		}
	default:
		DefaultFixArguments(args, context)
	}
}

// SocketService is the hprose socket service
type SocketService struct {
	BaseService
}

func (service *SocketService) serveConn(conn net.Conn) {
	context := NewSocketContext(nil, conn)
	event := service.Event
	defer func() {
		if e := recover(); e != nil {
			err := NewPanicError(e)
			fireErrorEvent(event, err, context)
		}
	}()
	if err := fireAcceptEvent(event, context); err != nil {
		fireErrorEvent(event, err, context)
		return
	}
	handler := new(connHandler)
	handler.conn = conn
	handler.serve(service)
	if err := fireCloseEvent(event, context); err != nil {
		fireErrorEvent(event, err, context)
	}
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

func fireAcceptEvent(event ServiceEvent, context *SocketContext) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = NewPanicError(e)
		}
	}()
	switch event := event.(type) {
	case acceptEvent:
		event.OnAccept(context)
	case acceptEvent2:
		err = event.OnAccept(context)
	}
	return err
}

func fireCloseEvent(event ServiceEvent, context *SocketContext) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = NewPanicError(e)
		}
	}()
	switch event := event.(type) {
	case closeEvent:
		event.OnClose(context)
	case closeEvent2:
		err = event.OnClose(context)
	}
	return err
}

type connHandler struct {
	sync.Mutex
	conn net.Conn
}

func (handler *connHandler) serve(service *SocketService) {
	header := make([]byte, 4)
	var size uint32
	var data packet
	var err error
	reader := bufio.NewReader(handler.conn)
	for {
		if _, err = reader.Read(header); err != nil {
			break
		}
		size = toUint32(header)
		data.fullDuplex = (size&0x80000000 != 0)
		if data.fullDuplex {
			size &= 0x7FFFFFFF
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
			go handler.handle(service, data)
		} else {
			handler.handle(service, data)
		}
	}
	handler.conn.Close()
}

func (handler *connHandler) handle(service *SocketService, data packet) {
	context := NewSocketContext(service, handler.conn)
	data.body = service.Handle(data.body, &context.ServiceContext)
	if data.fullDuplex {
		handler.Lock()
	}
	err := sendData(handler.conn, data)
	if data.fullDuplex {
		handler.Unlock()
	}
	if err != nil {
		fireErrorEvent(service.Event, err, &context.ServiceContext)
	}
}
