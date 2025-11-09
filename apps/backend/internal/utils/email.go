package utils

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendVerificationEmail(email, token string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	fmt.Printf("Debug - SendGrid API Key length: %d\n", len(apiKey))
	fmt.Printf("Debug - First 6 chars of API Key: %s\n", apiKey[:6])

	if apiKey == "" || fromEmail == "" {
		return fmt.Errorf("missing email configuration. Required: API_KEY, SMTP_FROM_EMAIL")
	}

	// Create new SendGrid client
	client := sendgrid.NewSendClient(apiKey)

	// Create email message
	from := mail.NewEmail("ClearRouter", fromEmail)
	to := mail.NewEmail("", email)
	subject := "ClearRouter - Verify Your Email"

	verificationLink := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", token)
	htmlContent := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Welcome to ClearRouter!</h2>
			<p>Please click the following link to verify your email:</p>
			<p>
				<a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">
					Click here to verify your email
				</a>
			</p>
			<p>If the button doesn't work, copy and paste this link in your browser:</p>
			<p>%s</p>
		</body>
		</html>
	`, verificationLink, verificationLink)

	// Create email with HTML content
	message := mail.NewV3MailInit(from, subject, to, mail.NewContent("text/html", htmlContent))

	fmt.Printf("Sending verification email to: %s using SendGrid\n", email)

	// Send email
	response, err := client.Send(message)
	if err != nil {
		fmt.Printf("Email sending error: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		fmt.Printf("SendGrid API error: Status Code %d, Body: %s\n", response.StatusCode, response.Body)
		return fmt.Errorf("SendGrid API error: %s", response.Body)
	}

	fmt.Printf("Email sent successfully to %s (Status: %d)\n", email, response.StatusCode)
	return nil
}
