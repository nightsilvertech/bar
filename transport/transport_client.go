package transport

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	ep "github.com/nightsilvertech/bar/endpoint"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/service/interface"
	"github.com/nightsilvertech/utl/console"
	"github.com/nightsilvertech/utl/jsonwebtoken"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	grpcgoogle "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func encodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

func DialBarService(
	host, port string,
	transportCred credentials.TransportCredentials,
) (_interface.BarService, *grpcgoogle.ClientConn, error) {
	hostAndPort := fmt.Sprintf("%s:%s", host, port)
	dialOptions := []grpcgoogle.DialOption{
		grpcgoogle.WithStatsHandler(new(ocgrpc.ClientHandler)),
	}
	if transportCred != nil {
		dialOptions = append(dialOptions, grpcgoogle.WithTransportCredentials(transportCred))
	} else {
		dialOptions = append(dialOptions, grpcgoogle.WithInsecure())
	}

	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		return nil, nil, err
	}
	conn, err := grpcgoogle.Dial(hostAndPort, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	return newGRPBarClient(conn), conn, nil
}

func newGRPBarClient(conn *grpc.ClientConn) _interface.BarService {
	var addBarEp endpoint.Endpoint
	{
		const (
			rpcName   = `api.v1.BarService`
			rpcMethod = `AddBar`
		)

		addBarEp = grpctransport.NewClient(
			conn,
			rpcName,
			rpcMethod,
			encodeRequest,
			decodeResponse,
			pb.Bar{},
			grpctransport.ClientBefore(
				console.ContextToRequestIDMetadata(),
				jsonwebtoken.ContextToBearerTokenMetadata(),
			),
		).Endpoint()
	}

	var editBarEp endpoint.Endpoint
	{
		const (
			rpcName   = `api.v1.BarService`
			rpcMethod = `EditBar`
		)

		addBarEp = grpctransport.NewClient(
			conn,
			rpcName,
			rpcMethod,
			encodeRequest,
			decodeResponse,
			pb.Bar{},
			grpctransport.ClientBefore(
				console.ContextToRequestIDMetadata(),
				jsonwebtoken.ContextToBearerTokenMetadata(),
			),
		).Endpoint()
	}

	var deleteBarEp endpoint.Endpoint
	{
		const (
			rpcName   = `api.v1.BarService`
			rpcMethod = `DeleteBar`
		)

		addBarEp = grpctransport.NewClient(
			conn,
			rpcName,
			rpcMethod,
			encodeRequest,
			decodeResponse,
			pb.Bar{},
			grpctransport.ClientBefore(
				console.ContextToRequestIDMetadata(),
				jsonwebtoken.ContextToBearerTokenMetadata(),
			),
		).Endpoint()
	}

	var getDetailBarEp endpoint.Endpoint
	{
		const (
			rpcName   = `api.v1.BarService`
			rpcMethod = `GetDetailBar`
		)

		addBarEp = grpctransport.NewClient(
			conn,
			rpcName,
			rpcMethod,
			encodeRequest,
			decodeResponse,
			pb.Bar{},
			grpctransport.ClientBefore(
				console.ContextToRequestIDMetadata(),
				jsonwebtoken.ContextToBearerTokenMetadata(),
			),
		).Endpoint()
	}

	var getAllBarEp endpoint.Endpoint
	{
		const (
			rpcName   = `api.v1.BarService`
			rpcMethod = `GetAllBar`
		)

		addBarEp = grpctransport.NewClient(
			conn,
			rpcName,
			rpcMethod,
			encodeRequest,
			decodeResponse,
			pb.Bars{},
			grpctransport.ClientBefore(
				console.ContextToRequestIDMetadata(),
				jsonwebtoken.ContextToBearerTokenMetadata(),
			),
		).Endpoint()
	}

	return &ep.BarEndpoint{
		AddBarEndpoint:       addBarEp,
		EditBarEndpoint:      editBarEp,
		DeleteBarEndpoint:    deleteBarEp,
		GetDetailBarEndpoint: getDetailBarEp,
		GetAllBarEndpoint:    getAllBarEp,
	}
}
