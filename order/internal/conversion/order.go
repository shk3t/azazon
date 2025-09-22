package conversion

import (
	orderapi "common/api/order"
	modelpkg "common/pkg/model"
	"order/internal/model"
)

func OrderEvent(o *model.Order) *modelpkg.OrderEvent {
	return &modelpkg.OrderEvent{
		OrderId: int(o.Id),
		UserId:  int(o.UserId),
	}
}

func GetOrderInfoResponse(o *model.Order) *orderapi.GetOrderInfoResponse {
	respItems := make([]*orderapi.Item, len(o.Items))
	for i, item := range o.Items {
		respItems[i] = &orderapi.Item{
			ProductId: int64(item.ProductId),
			Quantity:  int64(item.Quantity),
		}
	}

	return &orderapi.GetOrderInfoResponse{
		OrderId: int64(o.Id),
		Items:   respItems,
		Status:  orderapi.OrderStatus(o.Status),
		Address: o.Address,
		Track:   o.Track,
	}
}
