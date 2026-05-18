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
	from := os.Getenv("MAILGUN_FROM_EMAIL")

	if domain == "" || apiKey == "" || from == "" {
		return fmt.Errorf("missing mailgun configuration")
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	verifyLink := fmt.Sprintf("%s/auth/verify?token=%s", appURL, token)

	subject := "ClearRouter - Verify Your Email"
	body := fmt.Sprintf(`
		<h2>Welcome to ClearRouter</h2>
		<p>Click below to verify your email:</p>
		<a href="%s">Verify Email</a>
		<p>%s</p>
	`, verifyLink, verifyLink)

	mg := mailgun.NewMailgun(domain, apiKey)

	msg := mg.NewMessage(
		from,
		subject,
		"",        // plain text empty
		email,
	)

	msg.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, id, err := mg.Send(ctx, msg)
	if err != nil {
		fmt.Println("[MAILGUN] ERROR:", err)
		return err
	}

	fmt.Println("[MAILGUN] SENT ID:", id)
	return nil
}