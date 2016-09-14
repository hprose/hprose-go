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
 * rpc/fasthttp_service.go                                *
 *                                                        *
 * hprose fasthttp service for Go.                        *
 *                                                        *
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"io/ioutil"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hprose/hprose-golang/util"
	"github.com/valyala/fasthttp"
)

// FastHTTPContext is the hprose fasthttp context
type FastHTTPContext struct {
	*ServiceContext
	RequestCtx *fasthttp.RequestCtx
}

// FastHTTPService is the hprose fasthttp service
type FastHTTPService HTTPService

type fastSendHeaderEvent interface {
	OnSendHeader(context *FastHTTPContext)
}

type fastSendHeaderEvent2 interface {
	OnSendHeader(context *FastHTTPContext) error
}

type fasthttpFixer struct{}

func (fasthttpFixer) FixArguments(args []reflect.Value, context *ServiceContext) {
	i := len(args) - 1
	typ := args[i].Type()
	if typ == fasthttpContextType {
		if c, ok := context.TransportContext.(*FastHTTPContext); ok {
			args[i] = reflect.ValueOf(c)
		}
		return
	}
	if typ == fasthttpRequestCtxType {
		if c, ok := context.TransportContext.(*FastHTTPContext); ok {
			args[i] = reflect.ValueOf(c.RequestCtx)
		}
		return
	}
	fixArguments(args, context)
}

// NewFastHTTPService is the constructor of FastHTTPService
func NewFastHTTPService() (service *FastHTTPService) {
	t := time.Now().UTC()
	rand.Seed(t.UnixNano())
	service = new(FastHTTPService)
	service.BaseService = NewBaseService()
	service.P3P = true
	service.GET = true
	service.CrossDomain = true
	service.accessControlAllowOrigins = make(map[string]bool)
	service.lastModified = t.Format(time.RFC1123)
	service.etag = `"` + strconv.FormatInt(rand.Int63(), 16) + `"`
	service.fixer = fasthttpFixer{}
	return
}

func (service *FastHTTPService) crossDomainXMLHandler(
	ctx *fasthttp.RequestCtx) bool {
	path := "/crossdomain.xml"
	if service.crossDomainXMLContent == nil ||
		strings.ToLower(util.ByteString(ctx.Path())) != path {
		return false
	}
	header := ctx.Request.Header
	ifModifiedSince := util.ByteString(header.Peek("if-modified-since"))
	ifNoneMatch := util.ByteString(header.Peek("if-none-match"))
	if ifModifiedSince == service.lastModified && ifNoneMatch == service.etag {
		ctx.SetStatusCode(304)
	} else {
		contentLength := len(service.crossDomainXMLContent)
		ctx.Response.Header.Set("Last-Modified", service.lastModified)
		ctx.Response.Header.Set("Etag", service.etag)
		ctx.Response.Header.SetContentType("text/xml")
		ctx.Response.Header.Set("Content-Length", util.Itoa(contentLength))
		ctx.SetBody(service.crossDomainXMLContent)
	}
	return true
}

func (service *FastHTTPService) clientAccessPolicyXMLHandler(
	ctx *fasthttp.RequestCtx) bool {
	path := "/clientaccesspolicy.xml"
	if service.clientAccessPolicyXMLContent == nil ||
		strings.ToLower(util.ByteString(ctx.Path())) != path {
		return false
	}
	header := ctx.Request.Header
	ifModifiedSince := util.ByteString(header.Peek("if-modified-since"))
	ifNoneMatch := util.ByteString(header.Peek("if-none-match"))
	if ifModifiedSince == service.lastModified && ifNoneMatch == service.etag {
		ctx.SetStatusCode(304)
	} else {
		contentLength := len(service.clientAccessPolicyXMLContent)
		ctx.Response.Header.Set("Last-Modified", service.lastModified)
		ctx.Response.Header.Set("Etag", service.etag)
		ctx.Response.Header.SetContentType("text/xml")
		ctx.Response.Header.Set("Content-Length", util.Itoa(contentLength))
		ctx.SetBody(service.clientAccessPolicyXMLContent)
	}
	return true
}

func (service *FastHTTPService) fireSendHeaderEvent(
	context *FastHTTPContext) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = NewPanicError(e)
		}
	}()
	switch event := service.Event.(type) {
	case fastSendHeaderEvent:
		event.OnSendHeader(context)
	case fastSendHeaderEvent2:
		err = event.OnSendHeader(context)
	}
	return err
}

