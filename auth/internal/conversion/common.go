package conversion

import (
	"auth/internal/model"
	"common/api/common"
)

func OrderEventModel(o *common.OrderEvent) *model.OrderEvent {
	return &model.OrderEvent{
		OrderId:   int(o.OrderId),
		UserId:    int(o.UserId),
		FullPrice: int(o.FullPrice),
	}
}

func OrderEventProto(o *model.OrderEvent) *common.OrderEvent {
	return &common.OrderEvent{
		OrderId:   int64(o.OrderId),
		UserId:    int64(o.UserId),
		FullPrice: int64(o.FullPrice),
	}
}