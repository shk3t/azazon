package conversion

import (
	modelpkg "common/pkg/model"
	"order/internal/model"
)

func OrderEvent(o *model.Order) *modelpkg.OrderEvent {
	return &modelpkg.OrderEvent{
		OrderId: int(o.Id),
		UserId:  int(o.UserId),
	}
}