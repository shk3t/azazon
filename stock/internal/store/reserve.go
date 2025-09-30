package store

import (
	errpkg "common/pkg/errors"
	"context"
	"errors"
	"stock/internal/model"
	"stock/internal/query"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgreReserveStore struct {
	HashPassword func(password string) (string, error)
}

func (s *PostgreReserveStore) Get(
	ctx context.Context,
	orderId int,
	productId int,
) (model.Reserve, error) {
	reserve, err := query.GetReserveByOrderIdAndProductId(ctx, orderId, productId)

	if errors.Is(err, pgx.ErrNoRows) {
		return reserve, errpkg.NotFound
	}

	return reserve, err
}

func (s *PostgreReserveStore) GetOlder(
	ctx context.Context,
	olderThan time.Time,
) ([]model.Reserve, error) {
	reserves, err := query.GetReservesOlderThan(ctx, olderThan)

	if errors.Is(err, pgx.ErrNoRows) {
		return reserves, errpkg.NotFound
	}

	return reserves, err
}

func (s *PostgreReserveStore) Create(ctx context.Context, tx pgx.Tx, reserve model.Reserve) error {
	id, err := query.CreateReserve(ctx, tx, reserve)
	if id == 0 {
		return errpkg.Duplicate
	}
	return err
}

func (s *PostgreReserveStore) Delete(ctx context.Context, tx pgx.Tx, reserve model.Reserve) error {
	deleted, err := query.DeleteReserveByOrderIdAndProductId(
		ctx, tx, reserve.OrderId, reserve.ProductId,
	)
	if deleted == 0 {
		return errpkg.NotFound
	}
	return err
}