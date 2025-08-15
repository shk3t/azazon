package interceptor

import (
	"base/pkg/log"
	"context"

	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Loggers.Request.Printf(
		"Unary RPC: %s, request: %v",
		info.FullMethod, req,
	)

	resp, err := handler(ctx, req)

	log.Loggers.Request.Printf(
		"Unary RPC: %s, response: %v, error: %v",
		info.FullMethod, resp, err,
	)

	return resp, err
}