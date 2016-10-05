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
 * rpc/service_context.go                                 *
 *                                                        *
 * hprose service context for Go.                         *
 *                                                        *
 * LastModified: Oct 5, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

// ServiceContext is the hprose service context
type ServiceContext struct {
	BaseContext
	*Method
	Service
	TransportContext Context
	IsMissingMethod  bool
	ByRef            bool
}

func (context *ServiceContext) initServiceContext(service Service) {
	context.initBaseContext()
	context.Method = nil
	context.Service = service
	context.IsMissingMethod = false
	context.ByRef = false
}

// NewServiceContext is the constructor of ServiceContext
func NewServiceContext(service Service) (context *ServiceContext) {
	context = new(ServiceContext)
	context.initServiceContext(service)
	return
}
