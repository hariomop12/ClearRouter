package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(email, token string) error {
	from := os.Getenv("SMTP_FROM_EMAIL")
	if from == "" {
		return fmt.Errorf("SMTP_FROM_EMAIL not set")
	}

	username := os.Getenv("SMTP_USER")
	if username == "" {
		return fmt.Errorf("SMTP_USER not set")
	}

	password := os.Getenv("SMTP_PASS")
	if password == "" {
		return fmt.Errorf("SMTP_PASS not set")
	}

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return fmt.Errorf("SMTP_HOST not set")
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		return fmt.Errorf("SMTP_PORT not set")
	}

	fmt.Printf("Attempting to send email using:\nHost: %s\nPort: %s\nUsername: %s\nFrom: %s\n", host, port, username, from)

	verificationLink := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", token)

	// Create email headers and body
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
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
		from, email, verificationLink, verificationLink))

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	fmt.Printf("Sending email to: %s via SMTP: %s\n", email, addr)
	err := smtp.SendMail(addr, auth, from, []string{email}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
