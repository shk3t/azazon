package model

import "time"

type Product struct {
	Id    int
	Name  string
	Price float64
}

type Stock struct {
	ProductId int
	Quantity  int
}

type Reserve struct {
	UserId    int
	OrderId   int
	ProductId int
	Quantity  int
	CreatedAt time.Time
}