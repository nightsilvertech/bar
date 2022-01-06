package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	"github.com/nightsilvertech/bar/constant"
	_interface "github.com/nightsilvertech/bar/service/interface"
	"github.com/nightsilvertech/utl/middlewares"
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
		addBarEp = middlewares.CircuitBreakerMiddleware(constant.ServiceName)(addBarEp)
		addBarEp = kitoc.TraceEndpoint(name)(addBarEp)
	}

	var editBarEp endpoint.Endpoint
	{
		const name = `EditBar`
		editBarEp = makeEditBarEndpoint(svc)
		editBarEp = middlewares.CircuitBreakerMiddleware(constant.ServiceName)(editBarEp)
		editBarEp = kitoc.TraceEndpoint(name)(editBarEp)
	}

	var deleteBarEp endpoint.Endpoint
	{
		const name = `DeleteBar`
		deleteBarEp = makeDeleteBarEndpoint(svc)
		deleteBarEp = middlewares.CircuitBreakerMiddleware(constant.ServiceName)(deleteBarEp)
		deleteBarEp = kitoc.TraceEndpoint(name)(deleteBarEp)
	}

	var getAllBarEp endpoint.Endpoint
	{
		const name = `GetAllBar`
		getAllBarEp = makeGetAllBarEndpoint(svc)
		getAllBarEp = middlewares.CircuitBreakerMiddleware(constant.ServiceName)(getAllBarEp)
		getAllBarEp = kitoc.TraceEndpoint(name)(getAllBarEp)
	}

	var getDetailBarEp endpoint.Endpoint
	{
		const name = `GetDetailBar`
		getDetailBarEp = makeGetDetailBarEndpoint(svc)
		getDetailBarEp = middlewares.CircuitBreakerMiddleware(constant.ServiceName)(getDetailBarEp)
		getDetailBarEp = kitoc.TraceEndpoint(name)(getDetailBarEp)
	}

	return BarEndpoint{
		AddBarEndpoint:       addBarEp,
		EditBarEndpoint:      editBarEp,
		DeleteBarEndpoint:    deleteBarEp,
		GetAllBarEndpoint:    getAllBarEp,
		GetDetailBarEndpoint: getDetailBarEp,
	}
}
