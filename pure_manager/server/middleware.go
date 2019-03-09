package server

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func LogUnaryInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	log.Println("handling", info.FullMethod)
	resp, err := handler(ctx, req)
	end := time.Now()
	log.Println("handled", info.FullMethod, "in", end.Sub(start))
	return resp, err
}

func recoveryHandler(p interface{}) (err error) {
	// Just convert the panic to an error, which may be caught another
	// interceptor.
	return errors.Errorf("panic: %s", p)
}

var RecoveryUnaryInterceptor = grpc_recovery.UnaryServerInterceptor(
	grpc_recovery.WithRecoveryHandler(recoveryHandler))
