package service

import (
	"common/pkg/grpcutil"
	"context"
	"errors"
	"common/pkg/model"
	"time"
)

var NewErr = grpcutil.NewServiceError
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

	time.Sleep(1 * time.Second)

	if body.FullPrice > balance {
		return errors.New("Not enough money")
	}

	return nil
}