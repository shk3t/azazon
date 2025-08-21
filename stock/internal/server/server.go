package server

import (
	"common/api/stock"
	"stock/internal/service"

	"google.golang.org/grpc"
)

type StockServer struct {
	stock.UnimplementedStockServiceServer
	GrpcServer *grpc.Server
	service    *service.StockService
}

func NewStockServer(opts grpc.ServerOption) *StockServer {
	s := &StockServer{
		service: service.NewStockService(),
	}

	s.GrpcServer = grpc.NewServer(opts)
	stock.RegisterStockServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

var allServers []*StockServer

func Deinit() {
	for _, s := range allServers {
		s.GrpcServer.GracefulStop()
	}
}