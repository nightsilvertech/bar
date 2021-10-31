package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/nightsilvertech/bar/constant"
	"github.com/nightsilvertech/bar/middleware"
	_interface "github.com/nightsilvertech/bar/service/interface"
)

type BarEndpoint struct {
	AddBarEndpoint       endpoint.Endpoint
	EditBarEndpoint      endpoint.Endpoint
	DeleteBarEndpoint    endpoint.Endpoint
	GetAllBarEndpoint    endpoint.Endpoint
	GetDetailBarEndpoint endpoint.Endpoint
}

func NewBarEndpoint(svc _interface.BarService) BarEndpoint {
	var addBarEp endpoint.Endpoint
	{
		const name = `AddBar`
		addBarEp = makeAddBarEndpoint(svc)
		addBarEp = middleware.CircuitBreakerMiddleware(constant.ServiceName)(addBarEp)
	}

	var editBarEp endpoint.Endpoint
	{
		const name = `EditBar`
		editBarEp = makeEditBarEndpoint(svc)
		editBarEp = middleware.CircuitBreakerMiddleware(constant.ServiceName)(editBarEp)
	}

	var deleteBarEp endpoint.Endpoint
	{
		const name = `DeleteBar`
		deleteBarEp = makeDeleteBarEndpoint(svc)
		deleteBarEp = middleware.CircuitBreakerMiddleware(constant.ServiceName)(editBarEp)
	}

	var getAllBarEp endpoint.Endpoint
	{
		const name = `GetAllBar`
		getAllBarEp = makeGetAllBarEndpoint(svc)
		getAllBarEp = middleware.CircuitBreakerMiddleware(constant.ServiceName)(getAllBarEp)
	}

	var getDetailBarEp endpoint.Endpoint
	{
		const name = `GetDetailBar`
		getDetailBarEp = makeGetDetailBarEndpoint(svc)
		getDetailBarEp = middleware.CircuitBreakerMiddleware(constant.ServiceName)(getDetailBarEp)
	}

	return BarEndpoint{
		AddBarEndpoint:       addBarEp,
		EditBarEndpoint:      editBarEp,
		DeleteBarEndpoint:    deleteBarEp,
		GetAllBarEndpoint:    getAllBarEp,
		GetDetailBarEndpoint: getDetailBarEp,
	}
}
