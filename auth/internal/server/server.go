package server

import (
	conv "auth/internal/conversion"
	"auth/internal/service"
	"common/api/auth"
	"common/pkg/sugar"
	"context"

	"google.golang.org/grpc"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	GrpcServer *grpc.Server
	service    *service.AuthService
}

func NewAuthServer(opts grpc.ServerOption) *AuthServer {
	s := &AuthServer{
		service: service.NewAuthService(),
	}

	s.GrpcServer = grpc.NewServer(opts)
	auth.RegisterAuthServiceServer(s.GrpcServer, s)

	allServers = append(allServers, s)
	return s
}

func (s *AuthServer) Register(
	ctx context.Context,
	in *auth.RegisterRequest,
) (*auth.RegisterResponse, error) {
	resp, err := s.service.Register(ctx, *conv.User(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return conv.RegisterResponse(resp), nil
}

func (s *AuthServer) Login(
	ctx context.Context,
	in *auth.LoginRequest,
) (*auth.LoginResponse, error) {
	resp, err := s.service.Login(ctx, *conv.User(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return conv.LoginResponse(resp), nil
}

func (s *AuthServer) ValidateToken(
	ctx context.Context,
	in *auth.ValidateTokenRequest,
) (*auth.ValidateTokenResponse, error) {
	err := s.service.ValidateToken(ctx, in.Token)
	if err != nil {
		return nil, err.Grpc()
	}
	return &auth.ValidateTokenResponse{Valid: true}, nil
}

func (s *AuthServer) UpdateUser(
	ctx context.Context,
	in *auth.UpdateUserRequest,
) (*auth.UpdateUserResponse, error) {
	resp, err := s.service.UpdateUser(
		ctx,
		in.Token,
		*conv.User(in),
		sugar.Value(in.RoleKey),
	)
	if err != nil {
		return nil, err.Grpc()
	}
	return conv.UpdateUserResponse(resp), nil
}

var allServers []*AuthServer

func Deinit() {
	for _, s := range allServers {
		s.GrpcServer.GracefulStop()
	}
}