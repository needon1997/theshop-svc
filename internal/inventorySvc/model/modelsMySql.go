package model

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	GoodsID uint `gorm:"unique;not null"`
	Stocks  uint `gorm:"default:0""`
	Version uint `gorm:"default:0"`
}

type InventoryHistory struct {
	gorm.Model
	OrderSn string `gorm:"unique;not null"`
}
