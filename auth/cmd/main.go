package main

import (
	"auth/internal/config"
	"auth/internal/server"
	"auth/internal/setup"
	"base/pkg/interceptor"
	"base/pkg/log"
	"base/pkg/sugar"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	err := setup.InitAll(sugar.Default(os.Getwd()))
	if err != nil {
		panic(err)
	}

	logger := log.Loggers.Run
	env := config.Env
	port := sugar.If(env.Test, env.TestPort, env.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := server.CreateAuthServer(
		grpc.ChainUnaryInterceptor(interceptor.LoggingUnaryInterceptor),
	)

	logger.Printf(
		"%s server is running on :%d\n",
		config.AppName, port,
	)
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}