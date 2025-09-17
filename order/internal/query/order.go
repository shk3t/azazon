package query

import (
	"common/pkg/helper"
	"context"
	db "order/internal/database"
	"order/internal/model"

	"github.com/jackc/pgx/v5"
)

func GetOrderById(ctx context.Context, id int) (model.Order, error) {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, user_id, status, address, track
		FROM "order"
		WHERE id = $1`,
		id,
	)
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Order])
}

func CreateOrder(ctx context.Context, tx pgx.Tx, o model.Order) (int, error) {
	conn := helper.TxOrPool(tx, db.ConnPool)
	var id int

	err := conn.QueryRow(
		ctx, `
        INSERT INTO "order" (user_id, status, address, track)
        VALUES ($1, $2, $3, $4)
        RETURNING id`,
		o.UserId, o.Status, o.Address, o.Track,
	).Scan(&id)
	return id, err
}

func UpdateOrder(ctx context.Context, tx pgx.Tx, id int, o model.Order) error {
	conn := helper.TxOrPool(tx, db.ConnPool)
	_, err := conn.Exec(
		ctx, `
		UPDATE "order"
		SET user_id = $1, status = $2, address = $3, track = $4
		WHERE id = $5`,
		o.UserId, o.Status, o.Address, o.Track,
		id,
	)
	return err
}

func DeleteOrder(ctx context.Context, id int) {
	db.ConnPool.Exec(
		ctx,
		`DELETE FROM "order" WHERE id = $1`,
		id,
	)
}