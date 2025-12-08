package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint
	TotalCents int         `json:"total_cents"`
	Status     string      `json:"status"`
	Items      []OrderItem `json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     int     `json:"price"`
}
