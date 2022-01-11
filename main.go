package main

import (
	"context"
	"crypto/tls"
	"fmt"
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
	"github.com/nightsilvertech/utl/preparation"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
)

type Secure struct {
	ServerCertPath     string
	ServerKeyPath      string
	ServerNameOverride string
}

func (secure Secure) ServeGRPC(service pb.BarServiceServer) error {
	level.Info(gvar.Logger).Log(console.LogInfo, "serving grpc server")
	address := fmt.Sprintf("%s:%s", constant.Host, constant.GrpcPort)
	serverCert, err := tls.LoadX509KeyPair(secure.ServerCertPath, secure.ServerKeyPath)
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

func (secure Secure) ServeHTTP(service pb.BarServiceServer) error {
	level.Info(gvar.Logger).Log(console.LogInfo, "serving http server")
	httpAddress := fmt.Sprintf("%s:%s", constant.Host, constant.HttpPort)
	grpcAddress := fmt.Sprintf("%s:%s", constant.Host, constant.GrpcPort)
	clientCert, err := credentials.NewClientTLSFromFile(secure.ServerCertPath, secure.ServerNameOverride)
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
	return http.ListenAndServeTLS(httpAddress, secure.ServerCertPath, secure.ServerKeyPath, mux)
}

func Serve(service pb.BarServiceServer) {
	secure := Secure{
		ServerCertPath:     "C:\\Users\\Asus\\Desktop\\tls\\bar\\server.crt",
		ServerKeyPath:      "C:\\Users\\Asus\\Desktop\\tls\\bar\\server.key",
		ServerNameOverride: "0.0.0.0",
	}

	g := new(errgroup.Group)
	g.Go(func() error { return secure.ServeGRPC(service) })
	g.Go(func() error { return secure.ServeHTTP(service) })
	log.Fatal(g.Wait())
}

func main() {
	prepare := preparation.Data{
		LoggingFilePath:            "C:\\Users\\Asus\\Desktop\\service.log",
		TracerUrl:                  "http://localhost:9411/api/v2/spans",
		CircuitBreakerTimeout:      constant.CircuitBreakerTimout,
		ServiceName:                constant.ServiceName,
		ZipkinEndpointPort:         constant.ZipkinHostPort,
		Debug:                      false,
		FractionProbabilitySampler: 1,
	}
	prepare.CircuitBreaker()
	gvar.Logger = prepare.Logger()
	tracer := prepare.Tracer()

	repositories := repository.NewRepository(tracer)
	services := service.NewService(*repositories, tracer)
	endpoints := ep.NewBarEndpoint(services)
	server := transport.NewBarServer(endpoints)
	Serve(server)
}
