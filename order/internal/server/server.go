package server

import (
	"common/api/order"
	"order/internal/service"

	"google.golang.org/grpc"
)

type OrderServer struct {
	order.UnimplementedOrderServiceServer
	GrpcServer *grpc.Server
	service    *service.OrderService
}

func NewOrderServer(opts grpc.ServerOption) *OrderServer {
	s := &OrderServer{
		service: service.NewOrderService(),
	}

	s.GrpcServer = grpc.NewServer(opts)
	order.RegisterOrderServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

var allServers []*OrderServer

func Deinit() {
	for _, s := range allServers {
		s.GrpcServer.GracefulStop()
	}
}