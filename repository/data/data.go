package data

import (
	"context"
	"database/sql"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/repository/interface"
	"sync"
)

var mutex = &sync.RWMutex{}

type dataReadWrite struct {
	db *sql.DB
}

func (d *dataReadWrite) WriteBar(ctx context.Context, req *pb.Bar) (res *pb.Bar, err error) {
	panic("implement me")
}

func (d *dataReadWrite) ModifyBar(ctx context.Context, req *pb.Bar) (res *pb.Bar, err error) {
	panic("implement me")
}

func (d *dataReadWrite) RemoveBar(ctx context.Context, req *pb.Select) (res *pb.Bar, err error) {
	panic("implement me")
}

func (d *dataReadWrite) ReadDetailBar(ctx context.Context, req *pb.Select) (res *pb.Bar, err error) {
	panic("implement me")
}

func (d *dataReadWrite) ReadAllBar(ctx context.Context, req *pb.Pagination) (res *pb.Bars, err error) {
	panic("implement me")
}

func NewDataReadWriter() _interface.DRW {
	return &dataReadWrite{

	}
}
