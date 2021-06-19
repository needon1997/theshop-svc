package model

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	CategoryName     string `gorm:"not null;size:50;unique" json:"name"`
	Level            uint   `gorm:"default:1;not null" json:"level"`
	IsShow           bool   `gorm:"default:false" json:"is_show"`
	ParentCategoryID uint   `gorm:"default:null" json:"parent_category_id"`
}
type Brand struct {
	gorm.Model
	Name string `gorm:"unique;index;size:50"`
	Logo string `gorm:"size:256,default:;"`
}

type Goods struct {
	gorm.Model
	CategoryID  uint
	Category    Category
	BrandID     uint
	Brand       Brand
	IsShow      bool    `gorm:"default:false;not null"`
	GoodsSn     string  `gorm:"size:50;unique;not null"`
	Name        string  `gorm:"size:100;default:;not null"`
	ClickNum    uint    `gorm:"default:0"`
	SoldNum     uint    `gorm:"default:0"`
	FavNum      uint    `gorm:"default:0"`
	Stock       uint    `gorm:"default:0"`
	MarketPrice float32 `gorm:"default:0"`
	ShopPrice   float32 `gorm:"default:0"`
	Brief       string  `gorm:"size:200"`
	Description string  `gorm:"size:500"`
	ShipFree    bool    `gorm:"default:false"`
	Images      string
	DescImages  string
	FrontImage  string
	IsNew       bool `gorm:"default:false"`
	IsHot       bool `gorm:"default:false"`
}

type GoodsBrandCategory struct {
	gorm.Model
	CategoryID uint `gorm:"not null"`
	Category   Category
	BrandID    uint `gorm:"not null"`
	Brand      Brand
}

type Banner struct {
	gorm.Model
	Image string `gorm:"default:;size:256"`
	Url   string `gorm:"default:;size:256"`
	Index uint   `gorm:"default:0"`
}
