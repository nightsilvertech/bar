package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nightsilvertech/bar/constant"
	ep "github.com/nightsilvertech/bar/endpoint"
	"github.com/nightsilvertech/bar/gvar"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	"github.com/nightsilvertech/bar/repository"
	"github.com/nightsilvertech/bar/service"
	"github.com/nightsilvertech/bar/transport"
	"github.com/nightsilvertech/bar/util"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func ServeGRPC(listener net.Listener, service pb.BarServiceServer, serverOptions []grpc.ServerOption) error {
	var grpcServer *grpc.Server
	if len(serverOptions) > 0 {
		grpcServer = grpc.NewServer(serverOptions...)
	} else {
		grpcServer = grpc.NewServer()
	}
	pb.RegisterBarServiceServer(grpcServer, service)
	return grpcServer.Serve(listener)
}

func ServeHTTP(listener net.Listener, service pb.BarServiceServer) error {
	mux := runtime.NewServeMux()
	err := pb.RegisterBarServiceHandlerServer(context.Background(), mux, service)
	if err != nil {
		return err
	}
	return http.Serve(listener, mux)
}

func MergeServer(service pb.BarServiceServer, serverOptions []grpc.ServerOption) {
	port := fmt.Sprintf(":%s", "1900")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	m := cmux.New(listener)
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings(
		"content-type", "application/grpc",
	))
	httpListener := m.Match(cmux.HTTP1Fast())

	g := new(errgroup.Group)
	g.Go(func() error { return ServeGRPC(grpcListener, service, serverOptions) })
	g.Go(func() error { return ServeHTTP(httpListener, service) })
	g.Go(func() error { return m.Serve() })

	log.Fatal(g.Wait())
}

func main() {
	gvar.Logger = util.CreateStdGoKitLog(constant.ServiceName, false)

	repositories, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}

	services := service.NewService(*repositories)
	endpoints := ep.NewBarEndpoint(services)
	server := transport.NewBarServer(endpoints)
	MergeServer(server, nil)
}
