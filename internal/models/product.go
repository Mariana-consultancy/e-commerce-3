package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	SellerID    uint    `json:"seller_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"image_url"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
	Status      bool    `json:"status"`
	Orders      []Order `json:"orders" gorm:"many2many:order_items;"`
}

type Cart struct {
	gorm.Model
	Quantity  int  `json:"quantity"`
	ProductID uint `json:"product_id"`
	UserID    uint `json:"user_id"`
}
