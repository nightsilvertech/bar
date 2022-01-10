package main

import (
	"context"
	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	"crypto/tls"
	"fmt"
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

type TLSPrepare struct {
	ServerCertPath     string
	ServerKeyPath      string
	ServerNameOverride string
}

func (tlsp TLSPrepare) ServeGRPC(service pb.BarServiceServer) error {
	level.Info(gvar.Logger).Log(console.LogInfo, "serving grpc server")
	address := fmt.Sprintf("%s:%s", constant.Host, constant.GrpcPort)
	serverCert, err := tls.LoadX509KeyPair(tlsp.ServerCertPath, tlsp.ServerKeyPath)
	if err != nil {
		return err
	}
	serverOpts := []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&serverCert))}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterBarServiceServer(grpcServer, service)
	return grpcServer.Serve(listener)
}

func (tlsp TLSPrepare) ServeHTTP(service pb.BarServiceServer) error {
	level.Info(gvar.Logger).Log(console.LogInfo, "serving http server")
	httpAddress := fmt.Sprintf("%s:%s", constant.Host, constant.HttpPort)
	grpcAddress := fmt.Sprintf("%s:%s", constant.Host, constant.GrpcPort)
	clientCert, err := credentials.NewClientTLSFromFile(tlsp.ServerCertPath, tlsp.ServerNameOverride)
	if err != nil {
		return err
	}
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(clientCert)}

	mux := runtime.NewServeMux()
	pb.RegisterBarServiceServer(grpc.NewServer(), service)
	err = pb.RegisterBarServiceHandlerFromEndpoint(context.Background(), mux, grpcAddress, dialOptions)
	if err != nil {
		return err
	}
	return http.ListenAndServeTLS(httpAddress, tlsp.ServerCertPath, tlsp.ServerKeyPath, mux)
}

func Serve(service pb.BarServiceServer) {
	tlsp := TLSPrepare{
		ServerCertPath:     "C:\\Users\\Asus\\Desktop\\tls\\server.crt",
		ServerKeyPath:      "C:\\Users\\Asus\\Desktop\\tls\\server.key",
		ServerNameOverride: "0.0.0.0",
	}

	g := new(errgroup.Group)
	g.Go(func() error { return tlsp.ServeGRPC(service) })
	g.Go(func() error { return tlsp.ServeHTTP(service) })
	log.Fatal(g.Wait())
}

func main() {
	gvar.Logger = console.CreateStdGoKitLog(constant.ServiceName, false, "C:\\Users\\Asus\\Desktop\\service.log")

	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	localEndpoint, _ := zipkin.NewEndpoint(constant.ServiceName, constant.ZipkinHostPort)
	exporter := oczipkin.NewExporter(reporter, localEndpoint)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.RegisterExporter(exporter)
	tracer := trace.DefaultTracer
	hystrix.ConfigureCommand(constant.ServiceName, hystrix.CommandConfig{Timeout: constant.CircuitBreakerTimout})

	repositories, err := repository.NewRepository(tracer)
	if err != nil {
		panic(err)
	}
	services := service.NewService(*repositories, tracer)
	endpoints := ep.NewBarEndpoint(services)
	server := transport.NewBarServer(endpoints)
	Serve(server)
}
