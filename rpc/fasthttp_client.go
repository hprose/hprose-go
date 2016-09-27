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
 * LastModified: Sep 27, 2016                             *
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
	*BaseClient
	*fasthttp.Client
	Header      *fasthttp.RequestHeader
	compression bool
	keepAlive   bool
}

// NewFastHTTPClient is the constructor of FastHTTPClient
func NewFastHTTPClient(uri ...string) (client *FastHTTPClient) {
	client = new(FastHTTPClient)
	client.BaseClient = NewBaseClient()
	client.Client = new(fasthttp.Client)
	client.Header = new(fasthttp.RequestHeader)
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
	checkURLList(client, uriList)
	client.BaseClient.SetURIList(uriList)
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
	if err := client.Client.DoTimeout(req, resp, context.Timeout); err != nil {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}
	body := resp.Body()
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return body, nil
}
