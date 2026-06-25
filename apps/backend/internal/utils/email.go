package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/resend/resend-go/v2"
)

// SendVerificationEmail triggers a verification email via Resend and logs the entire lifecycle.
func SendVerificationEmail(email, token string) error {
	fmt.Println("[RESEND-LOG] ----- Starting Email Verification Process -----")

	// 1. Fetch and log environment configuration
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("RESEND_FROM_EMAIL")
	appURL := os.Getenv("APP_URL")

	// Mask API key for security logs while verifying its length
	maskedKey := "MISSING"
	if len(apiKey) > 8 {
		maskedKey = apiKey[:6] + "..." + apiKey[len(apiKey)-4:]
	} else if apiKey != "" {
		maskedKey = "PRESENT (Too short)"
	}

	fmt.Printf("[RESEND-LOG] Config Loaded -> API_KEY: %s | FROM: '%s' | APP_URL: '%s'\n", maskedKey, from, appURL)

	// 2. Validate critical configurations
	if apiKey == "" {
		err := fmt.Errorf("missing RESEND_API_KEY environment variable")
		fmt.Printf("[RESEND-LOG] Validation Failed: %v\n", err)
		return err
	}

	if from == "" {
		err := fmt.Errorf("missing RESEND_FROM_EMAIL environment variable")
		fmt.Printf("[RESEND-LOG] Validation Failed: %v\n", err)
		return err
	}

	// 3. Set default fallback URL if empty
	if appURL == "" {
		appURL = "http://localhost:8080"
		fmt.Printf("[RESEND-LOG] Warning: APP_URL empty, defaulting to: %s\n", appURL)
	}

	// 4. Construct payload metrics
	verifyLink := fmt.Sprintf("%s/auth/verify?token=%s", appURL, token)
	subject := "ClearRouter - Verify Your Email"
	
	fmt.Printf("[RESEND-LOG] Target Payload -> Recipient: '%s' | Sender Domain: '%s'\n", email, from)
	fmt.Printf("[RESEND-LOG] Generated Verification Link: %s\n", verifyLink)

	html := fmt.Sprintf(`
		<h2>Welcome to ClearRouter</h2>
		<p>Click below to verify your email:</p>
		<a href="%s">Verify Email</a>
		<p>%s</p>
	`, verifyLink, verifyLink)

	// 5. Initialize client and build payload parameters
	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{email},
		Subject: subject,
		Html:    html,
	}

	// Extra sanity check warning if the sender string does not contain your verified domain
	if !strings.Contains(from, "clearrouter.hariomop.in") {
		fmt.Printf("[RESEND-LOG] WARNING: The 'From' address ('%s') does not match your verified domain 'clearrouter.hariomop.in'. Resend may reject this request.\n", from)
	}

	// 6. Execute transmission
	fmt.Println("[RESEND-LOG] Dispatching API request to Resend...")
	sent, err := client.Emails.Send(params)
	
	if err != nil {
		fmt.Printf("[RESEND-LOG] API Execution Error: %v\n", err)
		fmt.Println("[RESEND-LOG] ----- Process Aborted (Failure) -----")
		return err
	}

	// 7. Track successful output metadata
	fmt.Printf("[RESEND-LOG] Success! Email dispatched successfully. Response Message ID: %s\n", sent.Id)
	fmt.Println("[RESEND-LOG] ----- Process Completed (Success) -----")
	return nil
}
