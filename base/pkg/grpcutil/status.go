package grpcutil

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func HttpToGrpcStatus(code int) codes.Code {
	switch code {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusInternalServerError:
		return codes.Internal
	default:
		return codes.Unknown
	}
}