package grpcutil

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type HandlerError struct {
	HttpCode int
	Message  string
}

func NewError(code int, msg string) *HandlerError {
	return &HandlerError{HttpCode: code, Message: msg}
}

func (err HandlerError) Error() string {
	return fmt.Sprintf("%d | %s", err.HttpCode, err.Message)
}

func (err HandlerError) Grpc() error {
	return status.Error(
		HttpToGrpcStatus(err.HttpCode),
		err.Message,
	)
}