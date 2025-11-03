package interceptor

import (
	"common/pkg/log"
	"context"
	"os"

	"google.golang.org/grpc"
)

var hostname, _ = os.Hostname()

func LoggingUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Loggers.Event.Printf(
		"Hostname: %s | Unary RPC: %s, request: %v",
		hostname, info.FullMethod, req,
	)

	resp, err := handler(ctx, req)

	log.Loggers.Event.Printf(
		"Hostname: %s | Unary RPC: %s, response: %v, error: %v",
		hostname, info.FullMethod, resp, err,
	)

	return resp, err
}