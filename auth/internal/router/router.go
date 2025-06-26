package router

import (
	"auth/internal/handler"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /auth/register", handler.Register)
	mux.HandleFunc("POST /auth/login", handler.Login)

	return mux
	// return middleware.SetupMiddlewares(mux, middleware.LoggingMiddleware).(*http.ServeMux)
}