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
 * LastModified: Oct 2, 2016                              *
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
	context.Service = service
}

// NewServiceContext is the constructor of ServiceContext
func NewServiceContext(service Service) (context *ServiceContext) {
	context = new(ServiceContext)
	context.initServiceContext(service)
	return
}
