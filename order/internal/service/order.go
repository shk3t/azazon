package service

import "common/pkg/grpcutil"

var NewErr = grpcutil.NewError
var NewInternalErr = grpcutil.NewInternalError

type orderStore interface {
}

type OrderService struct {
	store orderStore
}

func NewOrderService() *OrderService {
	return &OrderService{
		store: nil,
	}
}