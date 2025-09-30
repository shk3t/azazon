package store

import (
	errpkg "common/pkg/errors"
	"context"
	"errors"
	"stock/internal/model"
	"stock/internal/query"

	"github.com/jackc/pgx/v5"
)

type PostgreStockStore struct {
	HashPassword func(password string) (string, error)
}

func (s *PostgreStockStore) Get(ctx context.Context, productId int) (model.Stock, error) {
	stock, err := query.GetStockByProductId(ctx, productId)

	if errors.Is(err, pgx.ErrNoRows) {
		return stock, errpkg.NotFound
	}

	return stock, err
}

func (s *PostgreStockStore) Save(ctx context.Context, tx pgx.Tx, stock model.Stock) (model.Stock, error) {
	_, err := query.GetStockByProductId(ctx, stock.ProductId)

	if errors.Is(err, pgx.ErrNoRows) {
		_, err = query.CreateStock(ctx, tx, stock)
	} else {
		err = query.UpdateStockByProductId(ctx, tx, stock)
	}

	return stock, err
}