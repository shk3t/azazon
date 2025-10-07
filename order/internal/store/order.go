package store

import (
	errpkg "common/pkg/errors"
	"context"
	"errors"
	db "order/internal/database"
	"order/internal/model"
	"order/internal/query"

	"github.com/jackc/pgx/v5"
)

type PostgreOrderStore struct {
	HashPassword func(password string) (string, error)
}

func (s *PostgreOrderStore) Get(ctx context.Context, id int) (model.Order, error) {
	order, err := query.GetOrderById(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return order, errpkg.NotFound
	}

	return order, err
}

func (s *PostgreOrderStore) Save(
	ctx context.Context,
	tx pgx.Tx,
	order model.Order,
) (model.Order, error) {
	var err error

	txPassed := tx != nil
	if !txPassed {
		tx, _ = db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	}

	defer tx.Rollback(ctx)

	if order.Id == 0 {
		order.Id, err = query.CreateOrder(ctx, tx, order)
		if err != nil {
			return order, err
		}
		err = query.CreateOrderItems(ctx, tx, order.Id, order.Items)

	} else {
		err = query.UpdateOrder(ctx, tx, order)
		if err != nil {
			return order, err
		}
		query.DeleteOrderItems(ctx, tx, order.Id)
		err = query.CreateOrderItems(ctx, tx, order.Id, order.Items)
	}
	if err != nil {
		return order, err
	}

	if !txPassed {
		err = tx.Commit(ctx)
	}
	return order, err
}