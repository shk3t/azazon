package store

import (
	errpkg "common/pkg/errors"
	"context"
	"errors"
	"stock/internal/model"
	"stock/internal/query"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgreReserveStore struct {
	HashPassword func(password string) (string, error)
}

func (s *PostgreReserveStore) Get(
	ctx context.Context,
	orderId int,
) ([]model.Reserve, error) {
	reserves, err := query.GetReservesByOrderId(ctx, orderId)

	if errors.Is(err, pgx.ErrNoRows) {
		return reserves, errpkg.NotFound
	}

	return reserves, err
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
	_, err := query.CreateReserve(ctx, tx, reserve)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return ProductDoesNotExistError
			}
		}
	}
	return err
}

func (s *PostgreReserveStore) Delete(ctx context.Context, tx pgx.Tx, reserve model.Reserve) error {
	var deleted int
	var err error

	if reserve.ProductId != 0 {
		deleted, err = query.DeleteReserveByOrderIdAndProductId(
			ctx, tx, reserve.OrderId, reserve.ProductId,
		)
	} else {
		deleted, err = query.DeleteReservesByOrderId(ctx, tx, reserve.OrderId)
	}

	if deleted == 0 {
		return errpkg.NotFound
	}
	return err
}