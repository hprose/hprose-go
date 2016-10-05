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
 * LastModified: Oct 3, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type connEntry struct {
	conn  net.Conn
	timer *time.Timer
}

// SocketClient is base struct for TCPClient and UnixClient
type SocketClient struct {
	BaseClient
	ReadBuffer  int
	WriteBuffer int
	IdleTimeout time.Duration
	tlsConfig   *tls.Config
	pool        chan *connEntry
	connCount   int32
	nextid      uint32
	createConn  func() net.Conn
	cond        sync.Cond
}

func (client *SocketClient) initSocketClient() {
	client.initBaseClient()
	client.IdleTimeout = 30 * time.Second
	client.pool = make(chan *connEntry, runtime.NumCPU()*2)
	client.cond.L = &sync.Mutex{}
	client.SetFullDuplex(false)
}

// TLSClientConfig returns the tls.Config in hprose client
func (client *SocketClient) TLSClientConfig() *tls.Config {
	return client.tlsConfig
}

// SetTLSClientConfig sets the tls.Config
func (client *SocketClient) SetTLSClientConfig(config *tls.Config) {
	client.tlsConfig = config
}

// SetFullDuplex sets full duplex or half duplex mode of hprose socket client
func (client *SocketClient) SetFullDuplex(fullDuplex bool) {
	if fullDuplex {
		client.SendAndReceive = client.fullDuplexSendAndReceive
	} else {
		client.SendAndReceive = client.halfDuplexSendAndReceive
	}
}

// MaxPoolSize returns the max conn pool size of hprose socket client
func (client *SocketClient) MaxPoolSize() int {
	return cap(client.pool)
}

// SetMaxPoolSize sets the max conn pool size of hprose socket client
func (client *SocketClient) SetMaxPoolSize(size int) {
	pool := make(chan *connEntry, size)
	for i := 0; i < len(client.pool); i++ {
		select {
		case pool <- <-client.pool:
		default:
		}
	}
	client.pool = pool
}

func (client *SocketClient) getConn() net.Conn {
	for {
		select {
		case entry, closed := <-client.pool:
			if !closed {
				panic(errClientIsAlreadyClosed)
			}
			if entry.timer != nil {
				entry.timer.Stop()
				if entry.conn != nil {
					return entry.conn
				}
			}
			continue
		default:
			return nil
		}
	}
}

func (client *SocketClient) fetchConn() net.Conn {
	client.cond.L.Lock()
	for {
		conn := client.getConn()
		if conn != nil {
			client.cond.L.Unlock()
			return conn
		}
		if atomic.AddInt32(&client.connCount, 1) <= int32(cap(client.pool)) {
			client.cond.L.Unlock()
			return client.createConn()
		}
		atomic.AddInt32(&client.connCount, -1)
		client.cond.Wait()
	}
}

func ifErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Close the client
func (client *SocketClient) Close() {
	close(client.pool)
}

func (client *SocketClient) fullDuplexSendAndReceive(
	data []byte, context *ClientContext) ([]byte, error) {
	id := atomic.AddUint32(&client.nextid, 1)
	buf := make([]byte, len(data)+4)
	fromUint32(buf, id)
	copy(buf[4:], data)
	//response := make(chan websocketResponse)
	return nil, nil
}

func (client *SocketClient) close(conn net.Conn) {
	conn.Close()
	atomic.AddInt32(&client.connCount, -1)
}

func (client *SocketClient) halfDuplexSendAndReceive(
	data []byte, context *ClientContext) (response []byte, err error) {
	conn := client.fetchConn()
	err = conn.SetDeadline(time.Now().Add(context.Timeout))
	dataPacket := packet{body: data}
	if err == nil {
		err = sendData(conn, dataPacket)
	}
	if err == nil {
		err = recvData(conn, &dataPacket)
	}
	if err == nil {
		err = conn.SetDeadline(time.Time{})
	}
	if err != nil {
		client.close(conn)
		client.cond.Signal()
		return
	}
	entry := &connEntry{conn: conn}
	entry.timer = time.AfterFunc(client.IdleTimeout, func() {
		client.close(entry.conn)
		entry.conn = nil
		entry.timer = nil
	})
	client.pool <- entry
	client.cond.Signal()
	return dataPacket.body, nil
}
