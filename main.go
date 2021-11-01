package main

import (
	"context"
	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nightsilvertech/bar/constant"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"

	ep "github.com/nightsilvertech/bar/endpoint"
	"github.com/nightsilvertech/bar/gvar"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	"github.com/nightsilvertech/bar/repository"
	"github.com/nightsilvertech/bar/service"
	"github.com/nightsilvertech/bar/transport"
	"github.com/nightsilvertech/utl/console"
	"github.com/openzipkin/zipkin-go"
	"github.com/soheilhy/cmux"
	"go.opencensus.io/trace"
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
	gvar.Logger = console.CreateStdGoKitLog(constant.ServiceName, false)

	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	localEndpoint, _ := zipkin.NewEndpoint(constant.ServiceName, "http://localhost:0")
	exporter := oczipkin.NewExporter(reporter, localEndpoint)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.RegisterExporter(exporter)

	repositories, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}

	services := service.NewService(*repositories)
	endpoints := ep.NewBarEndpoint(services)
	server := transport.NewBarServer(endpoints)
	MergeServer(server, nil)
}
