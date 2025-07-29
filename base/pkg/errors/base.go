package errors

import "errors"

var (
	NotFound = errors.New("Not found")
	InvalidToken = errors.New("Invalid token")
)