package main

import (
	"auth/internal/middleware"
	"auth/internal/router"
	"auth/internal/server"
	"auth/internal/setup"
	api "base/api/go"
	"base/pkg/sugar"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

var dbPool *pgxpool.Pool

func main() {
	err := setup.InitAll("../.env", sugar.Default(os.Getwd()))
	if err != nil {
		panic(err)
	}

	gRPCServer := grpc.NewServer() // TODO: что с этим можно сделать; TODO: что это такое?
	api.RegisterAuthServiceServer(gRPCServer, &server.Server{}) // TODO: Зачем еще один уровень абстракций?

	mux := http.NewServeMux()
	router.SetupRoutes(mux)
	wrapped := middleware.LoggingMiddleware(mux)

	log.Printf("Server is running on port %d\n", setup.Env.Port)
	err = http.ListenAndServe(":"+strconv.Itoa(setup.Env.Port), wrapped)
	if err != nil {
		panic(err)
	}
}