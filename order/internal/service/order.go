package service

import (
	errpkg "common/pkg/errors"
	"common/pkg/grpcutil"
	"context"
	"errors"
	"net/http"
	"order/internal/model"
	"order/internal/store"

	"github.com/jackc/pgx/v5"
)

var NewErr = grpcutil.NewServiceError
var NewInternalErr = grpcutil.NewInternalError

type orderStore interface {
	Get(ctx context.Context, id int) (model.Order, error)
	Save(ctx context.Context, tx pgx.Tx, order model.Order) (model.Order, error)
}

type OrderService struct {
	store orderStore
}

func NewOrderService() *OrderService {
	return &OrderService{
		store: &store.PostgreOrderStore{},
	}
}

func (s *OrderService) SaveOrder(
	ctx context.Context,
	tx pgx.Tx,
	body model.Order,
) (*model.Order, *grpcutil.ServiceError) {
	order, err := s.store.Save(ctx, tx, body)
	if err != nil {
		return nil, NewInternalErr(err)
	}
	return &order, nil
}

func (s *OrderService) GetOrderInfo(
	ctx context.Context,
	orderId int,
) (*model.Order, *grpcutil.ServiceError) {
	order, err := s.store.Get(ctx, orderId)

	if err != nil {
		if errors.Is(err, errpkg.NotFound) {
			return nil, NewErr(http.StatusNotFound, "Order is not found")
		}
		return nil, NewInternalErr(err)
	}

	return &order, nil
}