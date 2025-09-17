package query

import (
	"common/pkg/helper"
	"context"
	db "order/internal/database"
	"order/internal/model"

	"github.com/jackc/pgx/v5"
)

func GetItemsByOrderId(ctx context.Context, orderId int) ([]model.Item, error) {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, user_id, status, address, track
		FROM item
		WHERE order_id = $1`,
		orderId,
	)
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Item])
}

func CreateOrderItems(ctx context.Context, tx pgx.Tx, orderId int, items []model.Item) error {
	conn := helper.TxOrPool(tx, db.ConnPool)

	entries := make([][]any, len(items))
	for i, item := range items {
		entries[i] = []any{orderId, item.Id}
	}

	_, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"item"},
		[]string{"order_id", "quantity"},
		pgx.CopyFromRows(entries),
	)

	return err
}

func DeleteOrderItems(ctx context.Context, tx pgx.Tx, orderId int) {
	conn := helper.TxOrPool(tx, db.ConnPool)
	conn.Exec(
		ctx,
		`DELETE FROM item WHERE order_id = $1`,
		orderId,
	)
}