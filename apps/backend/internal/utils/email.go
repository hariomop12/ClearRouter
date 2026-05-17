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
	smtpPassword := normalizeSMTPSecret(firstNonEmptyEnv("SMTP_PASSWORD", "SMTP_PASS"))
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = smtpUsername
	}

	fmt.Println("[EMAIL] === Email Configuration Check ===")
	fmt.Printf("[EMAIL] SMTP_HOST: %s (empty: %v)\n", smtpHost, smtpHost == "")
	fmt.Printf("[EMAIL] SMTP_PORT: %s (empty: %v)\n", smtpPort, smtpPort == "")
	fmt.Printf("[EMAIL] SMTP_USERNAME: %s (empty: %v)\n", smtpUsername, smtpUsername == "")
	fmt.Printf("[EMAIL] SMTP_PASSWORD: ****** (empty: %v)\n", smtpPassword == "")
	fmt.Printf("[EMAIL] SMTP_FROM_EMAIL: %s (empty: %v)\n", fromEmail, fromEmail == "")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || fromEmail == "" {
		fmt.Println("[EMAIL] ERROR: Missing email configuration")
		return fmt.Errorf("missing email configuration. Required: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM_EMAIL")
	}

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

	msg := fmt.Sprintf("From: Clear Route <%s>\r\nTo: %s\r\nSubject: ClearRouter - Verify Your Email\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s\r\n", fromEmail, email, htmlContent)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	fmt.Printf("[EMAIL] Attempting to send verification email\n")
	fmt.Printf("[EMAIL] To: %s\n", email)
	fmt.Printf("[EMAIL] From: %s\n", fromEmail)
	fmt.Printf("[EMAIL] SMTP Server: %s\n", addr)
	fmt.Printf("[EMAIL] Message length: %d bytes\n", len(msg))
	fmt.Println("[EMAIL] === Connecting to SMTP Server ===")

	// Test DNS resolution
	fmt.Printf("[EMAIL] Testing DNS resolution for %s...\n", smtpHost)
	ips, dnsErr := net.LookupIP(smtpHost)
	if dnsErr != nil {
		fmt.Printf("[EMAIL] ERROR: DNS resolution failed: %v\n", dnsErr)
	} else {
		fmt.Printf("[EMAIL] DNS resolved to: %v\n", ips)
	}

	// Test connectivity
	fmt.Printf("[EMAIL] Testing connectivity to %s...\n", addr)
	conn, connErr := net.DialTimeout("tcp", addr, 10*time.Second)
	if connErr != nil {
		fmt.Printf("[EMAIL] ERROR: Cannot connect to SMTP server: %v\n", connErr)
		return fmt.Errorf("cannot connect to SMTP server: %w", connErr)
	}
	conn.Close()
	fmt.Println("[EMAIL] ✓ TCP connection successful")

	// Attempt SMTP with TLS
	fmt.Println("[EMAIL] Creating TLS dialer...")
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	fmt.Println("[EMAIL] Attempting to send via smtp.SendMail with STARTTLS...")
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(addr, auth, fromEmail, []string{email}, []byte(msg))
	if err != nil {
		fmt.Printf("[EMAIL] ERROR: Failed to send email: %v\n", err)
		fmt.Printf("[EMAIL] Error Type: %T\n", err)
		fmt.Printf("[EMAIL] Error String: %s\n", err.Error())

		// Try alternative method with an explicit client + STARTTLS/implicit TLS
		fmt.Println("[EMAIL] === Trying alternative method with explicit SMTP client ===")
		err2 := sendEmailWithSMTPClient(addr, smtpHost, smtpPort, auth, fromEmail, email, msg, tlsConfig)
		if err2 != nil {
			fmt.Printf("[EMAIL] ERROR: Alternative method also failed: %v\n", err2)
			return fmt.Errorf("failed to send email (both methods): standard: %w, tls-explicit: %w", err, err2)
		}
		fmt.Println("[EMAIL] SUCCESS: Email sent via alternative TLS method")
		return nil
	}

	fmt.Printf("[EMAIL] SUCCESS: Email sent successfully to %s\n", email)
	return nil
}

