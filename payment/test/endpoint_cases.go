package notificationtest

import "payment/internal/model"

const balance = 10000

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
			FullPrice: 9000,
		},
	},
}