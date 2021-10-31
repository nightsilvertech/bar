package service

import (
	"context"
	"github.com/nightsilvertech/bar/gvar"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_repo "github.com/nightsilvertech/bar/repository"
	_interface "github.com/nightsilvertech/bar/service/interface"
	"github.com/nightsilvertech/bar/util"
	uuid "github.com/satori/go.uuid"
)

type service struct {
	repo _repo.Repository
}

func (s service) AddBar(ctx context.Context, bar *pb.Bar) (*pb.Bar, error) {
	const funcName = `AddBar`
	bar.Id = uuid.NewV4().String()
	return s.repo.Data.WriteBar(ctx, bar)
}

func (s service) EditBar(ctx context.Context, bar *pb.Bar) (*pb.Bar, error) {
	const funcName = `EditBar`
	return s.repo.Data.ModifyBar(ctx, bar)
}

func (s service) DeleteBar(ctx context.Context, selects *pb.Select) (*pb.Bar, error) {
	const funcName = `DeleteBar`
	return s.repo.Data.RemoveBar(ctx, selects)
}

func (s service) GetDetailBar(ctx context.Context, selects *pb.Select) (*pb.Bar, error) {
	const funcName = `GetDetailBar`
	return s.repo.Data.ReadDetailBar(ctx, selects)
}

func (s service) GetAllBar(ctx context.Context, pagination *pb.Pagination) (*pb.Bars, error) {
	const funcName = `GetAllBar`
	console := util.ConsoleLog(gvar.Logger)
	console.Log("test", "helo")
	return s.repo.Data.ReadAllBar(ctx, pagination)
}

func NewService(repo _repo.Repository) _interface.BarService {
	return &service{
		repo: repo,
	}
}
