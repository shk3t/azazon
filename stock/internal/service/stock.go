package service

import (
	errpkg "common/pkg/errors"
	"common/pkg/grpcutil"
	"context"
	"errors"
	"net/http"
	"stock/internal/model"
	"stock/internal/store"
)

var NewErr = grpcutil.NewServiceError
var NewInternalErr = grpcutil.NewInternalError

type stores struct {
	product productStore
	stock   stockStore
}

type productStore interface {
	Get(ctx context.Context, id int) (model.Product, error)
	Save(ctx context.Context, product model.Product) (model.Product, error)
	Delete(ctx context.Context, id int) error
}

type stockStore interface {
	Get(ctx context.Context, productId int) (model.Stock, error)
	Save(ctx context.Context, stock model.Stock) (model.Stock, error)
}

type StockService struct {
	stores stores
}

func NewStockService() *StockService {
	return &StockService{
		stores: stores{
			product: &store.PostgreProductStore{},
			stock:   &store.PostgreStockStore{},
		},
	}
}

func (s *StockService) SaveProduct(
	ctx context.Context,
	body model.Product,
) (*model.Product, *grpcutil.ServiceError) {
	product, err := s.stores.product.Save(ctx, body)
	if err != nil {
		return nil, NewInternalErr(err)
	}
	return &product, nil
}

func (s *StockService) IncreaseStockQuantity(
	ctx context.Context,
	productId int,
	quantityDelta int,
) (*model.Stock, *grpcutil.ServiceError) {
	stock, err := s.stores.stock.Get(ctx, productId)
	if err != nil {
		if errors.Is(err, errpkg.NotFound) {
			return nil, NewErr(http.StatusNotFound, "Product is not found")
		}
		return nil, NewInternalErr(err)
	}

	stock.Quantity += quantityDelta

	stock, err = s.stores.stock.Save(ctx, stock)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	return &stock, nil
}

func (s *StockService) GetStockInfo(
	ctx context.Context,
	productId int,
) (*model.Stock, *grpcutil.ServiceError) {
	stock, err := s.stores.stock.Get(ctx, productId)
	if err != nil {
		if errors.Is(err, errpkg.NotFound) {
			return nil, NewErr(http.StatusNotFound, "Product is not found")
		}
		return nil, NewInternalErr(err)
	}
	return &stock, nil
}

func (s *StockService) GetProductInfo(
	ctx context.Context,
	id int,
) (*model.Product, *grpcutil.ServiceError) {
	product, err := s.stores.product.Get(ctx, id)
	if err != nil {
		if errors.Is(err, errpkg.NotFound) {
			return nil, NewErr(http.StatusNotFound, "Product is not found")
		}
		return nil, NewInternalErr(err)
	}
	return &product, nil
}

func (s *StockService) DeleteProduct(
	ctx context.Context,
	id int,
) *grpcutil.ServiceError {
	err := s.stores.product.Delete(ctx, id)
	if err != nil {
		return NewInternalErr(err)
	}
	return nil
}