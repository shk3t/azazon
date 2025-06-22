package router

import (
	"auth/internal/handler"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", handler.Register)
	mux.HandleFunc("POST /auth/login", handler.Login)
}

func SetupMiddlewares(mux *http.ServeMux) {
	// middleware.LoggingMiddleware(mux)  // TODO
}