package model

type Order struct {
	Id      int
	UserId  int
	Items   []Item
	Status  OrderStatus
	Address string
	Track   string
}

type Item struct {
	ProductId int
	Quantity  int
}

type OrderStatus int

const (
	ConfirmedStatus OrderStatus = 1
	CancelledStatus OrderStatus = 3
)