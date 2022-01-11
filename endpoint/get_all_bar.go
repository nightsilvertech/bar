package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/service/interface"
)

func makeGetAllBarEndpoint(usecase _interface.BarService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := usecase.GetAllBar(ctx, request.(*pb.Pagination))
		return res, err
	}
}

func (e BarEndpoint) GetAllBar(ctx context.Context, req *pb.Pagination) (*pb.Bars, error) {
	res, err := e.GetAllBarEndpoint(ctx, req)
	if err != nil {
		return &pb.Bars{}, err
	}
	return res.(*pb.Bars), nil
}

