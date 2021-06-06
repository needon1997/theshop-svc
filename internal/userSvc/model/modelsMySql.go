package model

import (
	"database/sql"
)

type Gender string

const (
	male   Gender = "male"
	female Gender = "female"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"size:100;index;unique;not null"`
	Password string `gorm:"not null"`
	NickName string `gorm:"size:20"`
	HeadUrl  string `gorm:"default:default_url""`
	Birthday sql.NullTime
	Address  string `gorm:"size:200"`
	Desc     string `gorm:"size:400"`
	Gender   Gender `gorm:"default:male"`
	Role     uint8  `gorm:"default:1"`
}

func CountUser() (int64, error) {
	var count *int64 = new(int64)
	result := db.Model(&User{}).Count(count)
	return *count, result.Error
}

func GetAllUser() ([]User, error) {
	var users []User
	result := db.Find(&users)
	return users, result.Error
}

func GetUserByOffSetLimit(offset, limit int) ([]User, error) {
	var users []User
	result := db.Model(&User{}).Limit(limit).Offset(offset).Find(&users)
	return users, result.Error
}

func SaveUser(user User) (User, error) {
	result := db.Create(&user)
	return user, result.Error
}

func GetUserById(id interface{}) (User, error) {
	var user User
	result := db.First(&user, id)
	return user, result.Error
}

func GetUserByEmail(email interface{}) (User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func UpdateUser(user User) error {
	result := db.Save(&user)
	return result.Error
}
