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
 * rpc/unix_service.go                                    *
 *                                                        *
 * hprose unix service for Go.                            *
 *                                                        *
 * LastModified: Sep 30, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"net"
	"time"
)

// UnixService is the hprose unix service
type UnixService struct {
	BaseService
	TLSConfig *tls.Config
}

// NewUnixService is the constructor of UnixService
func NewUnixService() (service *UnixService) {
	service = new(UnixService)
	initBaseService(&service.BaseService)
	service.FixArguments = socketFixArguments
	return service
}

// ServeUnixConn runs on a single tcp connection. ServeUnixConn blocks, serving
// the connection until the client hangs up. The caller typically invokes
// ServeUnixConn in a go statement.
func (service *UnixService) ServeUnixConn(conn *net.UnixConn) {
	if service.TLSConfig != nil {
		tlsConn := tls.Server(conn, service.TLSConfig)
		tlsConn.Handshake()
		serveConn(&service.BaseService, tlsConn)
	} else {
		serveConn(&service.BaseService, conn)
	}
}

// ServeUnix runs on the UnixListener. ServeUnix blocks, serving the listener
// until the server is stop. The caller typically invokes ServeUnix in a go
// statement.
func (service *UnixService) ServeUnix(listener *net.UnixListener) {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := listener.AcceptUnix()
		if err != nil {
			tempDelay = nextTempDelay(err, service.Event, tempDelay)
			if tempDelay > 0 {
				continue
			}
			return
		}
		tempDelay = 0
		go service.ServeUnixConn(conn)
	}
}
