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
 * rpc/unix_client.go                                     *
 *                                                        *
 * hprose unx client for Go.                              *
 *                                                        *
 * LastModified: Oct 5, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"net"
	"net/url"
)

// UnixClient is hprose unix client
type UnixClient struct {
	SocketClient
}

// NewUnixClient is the constructor of UnixClient
func NewUnixClient(uri ...string) (client *UnixClient) {
	client = new(UnixClient)
	client.initSocketClient()
	client.createConn = client.createUnixConn
	client.SetURIList(uri)
	return
}

func checkUnixAddresses(client Client, uriList []string) {
	for _, uri := range uriList {
		if u, err := url.Parse(uri); err == nil {
			if u.Scheme != "unix" {
				panic("This client desn't support " + u.Scheme + " scheme.")
			}
		}
	}
}

// SetURIList set a list of server addresses
func (client *UnixClient) SetURIList(uriList []string) {
	checkUnixAddresses(client, uriList)
	client.BaseClient.SetURIList(uriList)
}

func (client *UnixClient) createUnixConn() net.Conn {
	u, err := url.Parse(client.uri)
	ifErrorPanic(err)
	unixaddr, err := net.ResolveUnixAddr(u.Scheme, u.Path)
	ifErrorPanic(err)
	conn, err := net.DialUnix(u.Scheme, nil, unixaddr)
	ifErrorPanic(err)
	if client.ReadBuffer > 0 {
		ifErrorPanic(conn.SetReadBuffer(client.ReadBuffer))
	}
	if client.WriteBuffer > 0 {
		ifErrorPanic(conn.SetWriteBuffer(client.WriteBuffer))
	}
	if client.TLSConfig != nil {
		return tls.Client(conn, client.TLSConfig)
	}
	return conn
}
