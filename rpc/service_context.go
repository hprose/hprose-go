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
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

// ServiceContext is the hprose base context
type ServiceContext struct {
	*BaseContext
	method        *serviceMethod
	methods       *serviceMethods
	missingMethod bool
	byref         bool
	Clients
}

// NewServiceContext is the constructor of ServiceContext
func NewServiceContext(clients Clients) (context *ServiceContext) {
	context = new(ServiceContext)
	context.BaseContext = NewBaseContext()
	context.Clients = clients
	return
}
