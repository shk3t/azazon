package consts

type TopicName string

var Topics = struct {
	OrderCreated    TopicName
	OrderConfirmed  TopicName
	OrderCanceled  TopicName
}{
	OrderCreated:    "order_created",
	OrderConfirmed:  "order_confirmed",
	OrderCanceled:  "order_canceled",
}