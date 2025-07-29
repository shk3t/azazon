package server

import (
	"auth/internal/service"
	"base/api/auth"
	"base/pkg/model"
	"context"

	"google.golang.org/grpc"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAuthServer() *AuthServer {
	return &AuthServer{
		service: *service.NewAuthService(),
	}
}

func CreateAuthServer(opts grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts)
	auth.RegisterAuthServiceServer(srv, NewAuthServer())

	runningServers = append(runningServers, srv)

	return srv
}

func (s *AuthServer) Register(ctx context.Context, in *auth.User) (*auth.AuthResponse, error) {
	resp, err := s.service.Register(ctx, model.NewUserFromGrpc(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return resp.Grpc(), nil
}

func (s *AuthServer) Login(ctx context.Context, in *auth.User) (*auth.AuthResponse, error) {
	resp, err := s.service.Login(ctx, model.NewUserFromGrpc(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return resp.Grpc(), nil
}

var runningServers []*grpc.Server

func Deinit() {
	for _, srv := range runningServers {
		srv.GracefulStop()
	}
}