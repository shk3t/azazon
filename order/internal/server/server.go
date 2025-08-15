package server

import (
	"order/internal/service"
	"base/api/order"

	"google.golang.org/grpc"
)

type OrderServer struct {
	order.UnimplementedOrderServiceServer
	service service.OrderService
}

func NewOrderServer() *OrderServer {
	return &OrderServer{
		service: *service.NewOrderService(),
	}
}

func CreateOrderServer(opts grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts)
	order.RegisterOrderServiceServer(srv, NewOrderServer())

	runningServers = append(runningServers, srv)

	return srv
}

var runningServers []*grpc.Server

func Deinit() {
	for _, srv := range runningServers {
		srv.GracefulStop()
	}
}