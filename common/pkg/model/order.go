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
	Id       int
	Quantity int
}

type OrderStatus int

const (
	Confirmed OrderStatus = iota + 1
	Cancelled
)