package server

import (
	"auth/internal/handler"
	"auth/internal/model"
	api "base/api/go"
	"context"
)

type AuthServer struct {
	api.UnimplementedAuthServiceServer
	auth api.AuthServiceServer
}

func (s *AuthServer) Register(ctx context.Context, in *api.User) (*api.AuthResponse, error) {
	resp, err := handler.Register(ctx, model.NewUserFromGrpc(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return resp.Grpc(), nil
}

func (s *AuthServer) Login(ctx context.Context, in *api.User) (*api.AuthResponse, error) {
	resp, err := handler.Login(ctx, model.NewUserFromGrpc(in))
	if err != nil {
		return nil, err.Grpc()
	}
	return resp.Grpc(), nil
}