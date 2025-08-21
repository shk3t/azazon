package server

import (
	"base/api/payment"
	"payment/internal/service"

	"google.golang.org/grpc"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
	GrpcServer *grpc.Server
	service    service.PaymentService
}

func NewPaymentServer(opts grpc.ServerOption) *PaymentServer {
	srv := &PaymentServer{
		service: *service.NewPaymentService(),
	}

	srv.GrpcServer = grpc.NewServer(opts)
	payment.RegisterPaymentServiceServer(srv.GrpcServer, srv)

	allServers = append(allServers, srv)
	return srv
}

var allServers []*PaymentServer

func Deinit() {
	for _, srv := range allServers {
		srv.GrpcServer.GracefulStop()
	}
}