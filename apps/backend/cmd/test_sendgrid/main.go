package main

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		fmt.Println("SENDGRID_API_KEY is not set")
		return
	}

	fmt.Printf("API Key length: %d\n", len(apiKey))
	fmt.Printf("API Key starts with: %s\n", apiKey[:6])

	from := mail.NewEmail("ClearRouter Test", "noreply@clearrouter.com")
	subject := "SendGrid Test Email"
	to := mail.NewEmail("Test User", "test@example.com") // Replace with your email for testing
	plainTextContent := "SendGrid test email"
	htmlContent := "<strong>SendGrid test email</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)

	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}

	fmt.Printf("Response status code: %d\n", response.StatusCode)
	fmt.Printf("Response body: %s\n", response.Body)
	fmt.Printf("Response headers: %v\n", response.Headers)
}
