package consts

type TopicName string

var Topics = struct {
	OrderCreated    TopicName
	OrderConfirmed  TopicName
	OrderCancelling TopicName
	OrderCancelled  TopicName
}{
	OrderCreated:    "order_created",
	OrderConfirmed:  "order_confirmed",
	OrderCancelling: "order_cancelling",
	OrderCancelled:  "order_canceled",
}