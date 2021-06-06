package utils

import (
	"fmt"
	config2 "github.com/needon1997/theshop-svc/internal/common/config"
	"go.uber.org/zap"
	"math/rand"
	"net/smtp"
)

const VERIFICATION_DIGITS = 7
const MAX = 10
const BASE = 48

func GenerateVerificationCode() string {
	code := make([]byte, VERIFICATION_DIGITS)
	for i := 0; i < VERIFICATION_DIGITS; i++ {
		code[i] = byte(BASE + (rand.Float32() * float32(MAX)))
	}
	return string(code)
}

func SendVerificationEmail(to, subject, body string) {
	host := config2.ServerConfig.EmailConfig.Host
	port := config2.ServerConfig.EmailConfig.Port
	from := config2.ServerConfig.EmailConfig.Username
	pass := config2.ServerConfig.EmailConfig.Password
	mime := config2.ServerConfig.EmailConfig.Mime
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n%s\n%s", from, to, subject, mime, body)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", host, port),
		smtp.PlainAuth("", from, pass, host),
		from, []string{to}, []byte(msg))

	if err != nil {
		zap.S().Errorw("smtp error", "error", err.Error())
		return
	}
}
