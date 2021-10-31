package _interface

import (
	"context"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
)

// DRW means data read writer this interface
// contains all data management function
type DRW interface {
	WriteBar(ctx context.Context, req *pb.Bar) (res *pb.Bar, err error)
	ModifyBar(ctx context.Context, req *pb.Bar) (res *pb.Bar, err error)
	RemoveBar(ctx context.Context, req *pb.Select) (res *pb.Bar, err error)
	ReadDetailBar(ctx context.Context, req *pb.Select) (res *pb.Bar, err error)
	ReadAllBar(ctx context.Context, req *pb.Pagination) (res *pb.Bars, err error)
}

// CRW means cache read writer this interface
// contains all cache management function
type CRW interface {

}
