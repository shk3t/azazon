package server

import (
	"payment/internal/service"
	"base/api/payment"

	"google.golang.org/grpc"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{
		service: *service.NewPaymentService(),
	}
}

func CreatePaymentServer(opts grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts)
	payment.RegisterPaymentServiceServer(srv, NewPaymentServer())

	runningServers = append(runningServers, srv)

	return srv
}

var runningServers []*grpc.Server

func Deinit() {
	for _, srv := range runningServers {
		srv.GracefulStop()
	}
}