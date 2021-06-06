package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) []byte {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashPassword
}