func sendEmailWithSMTPClient(addr, smtpHost, smtpPort string, auth smtp.Auth, from, to, msg string, tlsConfig *tls.Config) error {
	// Port 465 is implicit TLS. Port 587 is plaintext then STARTTLS.
	if smtpPort == "465" {
		fmt.Println("[EMAIL-TLS] Connecting with implicit TLS (port 465)...")
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			fmt.Printf("[EMAIL-TLS] ERROR: TLS dial failed: %v\n", err)
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, smtpHost)
		if err != nil {
			fmt.Printf("[EMAIL-TLS] ERROR: SMTP client creation failed: %v\n", err)
			return err
		}
		defer client.Close()

		return smtpClientSend(client, auth, from, to, msg)
	}

	fmt.Println("[EMAIL-STARTTLS] Connecting in plaintext then upgrading with STARTTLS...")
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		fmt.Printf("[EMAIL-STARTTLS] ERROR: TCP dial failed: %v\n", err)
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Printf("[EMAIL-STARTTLS] ERROR: SMTP client creation failed: %v\n", err)
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(tlsConfig); err != nil {
			fmt.Printf("[EMAIL-STARTTLS] ERROR: STARTTLS failed: %v\n", err)
			return err
		}
	} else {
		fmt.Println("[EMAIL-STARTTLS] WARNING: Server does not advertise STARTTLS; continuing without TLS")
	}

	return smtpClientSend(client, auth, from, to, msg)
}

func smtpClientSend(client *smtp.Client, auth smtp.Auth, from, to, msg string) error {
	fmt.Println("[EMAIL-CLIENT] Authenticating...")
	if err := client.Auth(auth); err != nil {
		fmt.Printf("[EMAIL-CLIENT] ERROR: Authentication failed: %v\n", err)
		return err
	}

	fmt.Println("[EMAIL-CLIENT] Sending mail...")
	if err := client.Mail(from); err != nil {
		fmt.Printf("[EMAIL-CLIENT] ERROR: From address error: %v\n", err)
		return err
	}
	if err := client.Rcpt(to); err != nil {
		fmt.Printf("[EMAIL-CLIENT] ERROR: Recipient error: %v\n", err)
		return err
	}

	w, err := client.Data()
	if err != nil {
		fmt.Printf("[EMAIL-CLIENT] ERROR: Data writer error: %v\n", err)
		return err
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		fmt.Printf("[EMAIL-CLIENT] ERROR: Write error: %v\n", err)
		return err
	}
	if err := w.Close(); err != nil {
		fmt.Printf("[EMAIL-CLIENT] ERROR: Close error: %v\n", err)
		return err
	}

	_ = client.Quit()
	fmt.Println("[EMAIL-CLIENT] ✓ Email sent successfully")
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

	fmt.Println("[EMAIL] === SendEmail Configuration Check ===")
	fmt.Printf("[EMAIL] To: %s\n", to)
	fmt.Printf("[EMAIL] Subject: %s\n", subject)
	fmt.Printf("[EMAIL] SMTP_HOST: %s\n", smtpHost)
	fmt.Printf("[EMAIL] SMTP_PORT: %s\n", smtpPort)

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || fromEmail == "" {
		fmt.Println("[EMAIL] ERROR: Missing email configuration in SendEmail")
		return fmt.Errorf("missing email configuration")
	}

	msg := fmt.Sprintf("From: Clear Route <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s\r\n", fromEmail, to, subject, htmlContent)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	fmt.Printf("[EMAIL] Connecting to SMTP server: %s\n", addr)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(msg))
	if err != nil {
		fmt.Printf("[EMAIL] ERROR in SendEmail: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("[EMAIL] SUCCESS: Email sent to %s\n", to)
	return nil
}

func SendEmailWithAttachments(to, subject, htmlContent string, attachments []string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := firstNonEmptyEnv("SMTP_USERNAME", "SMTP_USER")
	smtpPassword := normalizeSMTPSecret(firstNonEmptyEnv("SMTP_PASSWORD", "SMTP_PASS"))
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = smtpUsername
	}

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || fromEmail == "" {
		return fmt.Errorf("missing email configuration")
	}

	boundary := "myboundary"
	msg := fmt.Sprintf("From: Clear Route <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n", fromEmail, to, subject, boundary)

	msg += fmt.Sprintf("--%s\r\n", boundary)
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += htmlContent + "\r\n\r\n"

	for _, attachment := range attachments {
		// Add attachment handling if needed
		_ = attachment
	}

	msg += fmt.Sprintf("--%s--\r\n", boundary)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

func normalizeSMTPSecret(v string) string {
	// Common copy/paste format for Gmail "App Passwords" includes spaces.
	// SMTP auth expects the raw 16-char password (no spaces/newlines).
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, " ", "")
	return v
}
