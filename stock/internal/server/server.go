package server

import (
	"base/api/stock"
	"stock/internal/service"

	"google.golang.org/grpc"
)

type StockServer struct {
	stock.UnimplementedStockServiceServer
	GrpcServer *grpc.Server
	service    service.StockService
}

func NewStockServer(opts grpc.ServerOption) *StockServer {
	srv := &StockServer{
		service: *service.NewStockService(),
	}

	srv.GrpcServer = grpc.NewServer(opts)
	stock.RegisterStockServiceServer(srv.GrpcServer, srv)

	allServers = append(allServers, srv)
	return srv
}

var allServers []*StockServer

func Deinit() {
	for _, srv := range allServers {
		srv.GrpcServer.GracefulStop()
	}
}