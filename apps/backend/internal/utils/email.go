package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func SendVerificationEmail(email, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := firstNonEmptyEnv("SMTP_USERNAME", "SMTP_USER")
	smtpPassword := normalizeSMTPSecret(firstNonEmptyEnv("SMTP_PASSWORD", "SMTP_PASS"))
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	if fromEmail == "" {
		fromEmail = smtpUsername
	}

	fmt.Println("[EMAIL] === Email Configuration Check ===")
	fmt.Printf("[EMAIL] SMTP_HOST: %s\n", smtpHost)
	fmt.Printf("[EMAIL] SMTP_PORT: %s\n", smtpPort)
	fmt.Printf("[EMAIL] SMTP_USERNAME: %s\n", smtpUsername)
	fmt.Printf("[EMAIL] SMTP_PASSWORD empty: %v\n", smtpPassword == "")
	fmt.Printf("[EMAIL] SMTP_FROM_EMAIL: %s\n", fromEmail)

	if smtpHost == "" ||
		smtpPort == "" ||
		smtpUsername == "" ||
		smtpPassword == "" ||
		fromEmail == "" {
		return fmt.Errorf("missing email configuration")
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	verificationLink := fmt.Sprintf(
		"%s/auth/verify?token=%s",
		appURL,
		token,
	)

	htmlContent := fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6;">
	<h2>Welcome to ClearRouter!</h2>

	<p>Please verify your email address by clicking the button below:</p>

	<p>
		<a href="%s"
		   style="
				display:inline-block;
				padding:12px 20px;
				background:#4CAF50;
				color:white;
				text-decoration:none;
				border-radius:6px;
		   ">
			Verify Email
		</a>
	</p>

	<p>If the button does not work, copy and paste this link:</p>

	<p>%s</p>
</body>
</html>
`, verificationLink, verificationLink)

	msg := fmt.Sprintf(
		"From: Clear Route <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: ClearRouter - Verify Your Email\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		fromEmail,
		email,
		htmlContent,
	)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	fmt.Println("[EMAIL] === Sending Email ===")
	fmt.Printf("[EMAIL] SMTP Server: %s\n", addr)
	fmt.Printf("[EMAIL] From: %s\n", fromEmail)
	fmt.Printf("[EMAIL] To: %s\n", email)

	auth := smtp.PlainAuth(
		"",
		smtpUsername,
		smtpPassword,
		smtpHost,
	)

	err := smtp.SendMail(
		addr,
		auth,
		fromEmail,
		[]string{email},
		[]byte(msg),
	)

	if err != nil {
		fmt.Printf("[EMAIL] ERROR: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("[EMAIL] SUCCESS: Verification email sent to %s\n", email)

	return nil
}

func SendEmail(to, subject, htmlContent string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := firstNonEmptyEnv("SMTP_USERNAME", "SMTP_USER")
	smtpPassword := normalizeSMTPSecret(firstNonEmptyEnv("SMTP_PASSWORD", "SMTP_PASS"))
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	if fromEmail == "" {
		fromEmail = smtpUsername
	}

	if smtpHost == "" ||
		smtpPort == "" ||
		smtpUsername == "" ||
		smtpPassword == "" ||
		fromEmail == "" {
		return fmt.Errorf("missing email configuration")
	}

	msg := fmt.Sprintf(
		"From: Clear Route <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		fromEmail,
		to,
		subject,
		htmlContent,
	)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	auth := smtp.PlainAuth(
		"",
		smtpUsername,
		smtpPassword,
		smtpHost,
	)

	err := smtp.SendMail(
		addr,
		auth,
		fromEmail,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("[EMAIL] SUCCESS: Email sent to %s\n", to)

	return nil
}

func SendEmailWithAttachments(
	to,
	subject,
	htmlContent string,
	attachments []string,
) error {
	// TODO:
	// Implement attachment support later using multipart MIME.
	// For now this just sends normal HTML email.

	_ = attachments

	return SendEmail(to, subject, htmlContent)
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}

	return ""
}

func normalizeSMTPSecret(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, " ", "")
	v = strings.ReplaceAll(v, "\n", "")
	v = strings.ReplaceAll(v, "\r", "")

	return v
}