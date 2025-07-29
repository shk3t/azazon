package main

import (
	"auth/internal/config"
	"auth/internal/interceptor"
	"auth/internal/server"
	"auth/internal/setup"
	"base/pkg/log"
	"base/pkg/sugar"
	"net"
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

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.Env.Port))
	if err != nil {
		panic(err)
	}

	srv := server.CreateAuthServer(
		grpc.ChainUnaryInterceptor(interceptor.LoggingUnaryInterceptor),
	)

	log.RLog("Server is running on :" + strconv.Itoa(config.Env.Port))
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}
