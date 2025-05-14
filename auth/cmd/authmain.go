package main

import (
	"auth/internal/config"
	"auth/internal/database"
	"auth/internal/router"
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	ctx := context.Background()

	config.LoadEnvs()

	database.Connect(ctx)
	defer database.ConnPool.Close()

	mux := http.NewServeMux()
	router.SetupRoutes(mux)
	router.SetupMiddlewares(mux)

	http.ListenAndServe(":"+strconv.Itoa(config.Env.Port), mux)
}