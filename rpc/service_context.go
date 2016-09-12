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
 * LastModified: Sep 12, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

// ServiceContext is the hprose service context
type ServiceContext struct {
	*BaseContext
	*Method
	Clients
	TransportContext Context
	IsMissingMethod  bool
	ByRef            bool
}

// NewServiceContext is the constructor of ServiceContext
func NewServiceContext(clients Clients) (context *ServiceContext) {
	context = new(ServiceContext)
	context.BaseContext = NewBaseContext()
	context.Clients = clients
	return
}
