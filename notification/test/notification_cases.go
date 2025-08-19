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
}