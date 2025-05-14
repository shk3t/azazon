package router

import (
	"auth/internal/handler"
	"auth/internal/middleware"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", handler.Register)
	mux.HandleFunc("POST /login", handler.Login)
}

func SetupMiddlewares(mux *http.ServeMux) {
	middleware.LoggingMiddleware(mux)  // TODO
}