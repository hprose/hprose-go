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
 * rpc/tcp_server.go                                      *
 *                                                        *
 * hprose tcp server for Go.                              *
 *                                                        *
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

// TCPServer is a hprose tcp server
type TCPServer struct {
	*TCPService
	uri      string
	listener *net.TCPListener
	signal   chan os.Signal
}

// NewTCPServer is the constructor for TCPServer
func NewTCPServer(uri string) (server *TCPServer) {
	if uri == "" {
		uri = "tcp://127.0.0.1:0"
	}
	server = new(TCPServer)
	server.TCPService = NewTCPService()
	server.uri = uri
	return
}

// URI return the real address of this server
func (server *TCPServer) URI() string {
	if server.listener == nil {
		panic(errServerIsNotStarted)
	}
	u, err := url.Parse(server.uri)
	if err != nil {
		panic(err)
	}
	return u.Scheme + "://" + server.listener.Addr().String()
}

// Handle the hprose tcp server
func (server *TCPServer) Handle() (err error) {
	if server.listener != nil {
		return errServerIsAlreadyStarted
	}
	u, err := url.Parse(server.uri)
	if err != nil {
		return err
	}
	addr, err := net.ResolveTCPAddr(u.Scheme, u.Host)
	if err != nil {
		return err
	}
	if server.listener, err = net.ListenTCP(u.Scheme, addr); err != nil {
		return err
	}
	go server.ServeTCP(server.listener)
	return nil
}

// Start the hprose tcp server
func (server *TCPServer) Start() (err error) {
	for {
		if err = server.Handle(); err != nil {
			return err
		}
		server.signal = make(chan os.Signal, 1)
		signal.Notify(server.signal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGKILL)
		s := <-server.signal
		server.Stop()
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGKILL:
			return nil
		}
	}
}

// Stop the hprose tcp server
func (server *TCPServer) Stop() {
	if server.signal != nil {
		signal.Stop(server.signal)
		server.signal = nil
	}
	if server.listener != nil {
		listener := server.listener
		server.listener = nil
		listener.Close()
	}
}
