package consts

type ServiceName string

var Services = struct {
	Auth         ServiceName
	Notification ServiceName
	Order        ServiceName
	Payment      ServiceName
	Stock        ServiceName
}{
	Auth:         "auth",
	Notification: "notification",
	Order:        "order",
	Payment:      "payment",
	Stock:        "stock",
}