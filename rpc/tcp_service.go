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

// TCPService is the hprose tcp service
type TCPService struct {
	BaseService
	Linger          int
	NoDelay         bool
	KeepAlive       bool
	KeepAlivePeriod time.Duration
	TLSConfig       *tls.Config
}

// NewTCPService is the constructor of TCPService
func NewTCPService() (service *TCPService) {
	service = new(TCPService)
	initBaseService(&service.BaseService)
	service.FixArguments = socketFixArguments
	service.Linger = -1
	service.NoDelay = true
	service.KeepAlive = true
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
		serveConn(&service.BaseService, tlsConn)
	} else {
		serveConn(&service.BaseService, conn)
	}
}

// ServeTCP runs on the TCPListener. ServeTCP blocks, serving the listener
// until the server is stop. The caller typically invokes ServeTCP in a go
// statement.
func (service *TCPService) ServeTCP(listener *net.TCPListener) {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			tempDelay = nextTempDelay(err, service.Event, tempDelay)
			if tempDelay > 0 {
				continue
			}
			return
		}
		tempDelay = 0
		go service.ServeTCPConn(conn)
	}
}
