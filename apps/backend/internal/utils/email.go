package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(email, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	apiKey := os.Getenv("SMTP_USER") // SendGrid API Key
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	if smtpHost == "" || smtpPort == "" || apiKey == "" || fromEmail == "" {
		return fmt.Errorf("missing email configuration")
	}

	verificationLink := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", token)

	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: ClearRouter - Verify Your Email\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"<html><body>"+
		"<h2>Welcome to ClearRouter!</h2>"+
		"<p>Please click the following link to verify your email:</p>"+
		"<p><a href='%s'>Click here to verify your email</a></p>"+
		"<p>If the button doesn't work, copy and paste this link in your browser:</p>"+
		"<p>%s</p>"+
		"</body></html>",
		fromEmail, email, verificationLink, verificationLink))

	auth := smtp.PlainAuth("", "apikey", apiKey, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	fmt.Printf("Sending verification email to: %s using SendGrid\n", email)

	err := smtp.SendMail(addr, auth, fromEmail, []string{email}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
