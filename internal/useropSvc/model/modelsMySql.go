package model

import (
	"gorm.io/gorm"
	"time"
)

type Address struct {
	gorm.Model
	UserId   int32  `gorm:"not null"`
	Province string `gorm:"size:50;not null"`
	City     string `gorm:"size:50;not null"`
	Address  string `gorm:"size:100;not null"`
	Name     string `gorm:"size:100;not null"`
	Mobile   string `gorm:"size:11;not null"`
}

type UserFav struct {
	UserId    uint `gorm:"primaryKey;autoIncrement:false"`
	GoodsId   uint `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
}

type UserMessage struct {
	gorm.Model
	MessageType string `gorm:"type:enum('note', 'complain', 'inquiry', 'customer service', 'quote');default:'note'"`
	Subject     string `gorm:"type:varchar(100);not null"`
	Content     string `gorm:"not null"`
	File        string
	UserId      int32 `gorm:"not null"`
}
