package server

import (
	"notification/internal/service"
	"base/api/notification"

	"google.golang.org/grpc"
)

type NotificationServer struct {
	notification.UnimplementedNotificationServiceServer
	service service.NotificationService
}

func NewNotificationServer() *NotificationServer {
	return &NotificationServer{
		service: *service.NewNotificationService(),
	}
}

func CreateNotificationServer(opts grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts)
	notification.RegisterNotificationServiceServer(srv, NewNotificationServer())

	runningServers = append(runningServers, srv)

	return srv
}

var runningServers []*grpc.Server

func Deinit() {
	for _, srv := range runningServers {
		srv.GracefulStop()
	}
}