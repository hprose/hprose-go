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
 * rpc/tcp_service.go                                     *
 *                                                        *
 * hprose tcp service for Go.                             *
 *                                                        *
 * LastModified: Sep 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"net"
	"time"
	"unsafe"
)

// TCPService is the hprose tcp service
type TCPService struct {
	*StreamService
	Linger          int
	NoDelay         bool
	KeepAlive       bool
	KeepAlivePeriod time.Duration
	TLSConfig       *tls.Config
}

// NewTCPService is the constructor of TCPService
func NewTCPService() (service *TCPService) {
	service = new(TCPService)
	service.StreamService = (*StreamService)((unsafe.Pointer)(NewBaseService()))
	service.Linger = -1
	service.NoDelay = true
	service.KeepAlive = true
	service.KeepAlivePeriod = 0
	service.TLSConfig = nil
	return service
}

// ServeTCPConn runs on a single tcp connection. ServeTCPConn blocks, serving
// the connection until the client hangs up. The caller typically invokes
// ServeTCPConn in a go statement.
func (service *TCPService) ServeTCPConn(conn *net.TCPConn) {
	conn.SetLinger(service.Linger)
	conn.SetNoDelay(service.NoDelay)
	conn.SetKeepAlive(service.KeepAlive)
	if service.KeepAlivePeriod > 0 {
		conn.SetKeepAlivePeriod(service.KeepAlivePeriod)
	}
	if service.TLSConfig != nil {
		tlsConn := tls.Server(conn, service.TLSConfig)
		tlsConn.Handshake()
		service.ServeConn(tlsConn)
	} else {
		service.ServeConn(conn)
	}
}

// ServeTCP runs on the TCPListener. ServeTCP blocks, serving the listener
// until the server is stop. The caller typically invokes ServeTCP in a go
// statement.
func (service *TCPService) ServeTCP(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err == nil {
			break
		}
		go service.ServeTCPConn(conn)
	}
}
