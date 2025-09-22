package conversion

import (
	stockapi "common/api/stock"
	"stock/internal/model"
)

func StockProto(p *model.Product, s *model.Stock) *stockapi.Stock {
	return &stockapi.Stock{
		Product: &stockapi.Product{
			Id:    int64(p.Id),
			Name:  p.Name,
			Price: p.Price,
		},
		Quantity: int64(s.Quantity),
	}
}