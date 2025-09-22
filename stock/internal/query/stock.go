package query

import (
	"common/pkg/helper"
	"context"
	db "stock/internal/database"
	"stock/internal/model"

	"github.com/jackc/pgx/v5"
)

func GetStockByProductId(ctx context.Context, productId int) (model.Stock, error) {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, product_id, quantity
		FROM stock
		WHERE product_id = $1`,
		productId,
	)
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Stock])
}

func CreateStock(ctx context.Context, tx pgx.Tx, s model.Stock) (int, error) {
	conn := helper.TxOrPool(tx, db.ConnPool)
	var id int

	err := conn.QueryRow(
		ctx, `
        INSERT INTO stock (product_id, quantity)
        VALUES ($1, $2)
        RETURNING id`,
		s.ProductId, s.Quantity,
	).Scan(&id)
	return id, err
}

func UpdateStockByProductId(ctx context.Context, tx pgx.Tx, s model.Stock) error {
	conn := helper.TxOrPool(tx, db.ConnPool)
	_, err := conn.Exec(
		ctx, `
		UPDATE stock
		SET quantity = $1
		WHERE product_id = $2`,
		s.Quantity,
		s.ProductId,
	)
	return err
}

func DeleteStock(ctx context.Context, id int) {
	db.ConnPool.Exec(
		ctx,
		`DELETE FROM stock WHERE id = $1`,
		id,
	)
}