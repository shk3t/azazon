package notificationtest

import "payment/internal/model"

var startPaymentTestCases = []struct {
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
			FullPrice: 30000,
		},
	},
	{
		order: model.OrderEvent{
			OrderId:   30,
			UserId:    40,
			FullPrice: 2000,
		},
	},
	{
		order: model.OrderEvent{
			OrderId:   40,
			UserId:    50,
			FullPrice: 200000,
		},
	},
	{
		order: model.OrderEvent{
			OrderId:   50,
			UserId:    60,
			FullPrice: 13000,
		},
	},
}