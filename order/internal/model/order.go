package model

type Order struct {
	Id      int
	UserId  int
	Items   []Item `db:"-"`
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
	CreatedStatus   OrderStatus = 1
	ConfirmedStatus OrderStatus = 2
	CanceledStatus  OrderStatus = 3
)

func NewStatus(str string) OrderStatus {
	switch str {
	case "created":
		return 1
	case "confirmed":
		return 2
	case "canceled":
		return 3
	}
	return -1
}

func (s OrderStatus) String() string {
	switch s {
	case 1:
		return "created"
	case 2:
		return "confirmed"
	case 3:
		return "canceled"
	}
	return ""
}