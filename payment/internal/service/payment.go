package service

import (
	"common/pkg/grpcutil"
	"common/pkg/model"
	"context"
	"errors"
	"fmt"
	"payment/internal/query"
	"time"

	"github.com/jackc/pgerrcode"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
)

var NewErr = grpcutil.NewServiceError
var NewInternalErr = grpcutil.NewInternalError

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) StartPayment(
	ctx context.Context,
	tx pgx.Tx,
	body model.OrderEvent,
) error {
	balance := 10000.00

	time.Sleep(1 * time.Second)

	err := query.CreateProcessedPayment(ctx, tx, body.OrderId)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return NewInternalErr(fmt.Errorf("Duplicated order for payment | %w", err))
		}
		return NewInternalErr(err)
	}

	if body.FullPrice > balance {
		return errors.New("Not enough money")
	}

	return nil
}