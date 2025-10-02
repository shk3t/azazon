package service

import (
	"common/pkg/grpcutil"
	"common/pkg/model"
	"context"
	"errors"
	"time"

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

	if body.FullPrice > balance {
		return errors.New("Not enough money")
	}

	return nil
}