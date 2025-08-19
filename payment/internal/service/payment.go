package service

import "base/pkg/grpcutil"

var NewErr = grpcutil.NewError
var NewInternalErr = grpcutil.NewInternalError

type paymentStore interface {
}

type PaymentService struct {
	store paymentStore
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		store: nil,
	}
}