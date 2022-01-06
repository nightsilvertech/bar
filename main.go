package main

import (
	"context"
	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	"crypto/tls"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nightsilvertech/bar/constant"
	ep "github.com/nightsilvertech/bar/endpoint"
	"github.com/nightsilvertech/bar/gvar"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	"github.com/nightsilvertech/bar/repository"
	"github.com/nightsilvertech/bar/service"
	"github.com/nightsilvertech/bar/transport"
	"github.com/nightsilvertech/utl/console"
	"github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
)

const address = `localhost:1900`
const grpcAddress = `localhost:9081`
const httpAddress = `localhost:8081`

func ServeGRPC(service pb.BarServiceServer, serverOptions []grpc.ServerOption) error {
	level.Info(gvar.Logger).Log(console.LogInfo, "serving grpc server")

	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(serverOptions...)
	pb.RegisterBarServiceServer(grpcServer, service)
	return grpcServer.Serve(listener)
}

func ServeHTTP(service pb.BarServiceServer, clientOption []grpc.DialOption) error {
	level.Info(gvar.Logger).Log(console.LogInfo, "serving http server")

	mux := runtime.NewServeMux()
	pb.RegisterBarServiceServer(grpc.NewServer(), service)
	err := pb.RegisterBarServiceHandlerFromEndpoint(context.Background(), mux, address, clientOption)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8081", mux)
}

func Serve(service pb.BarServiceServer, useTls bool) {
	serverCert, err := tls.LoadX509KeyPair(
		"C:\\Users\\Asus\\Desktop\\tls\\server.crt",
		"C:\\Users\\Asus\\Desktop\\tls\\server.key")
	if err != nil {
		panic(err)
	}

	clientCert, err := credentials.NewClientTLSFromFile(
		"C:\\Users\\Asus\\Desktop\\tls\\server.crt",
		"localhost",
	)
	if err != nil {
		panic(err)
	}

	var serverOpts []grpc.ServerOption
	var clientOpts []grpc.DialOption

	if useTls {
		serverOpts = []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&serverCert))}
		clientOpts = []grpc.DialOption{grpc.WithTransportCredentials(clientCert)}
	} else {
		serverOpts = []grpc.ServerOption{}
		clientOpts = []grpc.DialOption{grpc.WithInsecure()}
	}

	g := new(errgroup.Group)
	g.Go(func() error { return ServeGRPC(service, serverOpts) })
	g.Go(func() error { return ServeHTTP(service, clientOpts) })
	log.Fatal(g.Wait())
}

func main() {
	gvar.Logger = console.CreateStdGoKitLog(constant.ServiceName, false, "C:\\Users\\Asus\\Desktop\\service.log")

	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	localEndpoint, _ := zipkin.NewEndpoint(constant.ServiceName, ":0")
	exporter := oczipkin.NewExporter(reporter, localEndpoint)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.RegisterExporter(exporter)
	tracer := trace.DefaultTracer
	hystrix.ConfigureCommand(constant.ServiceName, hystrix.CommandConfig{Timeout: 1000 * 30})

	repositories, err := repository.NewRepository(tracer)
	if err != nil {
		panic(err)
	}
	services := service.NewService(*repositories, tracer)
	endpoints := ep.NewBarEndpoint(services)
	server := transport.NewBarServer(endpoints)

	Serve(server, true)
}
