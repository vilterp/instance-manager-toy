package main

import (
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/vilterp/instance-manager-toy/proto"
	"github.com/vilterp/instance-manager-toy/server"
	"google.golang.org/grpc"
)

func main() {
	s := server.NewServer()
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				server.LogUnaryInterceptor,
				server.RecoveryUnaryInterceptor,
			),
		),
	)
	proto.RegisterGroupManagerServer(grpcServer, s)
	addr := "0.0.0.0:8888"
	log.Printf("listening at %s", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen", err)
	}
	log.Fatal(grpcServer.Serve(lis))
}
