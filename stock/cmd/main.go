package main

import (
	"base/pkg/helper"
	"base/pkg/interceptor"
	"base/pkg/log"
	"base/pkg/sugar"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"stock/internal/config"
	"stock/internal/server"
	"stock/internal/setup"

	"google.golang.org/grpc"
)

func main() {
	workDir, _ := helper.GetwdCdBack("stock", "cmd")
	workDir = filepath.Join(workDir, "stock")
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

	srv := server.CreateStockServer(
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