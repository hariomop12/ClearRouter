package main

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// SendGrid credentials
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	apiKey := os.Getenv("API_KEY")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	toEmail := "test@example.com" // Replace with your test email

	// Print configuration
	fmt.Printf("SMTP Host: %s\n", smtpHost)
	fmt.Printf("SMTP Port: %s\n", smtpPort)
	fmt.Printf("From Email: %s\n", fromEmail)
	fmt.Printf("API Key Length: %d\n", len(apiKey))

	// Create message
	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Test Email\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
		"This is a test email from ClearRouter",
		fromEmail, toEmail))

	// Authenticate
	auth := smtp.PlainAuth("", "apikey", apiKey, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Send email
	err := smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, message)
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}

	fmt.Println("Test email sent successfully!")
}
