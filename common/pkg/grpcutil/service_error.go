package grpcutil

import (
	"common/pkg/log"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/status"
)

type ServiceError struct {
	HttpCode int
	Message  string
}

func logWithStack(err error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		err = fmt.Errorf("%w (Message: %s, Code: %s, Detail: %s, Hint: %s, Position: %d)",
			err, pgErr.Message, pgErr.Code, pgErr.Detail, pgErr.Hint, pgErr.Position,
		)
	}

	stackTrace := string(debug.Stack())
	err = fmt.Errorf("%w\n%s", err, stackTrace)
	log.Loggers.Debug.Println(err)
}

func NewServiceError(code int, msg string) *ServiceError {
	return &ServiceError{HttpCode: code, Message: msg}
}

func NewInternalError(err error) *ServiceError {
	logWithStack(err)
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

func (err *ServiceError) Interface() error {
	if err == nil {
		return nil
	}
	return err
}

func NewGrpcError(code int, msg string) error {
	return (&ServiceError{HttpCode: code, Message: msg}).Grpc()
}

func NewInternalGrpcError(err error) error {
	return NewInternalError(err).Grpc()
}