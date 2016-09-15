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

// NewSocketContext is the constructor for SocketContext
func NewSocketContext(clients Clients, conn net.Conn) (context *SocketContext) {
	context = new(SocketContext)
	context.ServiceContext = NewServiceContext(clients)
	context.TransportContext = context
	context.Conn = conn
	return
}

// SocketService is the hprose socket service
type SocketService struct {
	*BaseService
}

// NewSocketService is the constructor of SocketService
func NewSocketService() (service *SocketService) {
	service = new(SocketService)
	service.BaseService = NewBaseService()
	service.fixer = socketFixer{}
	return service
}

// ServeConn runs on a single net connection. ServeConn blocks, serving the
// connection until the client hangs up. The caller typically invokes ServeConn
// in a go statement.
func (service *SocketService) ServeConn(conn net.Conn) {
	serveConn(service.BaseService, conn)
}

// Serve runs on the Listener. Serve blocks, serving the listener
// until the server is stop. The caller typically invokes Serve in a go
// statement.
func (service *SocketService) Serve(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go service.ServeConn(conn)
	}
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
	sendQueue chan packet
	conn      net.Conn
}

func serveConn(service *BaseService, conn net.Conn) {
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
	handler.sendQueue = make(chan packet, 10)
	handler.conn = conn
	go handler.init()
	handler.serve(service)
	if err := fireCloseEvent(event, context); err != nil {
		fireErrorEvent(event, err, context)
	}
}

func bytesToInt(b [4]byte) int {
	return int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
}

func (handler *connHandler) init() {
	for data := range handler.sendQueue {
		sendData(handler.conn, data)
	}
}

func (handler *connHandler) serve(service *BaseService) {
	var header [4]byte
	var size int
	var data packet
	var err error
	reader := bufio.NewReader(handler.conn)
	for {
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
			go handler.handle(service, data)
		} else {
			handler.handle(service, data)
		}
	}
	close(handler.sendQueue)
	handler.conn.Close()
}

func (handler *connHandler) handle(service *BaseService, data packet) {
	context := NewSocketContext(nil, handler.conn)
	data.body = service.Handle(data.body, context.ServiceContext)
	if data.fullDuplex {
		handler.sendQueue <- data
	} else {
		sendData(handler.conn, data)
	}
}
