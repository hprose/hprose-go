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
 * rpc/http_client.go                                     *
 *                                                        *
 * hprose http client for Go.                             *
 *                                                        *
 * LastModified: Oct 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"

	"github.com/valyala/fasthttp"
)

// FastHTTPClient is hprose fasthttp client
type FastHTTPClient struct {
	baseClient
	limiter
	fasthttp.Client
	Header      fasthttp.RequestHeader
	compression bool
	keepAlive   bool
}

// NewFastHTTPClient is the constructor of FastHTTPClient
func NewFastHTTPClient(uri ...string) (client *FastHTTPClient) {
	client = new(FastHTTPClient)
	client.initBaseClient()
	client.initLimiter()
	client.compression = false
	client.keepAlive = true
	client.SetURIList(uri)
	client.SendAndReceive = client.sendAndReceive
	return
}

func newFastHTTPClient(uri ...string) Client {
	return NewFastHTTPClient(uri...)
}

// SetURIList set a list of server addresses
func (client *FastHTTPClient) SetURIList(uriList []string) {
	if checkAddresses(uriList, httpSchemes) == "https" {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	client.baseClient.SetURIList(uriList)
}

// TLSClientConfig return the tls.Config in hprose client
func (client *FastHTTPClient) TLSClientConfig() *tls.Config {
	return client.TLSConfig
}

// SetTLSClientConfig set the tls.Config
func (client *FastHTTPClient) SetTLSClientConfig(config *tls.Config) {
	client.TLSConfig = config
}

// SetKeepAlive set the keepalive status of hprose client
func (client *FastHTTPClient) SetKeepAlive(enable bool) {
	client.keepAlive = enable
}

// Compression return the compression status of hprose client
func (client *FastHTTPClient) Compression() bool {
	return client.compression
}

// SetCompression set the compression status of hprose client
func (client *FastHTTPClient) SetCompression(enable bool) {
	client.compression = enable
}

func (client *FastHTTPClient) sendAndReceive(
	data []byte, context *ClientContext) ([]byte, error) {
	client.cond.L.Lock()
	client.limit()
	client.cond.L.Unlock()
	req := fasthttp.AcquireRequest()
	client.Header.CopyTo(&req.Header)
	req.Header.SetMethod("POST")
	req.SetRequestURI(client.uri)
	req.SetBody(data)
	req.Header.SetContentLength(len(data))
	req.Header.SetContentType("application/hprose")
	if client.keepAlive {
		req.Header.Set("Connection", "keep-alive")
	} else {
		req.Header.Set("Connection", "close")
	}
	if client.compression {
		req.Header.Set("Content-Encoding", "gzip")
	}
	resp := fasthttp.AcquireResponse()
	err := client.Client.DoTimeout(req, resp, context.Timeout)
	if err != nil {
		data = nil
	} else {
		data = resp.Body()
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	client.cond.L.Lock()
	client.unlimit()
	client.cond.L.Unlock()
	return data, err
}
