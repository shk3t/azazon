package notificationtest

import "common/pkg/model"

var orderCreatedTestCases = []struct {
	order model.OrderEvent
}{
	{
		order: model.OrderEvent{
			OrderId: 10,
			UserId:  20,
		},
	},
	{
		order: model.OrderEvent{
			OrderId: 20,
			UserId:  30,
		},
	},
	{
		order: model.OrderEvent{
			OrderId: 30,
			UserId:  40,
		},
	},
	{
		order: model.OrderEvent{
			OrderId: 40,
			UserId:  50,
		},
	},
	{
		order: model.OrderEvent{
			OrderId: 50,
			UserId:  60,
		},
	},
}