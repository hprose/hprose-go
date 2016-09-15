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
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import "net"

// UnixService is the hprose unix service
type UnixService struct {
	*BaseService
}

// NewUnixService is the constructor of UnixService
func NewUnixService() (service *UnixService) {
	service = new(UnixService)
	service.BaseService = NewBaseService()
	service.fixer = socketFixer{}
	return service
}

// ServeUnixConn runs on a single tcp connection. ServeUnixConn blocks, serving
// the connection until the client hangs up. The caller typically invokes
// ServeUnixConn in a go statement.
func (service *UnixService) ServeUnixConn(conn *net.UnixConn) {
	serveConn(service.BaseService, conn)
}

// ServeUnix runs on the UnixListener. ServeUnix blocks, serving the listener
// until the server is stop. The caller typically invokes ServeUnix in a go
// statement.
func (service *UnixService) ServeUnix(listener *net.UnixListener) {
	for {
		conn, err := listener.AcceptUnix()
		if err != nil {
			break
		}
		go service.ServeUnixConn(conn)
	}
}
