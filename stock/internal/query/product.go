package query

import (
	"common/pkg/helper"
	"context"
	db "stock/internal/database"
	"stock/internal/model"

	"github.com/jackc/pgx/v5"
)

func GetProductById(ctx context.Context, id int) (model.Product, error) {
	rows, _ := db.ConnPool.Query(
		ctx, `
		SELECT id, name, price
		FROM product
		WHERE id = $1`,
		id,
	)
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Product])
}

func CreateProduct(ctx context.Context, tx pgx.Tx, p model.Product) (int, error) {
	conn := helper.TxOrPool(tx, db.ConnPool)
	var id int

	err := conn.QueryRow(
		ctx, `
        INSERT INTO product (name, price)
        VALUES ($1, $2)
        RETURNING id`,
		p.Name, p.Price,
	).Scan(&id)
	return id, err
}

func UpdateProduct(ctx context.Context, tx pgx.Tx, id int, p model.Product) error {
	conn := helper.TxOrPool(tx, db.ConnPool)
	_, err := conn.Exec(
		ctx, `
		UPDATE product
		SET name = $1, price = $2
		WHERE id = $3`,
		p.Name, p.Price,
		id,
	)
	return err
}

func DeleteProduct(ctx context.Context, id int) {
	db.ConnPool.Exec(
		ctx,
		`DELETE FROM product WHERE id = $1`,
		id,
	)
}