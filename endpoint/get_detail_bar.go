package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/service/interface"
)

func makeGetDetailBarEndpoint(usecase _interface.BarService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.GetDetailBar(ctx, request.(*pb.Select))
		return res, err
	}
}

func (e BarEndpoint) GetDetailBar(ctx context.Context, req *pb.Select) (*pb.Bar, error) {
	res, err := e.GetDetailBarEndpoint(ctx, req)
	if err != nil {
		return res.(*pb.Bar), err
	}
	return res.(*pb.Bar), nil
}
