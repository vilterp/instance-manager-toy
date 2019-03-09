package main

import (
	"log"
	"net"

	"github.com/cockroachlabs/instance_manager/pure_manager/proto"
	"github.com/cockroachlabs/instance_manager/pure_manager/server"
	"github.com/grpc-ecosystem/go-grpc-middleware"
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
