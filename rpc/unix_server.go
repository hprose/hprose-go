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
 * rpc/unix_server.go                                     *
 *                                                        *
 * hprose unix server for Go.                             *
 *                                                        *
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"
)

// UnixServer is a hprose unix server
type UnixServer struct {
	*UnixService
	uri      string
	listener *net.UnixListener
	signal   chan os.Signal
}

// NewUnixServer is the constructor for UnixServer
func NewUnixServer(uri string) (server *UnixServer) {
	if uri == "" {
		uri = "unix:/tmp/hprose.sock"
	}
	server = new(UnixServer)
	server.UnixService = NewUnixService()
	server.uri = uri
	return
}

// URI return the real address of this server
func (server *UnixServer) URI() string {
	if server.listener == nil {
		panic(errServerIsNotStarted)
	}
	return "unix:" + server.listener.Addr().String()
}

// Handle the hprose unix server
func (server *UnixServer) Handle() (err error) {
	if server.listener != nil {
		return errServerIsAlreadyStarted
	}
	scheme, path := parseUnixURI(server.uri)
	if err != nil {
		return err
	}
	addr, err := net.ResolveUnixAddr(scheme, path)
	if err != nil {
		return err
	}
	if server.listener, err = net.ListenUnix(scheme, addr); err != nil {
		return err
	}
	go server.ServeUnix(server.listener)
	return nil
}

// Start the hprose unix server
func (server *UnixServer) Start() (err error) {
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
func (server *UnixServer) Stop() {
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
