package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendResetEmail(toEmail, resetLink string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("EMAIL_SENDER")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	subject := "Subject: Reset Your Password\n"
	body := fmt.Sprintf("Click the link to reset your password:\n%s\n", resetLink)
	msg := []byte(subject + "\n" + body)

	addr := smtpHost + ":" + smtpPort
	return smtp.SendMail(addr, auth, from, []string{toEmail}, msg)
}
