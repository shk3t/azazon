package service

import (
	errpkg "common/pkg/errors"
	"common/pkg/grpcutil"
	"common/pkg/sugar"
	"context"
	"errors"
	"fmt"
	"net/http"
	db "stock/internal/database"
	"stock/internal/model"
	"stock/internal/store"
	"time"

	"github.com/jackc/pgx/v5"
)

var NewErr = grpcutil.NewServiceError
var NewInternalErr = grpcutil.NewInternalError

type stores struct {
	product productStore
	stock   stockStore
	reserve reserveStore
}

type productStore interface {
	Get(ctx context.Context, id int) (model.Product, error)
	Save(ctx context.Context, product model.Product) (model.Product, error)
	Delete(ctx context.Context, id int) error
}

type stockStore interface {
	Get(ctx context.Context, productId int) (model.Stock, error)
	Save(ctx context.Context, tx pgx.Tx, stock model.Stock) (model.Stock, error)
}

type reserveStore interface {
	Get(ctx context.Context, orderId int, userId int) (model.Reserve, error)
	GetOlder(ctx context.Context, olderThan time.Time) ([]model.Reserve, error)
	Create(ctx context.Context, tx pgx.Tx, reserve model.Reserve) error
	Delete(ctx context.Context, tx pgx.Tx, reserve model.Reserve) error
}

type StockService struct {
	stores stores
}

func NewStockService() *StockService {
	return &StockService{
		stores: stores{
			product: &store.PostgreProductStore{},
			stock:   &store.PostgreStockStore{},
			reserve: &store.PostgreReserveStore{},
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

func (s *StockService) ChangeStockQuantity(
	ctx context.Context,
	tx pgx.Tx,
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
	if stock.Quantity < 0 {
		return nil, NewErr(
			http.StatusBadRequest,
			fmt.Sprintf("Product_%d: not enough items in stock", productId),
		)
	}

	stock, err = s.stores.stock.Save(ctx, tx, stock)
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

func (s *StockService) Reserve(
	ctx context.Context,
	body model.Reserve,
	undo bool,
) (*model.Stock, *grpcutil.ServiceError) {
	tx, _ := db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if undo {
		err := s.stores.reserve.Delete(ctx, tx, body)
		if err != nil {
			return nil, NewInternalErr(err)
		}
	} else {
		err := s.stores.reserve.Create(ctx, tx, body)
		if err != nil {
			return nil, NewInternalErr(err)
		}
	}

	sign := sugar.If(undo, 1, -1)
	stock, sErr := s.ChangeStockQuantity(ctx, tx, body.ProductId, sign*body.Quantity)
	if sErr != nil {
		return nil, sErr
	}


	if err := tx.Commit(ctx); err != nil {
		return nil, NewInternalErr(err)
	}
	return stock, nil
}

func (s *StockService) UndoOldReserves(
	ctx context.Context,
	olderThan time.Time,
) *grpcutil.ServiceError {
	reserves, err := s.stores.reserve.GetOlder(ctx, olderThan)
	if err != nil {
		return NewInternalErr(err)
	}

	for _, reserve := range reserves {
		_ = s.stores.reserve.Delete(ctx, nil, reserve)
	}

	return nil
}