package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/service/interface"
)

func makeAddBarEndpoint(usecase _interface.BarService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.AddBar(ctx, request.(*pb.Bar))
		return res, err
	}
}

func (e BarEndpoint) AddBar(ctx context.Context, req *pb.Bar) (*pb.Bar, error) {
	res, err := e.AddBarEndpoint(ctx, req)
	if err != nil {
		return res.(*pb.Bar), err
	}
	return res.(*pb.Bar), nil
}

