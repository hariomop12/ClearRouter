package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func SendVerificationEmail(email, token string) error {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")
	fromEmail := os.Getenv("MAILGUN_FROM_EMAIL")

	if domain == "" || apiKey == "" || fromEmail == "" {
		return fmt.Errorf("missing Mailgun config: MAILGUN_DOMAIN / MAILGUN_API_KEY / MAILGUN_FROM_EMAIL")
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	verifyLink := fmt.Sprintf("%s/auth/verify?token=%s", appURL, token)

	subject := "ClearRouter - Verify Your Email"

	html := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial;">
			<h2>Verify your email</h2>
			<p>Click below to verify:</p>

			<a href="%s"
			   style="padding:10px 18px;
			   background:#4CAF50;
			   color:white;
			   text-decoration:none;
			   border-radius:5px;">
			   Verify Email
			</a>

			<p>%s</p>
		</body>
		</html>
	`, verifyLink, verifyLink)

	fmt.Println("[MAILGUN] Initializing client...")

	mg := mailgun.NewMailgun(domain, apiKey)

	msg := mg.NewMessage(
		fromEmail,
		subject,
		"", // plain text empty
		email,
	)

	msg.SetHtml(html)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("[MAILGUN] Sending email to:", email)

	res, id, err := mg.Send(ctx, msg)
	if err != nil {
		fmt.Printf("[MAILGUN] ERROR: %v\n", err)
		return fmt.Errorf("mailgun send failed: %w", err)
	}

	fmt.Println("[MAILGUN] Response:", res)
	fmt.Println("[MAILGUN] Message ID:", id)
	fmt.Println("[MAILGUN] SUCCESS: Email sent")

	return nil
}