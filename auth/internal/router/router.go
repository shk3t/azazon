package router

import (
	"net/http"
	"auth/internal/handler"
)

func SetupRoutes() {
	http.HandleFunc("POST /register", handler.Register)
	http.HandleFunc("POST /login", handler.Login)
}