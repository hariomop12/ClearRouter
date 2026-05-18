package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func SendVerificationEmail(email, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := firstNonEmptyEnv("SMTP_USERNAME", "SMTP_USER")
	smtpPassword := normalizeSMTPSecret(
		firstNonEmptyEnv("SMTP_PASSWORD", "SMTP_PASS"),
	)

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
	fmt.Println("[EMAIL] Starting SMTP send...")

	auth := smtp.PlainAuth(
		"",
		smtpUsername,
		smtpPassword,
		smtpHost,
	)

	err := sendMailIPv4(
		addr,
		auth,
		fromEmail,
		[]string{email},
		[]byte(msg),
		smtpHost,
	)

	if err != nil {
		fmt.Printf("[EMAIL] ERROR: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("[EMAIL] SMTP send completed")
	fmt.Printf(
		"[EMAIL] SUCCESS: Verification email sent to %s\n",
		email,
	)

	return nil
}

func SendEmail(to, subject, htmlContent string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := firstNonEmptyEnv("SMTP_USERNAME", "SMTP_USER")
	smtpPassword := normalizeSMTPSecret(
		firstNonEmptyEnv("SMTP_PASSWORD", "SMTP_PASS"),
	)

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

	err := sendMailIPv4(
		addr,
		auth,
		fromEmail,
		[]string{to},
		[]byte(msg),
		smtpHost,
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
	// Add real multipart attachment support later.
	_ = attachments

	return SendEmail(to, subject, htmlContent)
}

func sendMailIPv4(
	addr string,
	auth smtp.Auth,
	from string,
	to []string,
	msg []byte,
	host string,
) error {

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	fmt.Println("[EMAIL] Dialing SMTP server over IPv4...")

	conn, err := dialer.Dial("tcp4", addr)
	if err != nil {
		return fmt.Errorf("tcp4 dial failed: %w", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf(
			"smtp client creation failed: %w",
			err,
		)
	}

	defer client.Close()

	fmt.Println("[EMAIL] SMTP client connected")

	if ok, _ := client.Extension("STARTTLS"); ok {
		fmt.Println("[EMAIL] STARTTLS supported, upgrading connection...")

		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}

		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf(
				"STARTTLS failed: %w",
				err,
			)
		}

		fmt.Println("[EMAIL] TLS upgrade successful")
	}

	fmt.Println("[EMAIL] Authenticating SMTP session...")

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf(
			"auth failed: %w",
			err,
		)
	}

	fmt.Println("[EMAIL] SMTP authentication successful")

	if err := client.Mail(from); err != nil {
		return fmt.Errorf(
			"MAIL FROM failed: %w",
			err,
		)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf(
				"RCPT TO failed: %w",
				err,
			)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf(
			"DATA command failed: %w",
			err,
		)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf(
			"message write failed: %w",
			err,
		)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf(
			"message close failed: %w",
			err,
		)
	}

	fmt.Println("[EMAIL] Message written successfully")

	err = client.Quit()
	if err != nil {
		return fmt.Errorf(
			"SMTP quit failed: %w",
			err,
		)
	}

	return nil
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