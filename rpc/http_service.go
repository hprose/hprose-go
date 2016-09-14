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
 * rpc/http_service.go                                    *
 *                                                        *
 * hprose http service for Go.                            *
 *                                                        *
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hprose/hprose-golang/util"
)

// HTTPContext is the hprose http context
type HTTPContext struct {
	*ServiceContext
	Response http.ResponseWriter
	Request  *http.Request
}

// HTTPService is the hprose http service
type HTTPService struct {
	*BaseService
	P3P                          bool
	GET                          bool
	CrossDomain                  bool
	accessControlAllowOrigins    map[string]bool
	lastModified                 string
	etag                         string
	crossDomainXMLFile           string
	crossDomainXMLContent        []byte
	clientAccessPolicyXMLFile    string
	clientAccessPolicyXMLContent []byte
}

type sendHeaderEvent interface {
	OnSendHeader(context *HTTPContext)
}

type sendHeaderEvent2 interface {
	OnSendHeader(context *HTTPContext) error
}

type httpFixer struct{}

func (httpFixer) FixArguments(args []reflect.Value, context *ServiceContext) {
	i := len(args) - 1
	typ := args[i].Type()
	if typ == httpContextType {
		if c, ok := context.TransportContext.(*HTTPContext); ok {
			args[i] = reflect.ValueOf(c)
		}
		return
	}
	if typ == httpRequestType {
		if c, ok := context.TransportContext.(*HTTPContext); ok {
			args[i] = reflect.ValueOf(c.Request)
		}
		return
	}
	fixArguments(args, context)
}

// NewHTTPService is the constructor of HTTPService
func NewHTTPService() (service *HTTPService) {
	t := time.Now().UTC()
	rand.Seed(t.UnixNano())
	service = new(HTTPService)
	service.BaseService = NewBaseService()
	service.P3P = true
	service.GET = true
	service.CrossDomain = true
	service.accessControlAllowOrigins = make(map[string]bool)
	service.lastModified = t.Format(time.RFC1123)
	service.etag = `"` + strconv.FormatInt(rand.Int63(), 16) + `"`
	service.fixer = httpFixer{}
	return
}

func (service *HTTPService) crossDomainXMLHandler(
	response http.ResponseWriter, request *http.Request) bool {
	if service.crossDomainXMLContent == nil ||
		strings.ToLower(request.URL.Path) != "/crossdomain.xml" {
		return false
	}
	if request.Header.Get("if-modified-since") == service.lastModified &&
		request.Header.Get("if-none-match") == service.etag {
		response.WriteHeader(304)
	} else {
		contentLength := len(service.crossDomainXMLContent)
		header := response.Header()
		header.Set("Last-Modified", service.lastModified)
		header.Set("Etag", service.etag)
		header.Set("Content-Type", "text/xml")
		header.Set("Content-Length", util.Itoa(contentLength))
		response.Write(service.crossDomainXMLContent)
	}
	return true
}

func (service *HTTPService) clientAccessPolicyXMLHandler(
	response http.ResponseWriter, request *http.Request) bool {
	if service.clientAccessPolicyXMLContent == nil ||
		strings.ToLower(request.URL.Path) != "/clientaccesspolicy.xml" {
		return false
	}
	if request.Header.Get("if-modified-since") == service.lastModified &&
		request.Header.Get("if-none-match") == service.etag {
		response.WriteHeader(304)
	} else {
		contentLength := len(service.clientAccessPolicyXMLContent)
		header := response.Header()
		header.Set("Last-Modified", service.lastModified)
		header.Set("Etag", service.etag)
		header.Set("Content-Type", "text/xml")
		header.Set("Content-Length", util.Itoa(contentLength))
		response.Write(service.clientAccessPolicyXMLContent)
	}
	return true
}

func (service *HTTPService) fireSendHeaderEvent(
	context *HTTPContext) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = NewPanicError(e)
		}
	}()
	switch event := service.Event.(type) {
	case sendHeaderEvent:
		event.OnSendHeader(context)
	case sendHeaderEvent2:
		err = event.OnSendHeader(context)
	}
	return err
}

