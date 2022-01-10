package service

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"github.com/nightsilvertech/bar/gvar"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_repo "github.com/nightsilvertech/bar/repository"
	_interface "github.com/nightsilvertech/bar/service/interface"
	"github.com/nightsilvertech/utl/console"
	"go.opencensus.io/trace"
)

type service struct {
	tracer trace.Tracer
	repo   _repo.Repository
}

func (s *service) AddBar(ctx context.Context, bar *pb.Bar) (res *pb.Bar, err error) {
	const funcName = `AddBar`
	ctx, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	// console log initialization
	ctx, consoleLog := console.Log(ctx, gvar.Logger, funcName)

	// upper log info
	level.Info(consoleLog).Log(console.LogInfo, "upper", console.LogData, bar)

	// logics
	res, err = s.repo.Data.WriteBar(ctx, bar)
	if err != nil {
		// error log
		level.Error(consoleLog).Log(console.LogErr, err)
		// span set status when error
		span.SetStatus(trace.Status{Code: int32(trace.StatusCodeInternal), Message: err.Error()})
		return res, err
	}

	// downer log info
	level.Info(consoleLog).Log(console.LogInfo, "downer")

	return res, nil
}

func (s *service) EditBar(ctx context.Context, bar *pb.Bar) (res *pb.Bar, err error) {
	const funcName = `EditBar`
	ctx, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	// console log initialization
	ctx, consoleLog := console.Log(ctx, gvar.Logger, funcName)

	// upper log info
	level.Info(consoleLog).Log(console.LogInfo, "upper", console.LogData, bar)

	// logics
	res, err = s.repo.Data.ModifyBar(ctx, bar)
	if err != nil {
		// error log
		level.Error(consoleLog).Log(console.LogErr, err)
		// span set status when error
		span.SetStatus(trace.Status{Code: int32(trace.StatusCodeInternal), Message: err.Error()})
		return res, err
	}

	// downer log info
	level.Info(consoleLog).Log(console.LogInfo, "downer")

	return res, nil
}

func (s *service) DeleteBar(ctx context.Context, selects *pb.Select) (res *pb.Bar, err error) {
	const funcName = `DeleteBar`
	ctx, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	// console log initialization
	ctx, consoleLog := console.Log(ctx, gvar.Logger, funcName)

	// upper log info
	level.Info(consoleLog).Log(console.LogInfo, "upper", console.LogData, selects)

	// logics
	res, err = s.repo.Data.RemoveBar(ctx, selects)
	if err != nil {
		// error log
		level.Error(consoleLog).Log(console.LogErr, err)
		// span set status when error
		span.SetStatus(trace.Status{Code: int32(trace.StatusCodeInternal), Message: err.Error()})
		return res, err
	}

	// downer log info
	level.Info(consoleLog).Log(console.LogInfo, "downer")

	return res, nil
}

func (s *service) GetDetailBar(ctx context.Context, selects *pb.Select) (res *pb.Bar, err error) {
	const funcName = `GetDetailBar`
	ctx, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	// console log initialization
	ctx, consoleLog := console.Log(ctx, gvar.Logger, funcName)

	// upper log info
	level.Info(consoleLog).Log(console.LogInfo, "upper", console.LogData, selects)

	// logics
	res, err = s.repo.Data.ReadDetailBar(ctx, selects)
	if err != nil {
		// error log
		level.Error(consoleLog).Log(console.LogErr, err)
		// span set status when error
		span.SetStatus(trace.Status{Code: int32(trace.StatusCodeInternal), Message: err.Error()})
		return res, err
	}

	// downer log info
	level.Info(consoleLog).Log(console.LogInfo, "downer")

	return res, nil
}

func (s *service) GetAllBar(ctx context.Context, pagination *pb.Pagination) (res *pb.Bars, err error) {
	const funcName = `GetAllBar`
	ctx, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	// console log initialization
	ctx, consoleLog := console.Log(ctx, gvar.Logger, funcName)

	// upper log info
	level.Info(consoleLog).Log(console.LogInfo, "upper", console.LogData, pagination)

	// logics
	res, err = s.repo.Data.ReadAllBar(ctx, pagination)
	if err != nil {
		// error log
		level.Error(consoleLog).Log(console.LogErr, err)
		// span set status when error
		span.SetStatus(trace.Status{Code: int32(trace.StatusCodeInternal), Message: err.Error()})
		return res, err
	}

	// downer log info
	level.Info(consoleLog).Log(console.LogInfo, "downer")

	return res, nil
}

func NewService(repo _repo.Repository, tracer trace.Tracer) _interface.BarService {
	return &service{
		tracer: tracer,
		repo:   repo,
	}
}
