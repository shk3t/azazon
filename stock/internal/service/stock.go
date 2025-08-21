package service

import "common/pkg/grpcutil"

var NewErr = grpcutil.NewError
var NewInternalErr = grpcutil.NewInternalError

type stockStore interface {
}

type StockService struct {
	store stockStore
}

func NewStockService() *StockService {
	return &StockService{
		store: nil,
	}
}