package model

import (
	"gorm.io/gorm"
	"time"
)

type ShoppingCart struct {
	gorm.Model
	UserId  uint `gorm:"not null"`
	GoodsId uint `gorm:"not null"`
	Num     uint `gorm:"default:0"`
	Checked bool `gorm:"default:true"`
}

type OrderInfo struct {
	gorm.Model
	UserId      uint    `gorm:"not null"`
	OrderSn     string  `gorm:"not null;unique;size:60"`
	Status      string  `gorm:"not null;default:paying"`
	OrderAmount float32 `gorm:"not null,default:0"`
	PayAt       *time.Time
	Address     string `gorm:"size:100;default:default_address"`
	BuyerName   string `gorm:"size:50:default:default_name"`
	BuyerMobile string `gorm:"size:15"`
	Note        string `gorm:"size:300;default:default_note"`
}

type OrderGoods struct {
	gorm.Model
	OrderId       uint
	GoodsId       uint
	GoodsName     string  `gorm:"size:50:default:default_name"`
	GoodsImageURL string  `gorm:"default:default_url"`
	GoodsPrice    float32 `gorm:"not null,default:0"`
	Num           uint    `gorm:"default:0"`
}
