package server

import (
	"common/api/order"
	"order/internal/service"

	"google.golang.org/grpc"
)

type OrderServer struct {
	order.UnimplementedOrderServiceServer
	GrpcServer *grpc.Server
	service    service.OrderService
}

func NewOrderServer(opts grpc.ServerOption) *OrderServer {
	srv := &OrderServer{
		service: *service.NewOrderService(),
	}

	srv.GrpcServer = grpc.NewServer(opts)
	order.RegisterOrderServiceServer(srv.GrpcServer, srv)

	allServers = append(allServers, srv)
	return srv
}

var allServers []*OrderServer

func Deinit() {
	for _, srv := range allServers {
		srv.GrpcServer.GracefulStop()
	}
}