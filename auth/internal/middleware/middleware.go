package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func SetupMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}