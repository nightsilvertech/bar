package transport

import (
	"context"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	ep "github.com/nightsilvertech/bar/endpoint"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	"github.com/nightsilvertech/utl/console"
)

type grpcBarServer struct {
	pb.UnimplementedBarServiceServer
	addBar       grpctransport.Handler
	editBar      grpctransport.Handler
	deleteBar    grpctransport.Handler
	getAllBar    grpctransport.Handler
	getDetailBar grpctransport.Handler
}

func decodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func encodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

func (g *grpcBarServer) AddBar(ctx context.Context, bar *pb.Bar) (*pb.Bar, error) {
	_, res, err := g.addBar.ServeGRPC(ctx, bar)
	if err != nil {
		return nil, err
	}
	return res.(*pb.Bar), nil
}

func (g *grpcBarServer) EditBar(ctx context.Context, bar *pb.Bar) (*pb.Bar, error) {
	_, res, err := g.editBar.ServeGRPC(ctx, bar)
	if err != nil {
		return nil, err
	}
	return res.(*pb.Bar), nil
}

func (g *grpcBarServer) DeleteBar(ctx context.Context, selects *pb.Select) (*pb.Bar, error) {
	_, res, err := g.deleteBar.ServeGRPC(ctx, selects)
	if err != nil {
		return nil, err
	}
	return res.(*pb.Bar), nil
}

func (g *grpcBarServer) GetAllBar(ctx context.Context, pagination *pb.Pagination) (*pb.Bars, error) {
	_, res, err := g.getAllBar.ServeGRPC(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return res.(*pb.Bars), nil
}

func (g *grpcBarServer) GetDetailBar(ctx context.Context, selects *pb.Select) (*pb.Bar, error) {
	_, res, err := g.getDetailBar.ServeGRPC(ctx, selects)
	if err != nil {
		return nil, err
	}
	return res.(*pb.Bar), nil
}

func NewBarServer(endpoints ep.BarEndpoint) pb.BarServiceServer {
	options := []grpctransport.ServerOption{
		kitoc.GRPCServerTrace(),
		grpctransport.ServerBefore(
			console.RequestIDMetadataToContext(),
		),
	}
	return &grpcBarServer{
		addBar: grpctransport.NewServer(
			endpoints.AddBarEndpoint,
			decodeRequest,
			encodeResponse,
			options...,
		),
		editBar: grpctransport.NewServer(
			endpoints.EditBarEndpoint,
			decodeRequest,
			encodeResponse,
			options...,
		),
		deleteBar: grpctransport.NewServer(
			endpoints.DeleteBarEndpoint,
			decodeRequest,
			encodeResponse,
			options...,
		),
		getAllBar: grpctransport.NewServer(
			endpoints.GetAllBarEndpoint,
			decodeRequest,
			encodeResponse,
			options...,
		),
		getDetailBar: grpctransport.NewServer(
			endpoints.GetDetailBarEndpoint,
			decodeRequest,
			encodeResponse,
			options...,
		),
	}
}
