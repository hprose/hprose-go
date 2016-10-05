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
 * rpc/tcp_client.go                                      *
 *                                                        *
 * hprose tcp client for Go.                              *
 *                                                        *
 * LastModified: Oct 3, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"net"
	"net/url"
	"time"
)

// TCPClient is hprose tcp client
type TCPClient struct {
	SocketClient
	Linger          int
	NoDelay         bool
	KeepAlive       bool
	KeepAlivePeriod time.Duration
}

// NewTCPClient is the constructor of TCPClient
func NewTCPClient(uri ...string) (client *TCPClient) {
	client = new(TCPClient)
	client.initSocketClient()
	client.Linger = -1
	client.NoDelay = true
	client.KeepAlive = true
	client.createConn = client.createTCPConn
	client.SetURIList(uri)
	return
}

func checkTCPAddresses(client Client, uriList []string) {
	for _, uri := range uriList {
		if u, err := url.Parse(uri); err == nil {
			if u.Scheme != "tcp" && u.Scheme != "tcp4" && u.Scheme != "tcp6" {
				panic("This client desn't support " + u.Scheme + " scheme.")
			}
		}
	}
}

// SetURIList set a list of server addresses
func (client *TCPClient) SetURIList(uriList []string) {
	checkTCPAddresses(client, uriList)
	client.BaseClient.SetURIList(uriList)
}

func (client *TCPClient) createTCPConn() net.Conn {
	u, err := url.Parse(client.uri)
	ifErrorPanic(err)
	tcpaddr, err := net.ResolveTCPAddr(u.Scheme, u.Host)
	ifErrorPanic(err)
	conn, err := net.DialTCP(u.Scheme, nil, tcpaddr)
	ifErrorPanic(err)
	ifErrorPanic(conn.SetLinger(client.Linger))
	ifErrorPanic(conn.SetNoDelay(client.NoDelay))
	ifErrorPanic(conn.SetKeepAlive(client.KeepAlive))
	if client.KeepAlivePeriod > 0 {
		ifErrorPanic(conn.SetKeepAlivePeriod(client.KeepAlivePeriod))
	}
	if client.ReadBuffer > 0 {
		ifErrorPanic(conn.SetReadBuffer(client.ReadBuffer))
	}
	if client.WriteBuffer > 0 {
		ifErrorPanic(conn.SetWriteBuffer(client.WriteBuffer))
	}
	if client.tlsConfig != nil {
		return tls.Client(conn, client.tlsConfig)
	}
	return conn
}
