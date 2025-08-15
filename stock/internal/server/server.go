package server

import (
	"stock/internal/service"
	"base/api/stock"

	"google.golang.org/grpc"
)

type StockServer struct {
	stock.UnimplementedStockServiceServer
	service service.StockService
}

func NewStockServer() *StockServer {
	return &StockServer{
		service: *service.NewStockService(),
	}
}

func CreateStockServer(opts grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts)
	stock.RegisterStockServiceServer(srv, NewStockServer())

	runningServers = append(runningServers, srv)

	return srv
}

var runningServers []*grpc.Server

func Deinit() {
	for _, srv := range runningServers {
		srv.GracefulStop()
	}
}