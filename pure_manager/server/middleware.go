package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

var LogUnaryInterceptor = func(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	log.Println("handling", info.FullMethod)
	resp, err := handler(ctx, req)
	end := time.Now()
	log.Println("handled", info.FullMethod, "in", end.Sub(start))
	return resp, err
}
