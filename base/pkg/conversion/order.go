package conversion

import (
	"base/api/order"
	"base/pkg/model"
)

func OrderEventModel(o *order.OrderEvent) *model.OrderEvent {
	return &model.OrderEvent{
		OrderId: int(o.OrderId),
		UserId:  int(o.UserId),
	}
}

func OrderEventProto(o *model.OrderEvent) *order.OrderEvent {
	return &order.OrderEvent{
		OrderId: int64(o.OrderId),
		UserId:  int64(o.UserId),
	}
}