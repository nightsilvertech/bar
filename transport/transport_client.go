package transport

import (
	"context"
	grpcgoogle "google.golang.org/grpc"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	ep "github.com/nightsilvertech/bar/endpoint"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/service/interface"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

func encodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

func DialBarService(ctx context.Context, hostAndPort string) (_interface.BarService, *grpcgoogle.ClientConn, error) {
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		return nil, nil, err
	}
	dialOptions := []grpcgoogle.DialOption{
		grpcgoogle.WithInsecure(),
		grpcgoogle.WithStatsHandler(new(ocgrpc.ClientHandler)),
	}
	conn, err := grpcgoogle.DialContext(ctx, hostAndPort, dialOptions...)
	if err != nil {
		panic(err)
	}
	return newGRPBarClient(conn), conn, nil
}

func newGRPBarClient(conn *grpc.ClientConn) _interface.BarService {
	var addBarEp endpoint.Endpoint
	{
		const (
			rpcName   = `api.v1.ProgramService`
			rpcMethod = `AddProgram`
		)

		addBarEp = grpctransport.NewClient(
			conn,
			rpcName,
			rpcMethod,
			encodeRequest,
			decodeResponse,
			pb.Bar{},
		).Endpoint()
	}
	return &ep.BarEndpoint{
		AddBarEndpoint: addBarEp,
	}
}