func (service *HTTPService) sendHeader(context *HTTPContext) (err error) {
	if err = service.fireSendHeaderEvent(context); err != nil {
		return err
	}
	header := context.Response.Header()
	header.Set("Content-Type", "text/plain")
	if service.P3P {
		header.Set("P3P", `CP="CAO DSP COR CUR ADM DEV TAI PSA PSD IVAi IVDi `+
			`CONi TELo OTPi OUR DELi SAMi OTRi UNRi PUBi IND PHY ONL `+
			`UNI PUR FIN COM NAV INT DEM CNT STA POL HEA PRE GOV"`)
	}
	if service.CrossDomain {
		origin := context.Request.Header.Get("origin")
		if origin != "" && origin != "null" {
			if len(service.accessControlAllowOrigins) == 0 || service.accessControlAllowOrigins[origin] {
				header.Set("Access-Control-Allow-Origin", origin)
				header.Set("Access-Control-Allow-Credentials", "true")
			}
		} else {
			header.Set("Access-Control-Allow-Origin", "*")
		}
	}
	return nil
}

// AddAccessControlAllowOrigin add access control allow origin
func (service *HTTPService) AddAccessControlAllowOrigin(origins ...string) {
	for _, origin := range origins {
		service.accessControlAllowOrigins[origin] = true
	}
}

// RemoveAccessControlAllowOrigin remove access control allow origin
func (service *HTTPService) RemoveAccessControlAllowOrigin(origins ...string) {
	for _, origin := range origins {
		delete(service.accessControlAllowOrigins, origin)
	}
}

// CrossDomainXMLFile return the cross domain xml file
func (service *HTTPService) CrossDomainXMLFile() string {
	return service.crossDomainXMLFile
}

// CrossDomainXMLContent return the cross domain xml content
func (service *HTTPService) CrossDomainXMLContent() []byte {
	return service.crossDomainXMLContent
}

// ClientAccessPolicyXMLFile return the client access policy xml file
func (service *HTTPService) ClientAccessPolicyXMLFile() string {
	return service.clientAccessPolicyXMLFile
}

// ClientAccessPolicyXMLContent return the client access policy xml content
func (service *HTTPService) ClientAccessPolicyXMLContent() []byte {
	return service.clientAccessPolicyXMLContent
}

// SetCrossDomainXMLFile set the cross domain xml file
func (service *HTTPService) SetCrossDomainXMLFile(filename string) {
	service.crossDomainXMLFile = filename
	service.crossDomainXMLContent, _ = ioutil.ReadFile(filename)
}

// SetClientAccessPolicyXMLFile set the client access policy xml file
func (service *HTTPService) SetClientAccessPolicyXMLFile(filename string) {
	service.clientAccessPolicyXMLFile = filename
	service.clientAccessPolicyXMLContent, _ = ioutil.ReadFile(filename)
}

// SetCrossDomainXMLContent set the cross domain xml content
func (service *HTTPService) SetCrossDomainXMLContent(content []byte) {
	service.crossDomainXMLFile = ""
	service.crossDomainXMLContent = content
}

// SetClientAccessPolicyXMLContent set the client access policy xml content
func (service *HTTPService) SetClientAccessPolicyXMLContent(content []byte) {
	service.clientAccessPolicyXMLFile = ""
	service.clientAccessPolicyXMLContent = content
}

func readAllFromHTTPRequest(request *http.Request) ([]byte, error) {
	if request.ContentLength > 0 {
		data := make([]byte, request.ContentLength)
		_, err := io.ReadFull(request.Body, data)
		return data, err
	}
	if request.ContentLength < 0 {
		return ioutil.ReadAll(request.Body)
	}
	return nil, nil
}

// Serve is the hprose http handler method with the userData
func (service *HTTPService) Serve(
	response http.ResponseWriter, request *http.Request,
	userData map[string]interface{}) {
	if service.clientAccessPolicyXMLHandler(response, request) ||
		service.crossDomainXMLHandler(response, request) {
		return
	}
	context := new(HTTPContext)
	context.ServiceContext = NewServiceContext(nil)
	context.ServiceContext.TransportContext = context
	context.Response = response
	context.Request = request
	if userData != nil {
		for k, v := range userData {
			context.SetInterface(k, v)
		}
	}
	var resp []byte
	err := service.sendHeader(context)
	if err == nil {
		switch request.Method {
		case "GET":
			if service.GET {
				resp = service.doFunctionList(context.ServiceContext)
			} else {
				response.WriteHeader(403)
			}
		case "POST":
			req, err := readAllFromHTTPRequest(request)
			request.Body.Close()
			if err == nil {
				resp, err = service.Handle(req, context.ServiceContext)
			}
		}
	}
	if err != nil {
		resp = service.endError(err, context)
	}
	response.Header().Set("Content-Length", util.Itoa(len(resp)))
	response.Write(resp)
}

// ServeHTTP is the hprose http handler method
func (service *HTTPService) ServeHTTP(
	response http.ResponseWriter, request *http.Request) {
	service.Serve(response, request, nil)
}
