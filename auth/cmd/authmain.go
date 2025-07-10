package main

import (
	"auth/internal/server"
	"auth/internal/setup"
	api "base/api/go"
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

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(setup.Env.Port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	api.RegisterAuthServiceServer(srv, &server.AuthServer{})

	log.RLog("Server is running on :" + strconv.Itoa(setup.Env.Port))
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}