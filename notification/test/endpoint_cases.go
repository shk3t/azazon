package notificationtest

import "notification/internal/model"

var orderCreatedTestCases = []struct {
	order model.OrderEvent
}{
	{
		order: model.OrderEvent{
			OrderId:   10,
			UserId:    20,
			FullPrice: 300,
		},
	},
	{
		order: model.OrderEvent{
			OrderId:   20,
			UserId:    30,
			FullPrice: 300,
		},
	},
	{
		order: model.OrderEvent{
			OrderId:   30,
			UserId:    40,
			FullPrice: 300,
		},
	},
	{
		order: model.OrderEvent{
			OrderId:   40,
			UserId:    50,
			FullPrice: 300,
		},
	},
	{
		order: model.OrderEvent{
			OrderId:   50,
			UserId:    60,
			FullPrice: 300,
		},
	},
}