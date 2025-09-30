package query

import (
	"common/pkg/helper"
	"context"
	db "stock/internal/database"
	"stock/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
)

func GetReserveByOrderIdAndProductId(
	ctx context.Context,
	orderId int,
	productId int,
) (model.Reserve, error) {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, user_id, order_id, product_id, quantity
		FROM reserve
		WHERE order_id = $1 AND product_id = $2`,
		orderId, productId,
	)
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Reserve])
}

func GetReservesOlderThan(ctx context.Context, dt time.Time) ([]model.Reserve, error) {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, user_id, order_id, product_id, quantity
		FROM reserve
		WHERE created_at < $1`,
		dt,
	)
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Reserve])
}

func CreateReserve(ctx context.Context, tx pgx.Tx, r model.Reserve) (int, error) {
	conn := helper.TxOrPool(tx, db.ConnPool)

	var id int
	err := conn.QueryRow(
		ctx, `
        INSERT INTO reserve (user_id, order_id, product_id, quantity, created_at)
        VALUES ($1, $2, $3, $4 NOW())
        RETURNING id`,
		r.UserId, r.OrderId, r.ProductId, r.Quantity,
	).Scan(&id)
	return id, err
}

func DeleteReserveByOrderIdAndProductId(
	ctx context.Context,
	tx pgx.Tx,
	orderId int,
	productId int,
) (int, error) {
	conn := helper.TxOrPool(tx, db.ConnPool)

	var id int
	err := conn.QueryRow(
		ctx,
		`DELETE FROM reserve WHERE order_id = $1 AND product_id = $2`,
		orderId, productId,
	).Scan(&id)
	return id, err
}