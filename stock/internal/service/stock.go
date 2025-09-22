package service

import (
	"common/pkg/grpcutil"
	"context"
	"stock/internal/model"
)

var NewErr = grpcutil.NewServiceError
var NewInternalErr = grpcutil.NewInternalError

type stockStore interface {
}

type StockService struct {
	store stockStore
}

func NewStockService() *StockService {
	return &StockService{
		store: nil,
	}
}

func (s *StockService) SaveProduct(
	ctx context.Context,
	body model.Product,
) (*model.Product, *grpcutil.ServiceError) {
	return nil, nil
}

func (s *StockService) IncreaseStockQuantity(
	ctx context.Context,
	productId int,
	quantityDelta int,
) (*model.Stock, *grpcutil.ServiceError) {
	return nil, nil
}

func (s *StockService) GetStockInfo(
	ctx context.Context,
	id int,
) (*model.Stock, *grpcutil.ServiceError) {
	return nil, nil
}

func (s *StockService) GetProductInfo(
	ctx context.Context,
	productId int,
) (*model.Product, *grpcutil.ServiceError) {
	return nil, nil
}

func (s *StockService) DeleteProduct(
	ctx context.Context,
	productId int,
) *grpcutil.ServiceError {
	return nil
}