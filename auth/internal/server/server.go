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

func (s *AuthServer) Register(
	ctx context.Context,
	in *auth.RegisterRequest,
) (*auth.RegisterResponse, error) {
	resp, err := s.service.Register(ctx, *model.UserFromRegisterRequest(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return resp.RegisterResponse(), nil
}

func (s *AuthServer) Login(
	ctx context.Context,
	in *auth.LoginRequest,
) (*auth.LoginResponse, error) {
	resp, err := s.service.Login(ctx, *model.UserFromLoginRequest(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return resp.LoginResponse(), nil
}

func (s *AuthServer) ValidateToken(
	ctx context.Context,
	in *auth.ValidateTokenRequest,
) (*auth.ValidateTokenResponse, error) {
	resp := s.service.ValidateToken(ctx, in.Token)
	return &auth.ValidateTokenResponse{Valid: resp}, nil
}

var runningServers []*grpc.Server

func Deinit() {
	for _, srv := range runningServers {
		srv.GracefulStop()
	}
}