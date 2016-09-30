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
 * rpc/websocket_client.go                                *
 *                                                        *
 * hprose websocket client for Go.                        *
 *                                                        *
 * LastModified: Sep 30, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type websocketReqeust struct {
	id   uint32
	data []byte
}

type websocketResponse struct {
	data []byte
	err  error
}

// WebSocketClient is hprose websocket client
type WebSocketClient struct {
	*BaseClient
	*http.Header
	MaxConcurrentRequests int
	dialer                *websocket.Dialer
	cond                  *sync.Cond
	conn                  *websocket.Conn
	nextid                uint32
	requests              chan websocketReqeust
	responses             map[uint32]chan websocketResponse
}

// NewWebSocketClient is the constructor of WebSocketClient
func NewWebSocketClient(uri ...string) (client *WebSocketClient) {
	client = new(WebSocketClient)
	client.BaseClient = NewBaseClient()
	client.Header = new(http.Header)
	client.MaxConcurrentRequests = 10
	client.dialer = new(websocket.Dialer)
	client.cond = sync.NewCond(new(sync.Mutex))
	client.SetURIList(uri)
	client.SendAndReceive = client.sendAndReceive
	return
}

func newWebSocketClient(uri ...string) Client {
	return NewWebSocketClient(uri...)
}

func checkWebSocketAddresses(client Client, uriList []string) {
	for _, uri := range uriList {
		if u, err := url.Parse(uri); err == nil {
			if u.Scheme != "ws" && u.Scheme != "wss" {
				panic("This client desn't support " + u.Scheme + " scheme.")
			}
			if u.Scheme == "wss" {
				client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
			}
		}
	}
}

// SetURIList set a list of server addresses
func (client *WebSocketClient) SetURIList(uriList []string) {
	checkWebSocketAddresses(client, uriList)
	client.BaseClient.SetURIList(uriList)
}

func (client *WebSocketClient) close(err error) {
	client.cond.L.Lock()
	if err != nil && client.responses != nil {
		for _, response := range client.responses {
			response <- websocketResponse{nil, err}
		}
	}
	client.responses = nil
	if client.conn != nil {
		client.conn.Close()
		client.conn = nil
	}
	client.cond.Broadcast()
	client.cond.L.Unlock()
}

// Close the client
func (client *WebSocketClient) Close() {
	client.close(errClientIsAlreadyClosed)
}

// TLSClientConfig returns the tls.Config in hprose client
func (client *WebSocketClient) TLSClientConfig() *tls.Config {
	return client.dialer.TLSClientConfig
}

// SetTLSClientConfig sets the tls.Config
func (client *WebSocketClient) SetTLSClientConfig(config *tls.Config) {
	client.dialer.TLSClientConfig = config
}

func (client *WebSocketClient) sendLoop() {
	conn := client.conn
	for request := range client.requests {
		err := conn.WriteMessage(websocket.BinaryMessage, request.data)
		if err != nil {
			client.close(err)
			break
		}
	}
	client.requests = nil
}

func (client *WebSocketClient) recvLoop() {
	conn := client.conn
	count := client.MaxConcurrentRequests
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			client.close(err)
			break
		}
		if msgType == websocket.BinaryMessage {
			id := toUint32(data)
			client.cond.L.Lock()
			response := client.responses[id]
			if response != nil {
				response <- websocketResponse{data[4:], nil}
				delete(client.responses, id)
			}
			if len(client.responses) < count {
				client.cond.Signal()
			}
			client.cond.L.Unlock()
		}
	}
	close(client.requests)
}

func (client *WebSocketClient) getConn(uri string) (err error) {
	if client.conn == nil {
		client.conn, _, err = client.dialer.Dial(uri, *client.Header)
		if err != nil {
			return err
		}
		count := client.MaxConcurrentRequests
		client.requests = make(chan websocketReqeust, count)
		client.responses = make(map[uint32]chan websocketResponse, count)
		go client.sendLoop()
		go client.recvLoop()
	}
	return nil
}

func (client *WebSocketClient) sendAndReceive(
	data []byte, context *ClientContext) ([]byte, error) {
	id := atomic.AddUint32(&client.nextid, 1)
	buf := make([]byte, len(data)+4)
	fromUint32(buf, id)
	copy(buf[4:], data)
	response := make(chan websocketResponse)
	client.cond.L.Lock()
	for {
		if len(client.responses) < client.MaxConcurrentRequests {
			break
		}
		client.cond.Wait()
	}
	if err := client.getConn(client.uri); err != nil {
		client.cond.L.Unlock()
		return nil, err
	}
	client.responses[id] = response
	client.cond.L.Unlock()
	client.requests <- websocketReqeust{id, buf}
	select {
	case resp := <-response:
		return resp.data, resp.err
	case <-time.After(client.timeout):
		client.cond.L.Lock()
		delete(client.responses, id)
		if len(client.responses) < client.MaxConcurrentRequests {
			client.cond.Signal()
		}
		client.cond.L.Unlock()
		return nil, ErrTimeout
	}
}
