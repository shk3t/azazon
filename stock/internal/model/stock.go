package model

type Product struct {
	Id    int
	Name  string
	Price float64
}

type Stock struct {
	ProductId int
	Quantity int
}