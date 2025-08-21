package grpcutil

import (
	"common/pkg/log"
	"fmt"
	"net/http"

	"google.golang.org/grpc/status"
)

type HandlerError struct {
	HttpCode int
	Message  string
}

func NewError(code int, msg string) *HandlerError {
	return &HandlerError{HttpCode: code, Message: msg}
}

func NewInternalError(err error) *HandlerError {
	log.Loggers.Debug.Println(err)
	return &HandlerError{HttpCode: http.StatusInternalServerError, Message: ""}
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