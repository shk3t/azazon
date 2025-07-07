package server

import (
	api "base/api/go"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	api.UnimplementedAuthServiceServer
	auth api.AuthServiceServer
}

func (s *Server) Register(context.Context, *api.User) (*api.AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func (s *Server) Login(context.Context, *api.User) (*api.AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}