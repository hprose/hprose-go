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
 * LastModified: Oct 2, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"net"
	"os"
	"strings"
)

// UnixServer is a hprose unix server
type UnixServer struct {
	UnixService
	starter
	uri      string
	listener *net.UnixListener
}

// NewUnixServer is the constructor for UnixServer
func NewUnixServer(uri string) (server *UnixServer) {
	if uri == "" {
		uri = "unix:/tmp/hprose.sock"
	}
	server = new(UnixServer)
	server.initUnixService()
	server.starter.server = server
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

// Close the hprose unix server
func (server *UnixServer) Close() {
	if server.listener != nil {
		listener := server.listener
		server.listener = nil
		listener.Close()
	}
}

func (server *UnixServer) signal() chan os.Signal {
	return server.c
}

func parseUnixURI(uri string) (scheme, path string) {
	t := strings.SplitN(uri, ":", 2)
	return t[0], t[1]
}
