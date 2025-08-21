package service

import (
	"common/pkg/grpcutil"
	"context"
	"errors"
	"payment/internal/model"
	"time"
)

var NewErr = grpcutil.NewError
var NewInternalErr = grpcutil.NewInternalError

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) StartPayment(
	ctx context.Context,
	body model.OrderEvent,
) error {
	balance := 10000

	time.Sleep(2 * time.Second)

	if body.FullPrice > balance {
		return errors.New("Not enough money")
	}

	return nil
}