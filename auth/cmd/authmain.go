package main

import (
	"auth/internal/middleware"
	"auth/internal/router"
	"auth/internal/setup"
	"base/pkg/sugar"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	err := setup.InitAll("../.env", sugar.Default(os.Getwd()))
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	router.SetupRoutes(mux)
	wrapped := middleware.LoggingMiddleware(mux)

	log.Printf("Server is running on port %d\n", setup.Env.Port)
	err = http.ListenAndServe(":"+strconv.Itoa(setup.Env.Port), wrapped)
	if err != nil {
		panic(err)
	}
}