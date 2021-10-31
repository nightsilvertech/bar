package _interface

import (
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
)

type BarService interface {
	pb.BarServiceServer
}
