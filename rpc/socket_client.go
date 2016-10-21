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
 * rpc/socket_client.go                                   *
 *                                                        *
 * hprose socket client for Go.                           *
 *                                                        *
 * LastModified: Oct 21, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"net"
)

type socketTransport interface {
	MaxPoolSize() int
	SetMaxPoolSize(size int)
	sendAndReceive(data []byte, context *ClientContext) ([]byte, error)
	setCreateConn(createConn func() net.Conn)
}

// SocketClient is base struct for TCPClient and UnixClient
type SocketClient struct {
	baseClient
	socketTransport
	ReadBuffer  int
	WriteBuffer int
	TLSConfig   *tls.Config
}

func (client *SocketClient) initSocketClient() {
	client.initBaseClient()
	client.socketTransport = newHalfDuplexSocketTransport()
	client.ReadBuffer = 0
	client.WriteBuffer = 0
	client.TLSConfig = nil
	client.SendAndReceive = client.sendAndReceive
}

// TLSClientConfig returns the tls.Config in hprose client
func (client *SocketClient) TLSClientConfig() *tls.Config {
	return client.TLSConfig
}

// SetTLSClientConfig sets the tls.Config
func (client *SocketClient) SetTLSClientConfig(config *tls.Config) {
	client.TLSConfig = config
}
