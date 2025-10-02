package query

import (
	"common/pkg/helper"
	"context"
	db "payment/internal/database"

	"github.com/jackc/pgx/v5"
)

func CreateProcessedPayment(ctx context.Context, tx pgx.Tx, orderId int) error {
	conn := helper.TxOrPool(tx, db.ConnPool)

	_, err := conn.Exec(
		ctx, `
        INSERT INTO processed_payment (order_id)
        VALUES ($1)`,
		orderId,
	)

	return err
}