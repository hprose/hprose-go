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
 * LastModified: Oct 8, 2016                              *
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

type socketResponse struct {
	data []byte
	err  error
}

type connEntry struct {
	conn      net.Conn
	timer     *time.Timer
	reqCount  int32
	cond      *sync.Cond
	responses map[uint32]chan socketResponse
}

// SocketClient is base struct for TCPClient and UnixClient
type SocketClient struct {
	BaseClient
	ReadBuffer  int
	WriteBuffer int
	IdleTimeout time.Duration
	TLSConfig   *tls.Config
	connPool    chan *connEntry
	connCount   int32
	nextid      uint32
	createConn  func() net.Conn
	cond        sync.Cond
}

func (client *SocketClient) initSocketClient() {
	client.initBaseClient()
	client.ReadBuffer = 0
	client.WriteBuffer = 0
	client.IdleTimeout = 30 * time.Second
	client.TLSConfig = nil
	client.connPool = make(chan *connEntry, runtime.NumCPU()*2)
	client.connCount = 0
	client.nextid = 0
	client.cond.L = &sync.Mutex{}
	client.SetFullDuplex(false)
}

// TLSClientConfig returns the tls.Config in hprose client
func (client *SocketClient) TLSClientConfig() *tls.Config {
	return client.TLSConfig
}

// SetTLSClientConfig sets the tls.Config
func (client *SocketClient) SetTLSClientConfig(config *tls.Config) {
	client.TLSConfig = config
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
	return cap(client.connPool)
}

// SetMaxPoolSize sets the max conn pool size of hprose socket client
func (client *SocketClient) SetMaxPoolSize(size int) {
	pool := make(chan *connEntry, size)
	for i := 0; i < len(client.connPool); i++ {
		select {
		case pool <- <-client.connPool:
		default:
		}
	}
	client.connPool = pool
}

func (client *SocketClient) getConn() *connEntry {
	for {
		select {
		case entry, closed := <-client.connPool:
			if !closed {
				panic(errClientIsAlreadyClosed)
			}
			if entry.timer != nil {
				entry.timer.Stop()
			}
			if entry.conn != nil {
				return entry
			}
			continue
		default:
			return nil
		}
	}
}

func (client *SocketClient) fullDuplexReceive(entry *connEntry) {
	conn := entry.conn
	var dataPacket packet
	for {
		err := recvData(conn, &dataPacket)
		if err != nil {
			if entry.responses != nil {
				entry.cond.L.Lock()
				responses := entry.responses
				entry.conn = nil
				entry.reqCount = 0
				entry.responses = nil
				entry.cond.L.Unlock()
				entry.cond.Broadcast()
				client.close(conn)
				for _, response := range responses {
					response <- socketResponse{nil, err}
				}
			}
			break
		}
		id := toUint32(dataPacket.id[:])
		entry.cond.L.Lock()
		response := entry.responses[id]
		delete(entry.responses, id)
		entry.reqCount--
		entry.cond.L.Unlock()
		entry.cond.Signal()
		if response != nil {
			response <- socketResponse{dataPacket.body, nil}
		}
	}
}

func (client *SocketClient) fetchConn(fullDuplex bool) *connEntry {
	client.cond.L.Lock()
	for {
		entry := client.getConn()
		if entry != nil && entry.conn != nil {
			client.cond.L.Unlock()
			return entry
		}
		if int(atomic.AddInt32(&client.connCount, 1)) <= cap(client.connPool) {
			client.cond.L.Unlock()
			entry := &connEntry{conn: client.createConn()}
			if fullDuplex {
				entry.cond = sync.NewCond(&sync.Mutex{})
				entry.responses = make(map[uint32]chan socketResponse, 10)
				go client.fullDuplexReceive(entry)
			}
			return entry
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
	close(client.connPool)
}

func (client *SocketClient) close(conn net.Conn) {
	conn.Close()
	atomic.AddInt32(&client.connCount, -1)
}

func (client *SocketClient) fullDuplexSendAndReceive(
	data []byte, context *ClientContext) (resp []byte, err error) {
	var entry *connEntry
	for {
		entry = client.fetchConn(true)
		entry.cond.L.Lock()
		for entry.reqCount > 10 {
			entry.cond.Wait()
		}
		entry.cond.L.Unlock()
		if entry.conn != nil {
			break
		}
		entry.reqCount = 0
		entry.cond.Signal()
	}
	conn := entry.conn
	id := atomic.AddUint32(&client.nextid, 1)
	deadline := time.Now().Add(context.Timeout)
	err = conn.SetDeadline(deadline)
	response := make(chan socketResponse)
	if err == nil {
		entry.cond.L.Lock()
		entry.responses[id] = response
		entry.reqCount++
		entry.cond.L.Unlock()
		dataPacket := packet{fullDuplex: true, body: data}
		fromUint32(dataPacket.id[:], id)
		err = sendData(conn, dataPacket)
	}
	if err == nil {
		err = conn.SetDeadline(time.Time{})
	}
	if err != nil {
		client.close(conn)
		client.cond.Signal()
		return
	}
	client.connPool <- entry
	client.cond.Signal()
	select {
	case resp := <-response:
		return resp.data, resp.err
	case <-time.After(deadline.Sub(time.Now())):
		entry.cond.L.Lock()
		delete(entry.responses, id)
		entry.reqCount--
		entry.cond.L.Unlock()
		entry.cond.Signal()
		return nil, ErrTimeout
	}
}

func (client *SocketClient) halfDuplexSendAndReceive(
	data []byte, context *ClientContext) ([]byte, error) {
	entry := client.fetchConn(false)
	conn := entry.conn
	err := conn.SetDeadline(time.Now().Add(context.Timeout))
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
		return nil, err
	}
	if entry.timer == nil {
		entry.timer = time.AfterFunc(client.IdleTimeout, func() {
			client.close(conn)
			entry.conn = nil
			entry.timer = nil
		})
	} else {
		entry.timer.Reset(client.IdleTimeout)
	}
	client.connPool <- entry
	client.cond.Signal()
	return dataPacket.body, nil
}
