package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	ep "github.com/nightsilvertech/bar/endpoint"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/service/interface"
	"github.com/nightsilvertech/utl/console"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	grpcgoogle "google.golang.org/grpc"
)

func encodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

func DialBarService(hostAndPort string) (_interface.BarService, *grpcgoogle.ClientConn, error) {
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		return nil, nil, err
	}
	dialOptions := []grpcgoogle.DialOption{
		grpcgoogle.WithInsecure(),
		grpcgoogle.WithStatsHandler(new(ocgrpc.ClientHandler)),
	}
	conn, err := grpcgoogle.Dial(hostAndPort, dialOptions...)
	if err != nil {
		panic(err)
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
			),
		).Endpoint()
	}
	return &ep.BarEndpoint{
		AddBarEndpoint: addBarEp,
	}
}
