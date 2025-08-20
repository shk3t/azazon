package notificationtest

import "base/pkg/model"

var orderCreatedTestCases = []struct {
	order model.OrderEvent
}{
	{
		order: model.OrderEvent{
			Id:     10,
			UserId: 20,
		},
	},
	{
		order: model.OrderEvent{
			Id:     20,
			UserId: 30,
		},
	},
	{
		order: model.OrderEvent{
			Id:     30,
			UserId: 40,
		},
	},
	{
		order: model.OrderEvent{
			Id:     40,
			UserId: 50,
		},
	},
	{
		order: model.OrderEvent{
			Id:     50,
			UserId: 60,
		},
	},
}