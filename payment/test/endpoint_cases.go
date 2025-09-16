package notificationtest

import "common/pkg/model"

const balance = 10000

var startPaymentTestCases = []struct {
	event model.OrderEvent
}{
	{
		event: model.OrderEvent{
			OrderId:   10,
			UserId:    20,
			FullPrice: 300,
		},
	},
	{
		event: model.OrderEvent{
			OrderId:   20,
			UserId:    30,
			FullPrice: 30000,
		},
	},
	{
		event: model.OrderEvent{
			OrderId:   30,
			UserId:    40,
			FullPrice: 9000,
		},
	},
}