func (service *FastHTTPService) sendHeader(
	context *FastHTTPContext) (err error) {
	if err = service.fireSendHeaderEvent(context); err != nil {
		return err
	}
	ctx := context.RequestCtx
	ctx.Response.Header.Set("Content-Type", "text/plain")
	if service.P3P {
		ctx.Response.Header.Set("P3P", `CP="CAO DSP COR CUR ADM DEV TAI PSA PSD IVAi IVDi `+
			`CONi TELo OTPi OUR DELi SAMi OTRi UNRi PUBi IND PHY ONL `+
			`UNI PUR FIN COM NAV INT DEM CNT STA POL HEA PRE GOV"`)
	}
	if service.CrossDomain {
		origin := util.ByteString(ctx.Request.Header.Peek("origin"))
		if origin != "" && origin != "null" {
			if len(service.accessControlAllowOrigins) == 0 || service.accessControlAllowOrigins[origin] {
				ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
				ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			}
		} else {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		}
	}
	return nil
}

// AddAccessControlAllowOrigin add access control allow origin
func (service *FastHTTPService) AddAccessControlAllowOrigin(origins ...string) {
	for _, origin := range origins {
		service.accessControlAllowOrigins[origin] = true
	}
}

// RemoveAccessControlAllowOrigin remove access control allow origin
func (service *FastHTTPService) RemoveAccessControlAllowOrigin(origins ...string) {
	for _, origin := range origins {
		delete(service.accessControlAllowOrigins, origin)
	}
}

// CrossDomainXMLFile return the cross domain xml file
func (service *FastHTTPService) CrossDomainXMLFile() string {
	return service.crossDomainXMLFile
}

// CrossDomainXMLContent return the cross domain xml content
func (service *FastHTTPService) CrossDomainXMLContent() []byte {
	return service.crossDomainXMLContent
}

// ClientAccessPolicyXMLFile return the client access policy xml file
func (service *FastHTTPService) ClientAccessPolicyXMLFile() string {
	return service.clientAccessPolicyXMLFile
}

// ClientAccessPolicyXMLContent return the client access policy xml content
func (service *FastHTTPService) ClientAccessPolicyXMLContent() []byte {
	return service.clientAccessPolicyXMLContent
}

// SetCrossDomainXMLFile set the cross domain xml file
func (service *FastHTTPService) SetCrossDomainXMLFile(filename string) {
	service.crossDomainXMLFile = filename
	service.crossDomainXMLContent, _ = ioutil.ReadFile(filename)
}

// SetClientAccessPolicyXMLFile set the client access policy xml file
func (service *FastHTTPService) SetClientAccessPolicyXMLFile(filename string) {
	service.clientAccessPolicyXMLFile = filename
	service.clientAccessPolicyXMLContent, _ = ioutil.ReadFile(filename)
}

// SetCrossDomainXMLContent set the cross domain xml content
func (service *FastHTTPService) SetCrossDomainXMLContent(content []byte) {
	service.crossDomainXMLFile = ""
	service.crossDomainXMLContent = content
}

// SetClientAccessPolicyXMLContent set the client access policy xml content
func (service *FastHTTPService) SetClientAccessPolicyXMLContent(content []byte) {
	service.clientAccessPolicyXMLFile = ""
	service.clientAccessPolicyXMLContent = content
}

// Serve is the hprose fasthttp handler method with the userData
func (service *FastHTTPService) Serve(
	ctx *fasthttp.RequestCtx, userData map[string]interface{}) {
	if service.clientAccessPolicyXMLHandler(ctx) ||
		service.crossDomainXMLHandler(ctx) {
		return
	}
	context := new(FastHTTPContext)
	context.ServiceContext = NewServiceContext(nil)
	context.ServiceContext.TransportContext = context
	context.RequestCtx = ctx
	if userData != nil {
		for k, v := range userData {
			context.SetInterface(k, v)
		}
	}
	var resp []byte
	err := service.sendHeader(context)
	if err == nil {
		switch util.ByteString(ctx.Method()) {
		case "GET":
			if service.GET {
				resp = service.doFunctionList(context.ServiceContext)
			} else {
				ctx.SetStatusCode(403)
			}
		case "POST":
			req := ctx.Request.Body()
			resp, err = service.Handle(req, context.ServiceContext)
		}
	}
	if err != nil {
		resp = service.endError(err, context)
	}
	context.RequestCtx = nil
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(resp)))
	ctx.SetBody(resp)
}

// ServeFastHTTP is the hprose fasthttp handler method
func (service *FastHTTPService) ServeFastHTTP(ctx *fasthttp.RequestCtx) {
	service.Serve(ctx, nil)
}
