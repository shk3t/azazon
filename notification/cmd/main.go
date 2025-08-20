package main

import (
	"base/pkg/helper"
	"base/pkg/interceptor"
	"base/pkg/log"
	"base/pkg/sugar"
	"fmt"
	"net"
	"notification/internal/config"
	"notification/internal/server"
	"notification/internal/setup"
	"path/filepath"

	"google.golang.org/grpc"
)

func main() {
	workDir, _ := helper.GetwdCdBack("notification", "cmd")
	workDir = filepath.Join(workDir, "notification")
	err := setup.InitAll(workDir)
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

	srv := server.NewNotificationServer(
		grpc.ChainUnaryInterceptor(interceptor.LoggingUnaryInterceptor),
	)

	logger.Printf(
		"%s server is running on :%d\n",
		config.AppName, port,
	)
	err = srv.GrpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}