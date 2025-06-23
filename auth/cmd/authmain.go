package main

import (
	"auth/internal/config"
	"auth/internal/database"
	"auth/internal/router"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	ctx := context.Background()

	if err := config.LoadEnvs("../.env"); err != nil {
		panic(err)
	}

	database.Connect(ctx)
	defer database.ConnPool.Close()

	mux := http.NewServeMux()
	router.SetupRoutes(mux)
	router.SetupMiddlewares(mux)

	log.Printf("Server is running on port %d\n", config.Env.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(config.Env.Port), mux)
	if err != nil {
		panic(err)
	}
}