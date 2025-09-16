package grpcutil

import (
	"common/pkg/log"
	"fmt"
	"net/http"

	"google.golang.org/grpc/status"
)

type ServiceError struct {
	HttpCode int
	Message  string
}

func NewServiceError(code int, msg string) *ServiceError {
	return &ServiceError{HttpCode: code, Message: msg}
}

func NewInternalError(err error) *ServiceError {
	log.Loggers.Debug.Println(err)
	return &ServiceError{HttpCode: http.StatusInternalServerError, Message: ""}
}

func (err ServiceError) Error() string {
	return fmt.Sprintf("%d | %s", err.HttpCode, err.Message)
}

func (err ServiceError) Grpc() error {
	return status.Error(
		HttpToGrpcStatus(err.HttpCode),
		err.Message,
	)
}

func NewGrpcError(code int, msg string) error {
	return (&ServiceError{HttpCode: code, Message: msg}).Grpc()
}

func NewInternalGrpcError(err error) error {
	return NewInternalError(err).Grpc()
}