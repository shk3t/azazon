package service

import (
	"common/pkg/grpcutil"
	"context"
	"order/internal/model"
)

var NewErr = grpcutil.NewServiceError
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

func (s *OrderService) CreateOrder(
	ctx context.Context,
	body model.Order,
) (orderId int, err *grpcutil.ServiceError) {
	return 0, nil
}

func (s *OrderService) GetOrderInfo(
	ctx context.Context,
	orderId int,
) (*model.Order, *grpcutil.ServiceError) {
	return nil, nil
}