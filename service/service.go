package service

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"github.com/nightsilvertech/bar/gvar"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_repo "github.com/nightsilvertech/bar/repository"
	_interface "github.com/nightsilvertech/bar/service/interface"
	"github.com/nightsilvertech/utl/console"
	uuid "github.com/satori/go.uuid"
	"go.opencensus.io/trace"
)

type service struct {
	tracer trace.Tracer
	repo   _repo.Repository
}

func (s *service) AddBar(ctx context.Context, bar *pb.Bar) (*pb.Bar, error) {
	const funcName = `AddBar`
	ctx, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	span.SetStatus(trace.Status{Code: int32(trace.StatusCodeNotFound), Message: "Cache miss"})

	ctx, consoleLog := console.Log(ctx, gvar.Logger, funcName)

	level.Info(consoleLog).Log(console.LogInfo, "upper", console.LogData, bar)

	bar.Id = uuid.NewV4().String()

	level.Info(consoleLog).Log(console.LogInfo, "downer")

	return s.repo.Data.WriteBar(ctx, bar)
}

func (s *service) EditBar(ctx context.Context, bar *pb.Bar) (*pb.Bar, error) {
	const funcName = `EditBar`
	return s.repo.Data.ModifyBar(ctx, bar)
}

func (s *service) DeleteBar(ctx context.Context, selects *pb.Select) (*pb.Bar, error) {
	const funcName = `DeleteBar`
	return s.repo.Data.RemoveBar(ctx, selects)
}

func (s *service) GetDetailBar(ctx context.Context, selects *pb.Select) (*pb.Bar, error) {
	const funcName = `GetDetailBar`
	return s.repo.Data.ReadDetailBar(ctx, selects)
}

func (s *service) GetAllBar(ctx context.Context, pagination *pb.Pagination) (*pb.Bars, error) {
	const funcName = `GetAllBar`
	return s.repo.Data.ReadAllBar(ctx, pagination)
}

func NewService(repo _repo.Repository, tracer trace.Tracer) _interface.BarService {
	return &service{
		tracer: tracer,
		repo:   repo,
	}
}
