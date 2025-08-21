package service

import (
	"common/pkg/grpcutil"
	"common/pkg/model"
	"context"
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

	if body.FullPrice > balance {
		// TODO: cancel
	}

	// TODO: confirm

	return nil
}