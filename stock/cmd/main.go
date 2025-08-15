package main

import (
	"base/pkg/interceptor"
	"base/pkg/log"
	"base/pkg/sugar"
	"net"
	"os"
	"stock/internal/config"
	"stock/internal/server"
	"stock/internal/setup"
	"strconv"

	"google.golang.org/grpc"
)

func main() {
	err := setup.InitAll(sugar.Default(os.Getwd()))
	if err != nil {
		panic(err)
	}

	logger := log.Loggers.Run

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.Env.Port))
	if err != nil {
		panic(err)
	}

	srv := server.CreateStockServer(
		grpc.ChainUnaryInterceptor(interceptor.LoggingUnaryInterceptor),
	)

	logger.Printf("Server is running on :%d\n", config.Env.Port)
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}