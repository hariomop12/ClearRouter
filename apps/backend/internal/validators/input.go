package validators

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain uppercase, lowercase, and numeric characters")
	}

	return nil
}

// SanitizeInput removes potentially harmful characters
func SanitizeInput(input string) string {
	// Remove HTML tags and trim whitespace
	input = strings.TrimSpace(input)
	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")
	return input
}

// ValidateChatMessage validates chat message content
func ValidateChatMessage(content string) error {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return errors.New("message content cannot be empty")
	}
	if len(content) > 10000 {
		return errors.New("message content exceeds maximum length of 10,000 characters")
	}
	return nil
}
