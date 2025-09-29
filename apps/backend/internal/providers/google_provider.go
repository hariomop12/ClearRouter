package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
)

type GoogleProvider struct {
	apiKey     string
	definition ProviderDefinition
	httpClient *http.Client
}

func NewGoogleProvider() *GoogleProvider {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		// Return a provider but log the error - don't panic on startup
		fmt.Printf("[WARNING] GOOGLE_API_KEY environment variable is not set\n")
	}

	return &GoogleProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		definition: ProviderDefinition{
			ID:           "google",
			Name:         "Google AI",
			Description:  "Google's Gemini models with support for text and multimodal tasks",
			Streaming:    true,
			Cancellation: true,
			JSONOutput:   true,
		},
	}
}

func (p *GoogleProvider) GetName() string {
	return p.definition.Name
}

func (p *GoogleProvider) GetDefinition() ProviderDefinition {
	return p.definition
}

func (p *GoogleProvider) GetSupportedModels() []string {
	return []string{
		// Gemini 1.5 Series
		"gemini-1.5-flash",
		"gemini-1.5-flash-8b",
		"gemini-1.5-pro",
		// Gemini 2.0 Series
		"gemini-2.0-flash",
		"gemini-2.0-flash-lite",
		// Gemini 2.5 Series
		"gemini-2.5-pro",
		"gemini-2.5-pro-preview-05-06",
		"gemini-2.5-pro-preview-06-05",
		"gemini-2.5-flash",
		"gemini-2.5-flash-lite",
		"gemini-2.5-flash-preview-04-17",
		"gemini-2.5-flash-preview-05-20",
		"gemini-2.5-flash-preview-09-2025",
		"gemini-2.5-flash-lite-preview-09-2025",
		"gemini-2.5-flash-image-preview",
		"gemini-2.5-flash-preview-04-17-thinking",
		// Gemma Series
		"gemma-3n-e2b-it",
		"gemma-3n-e4b-it",
		"gemma-3-1b-it",
		"gemma-3-4b-it",
		"gemma-3-12b-it",
	}
}

// getActualModelName maps our model IDs to the actual model names used by Google API
func (p *GoogleProvider) getActualModelName(modelID string) string {
	modelMap := map[string]string{
		// Gemini 1.5 Series
		"gemini-1.5-flash":    "gemini-1.5-flash",
		"gemini-1.5-flash-8b": "gemini-1.5-flash-8b",
		"gemini-1.5-pro":      "gemini-1.5-pro",
		// Gemini 2.0 Series
		"gemini-2.0-flash":      "gemini-2.0-flash",
		"gemini-2.0-flash-lite": "gemini-2.0-flash-lite",
		// Gemini 2.5 Series
		"gemini-2.5-pro":                              "gemini-2.5-pro",
		"gemini-2.5-pro-preview-05-06":               "gemini-2.5-pro-preview-05-06",
		"gemini-2.5-pro-preview-06-05":               "gemini-2.5-pro-preview-06-05",
		"gemini-2.5-flash":                            "gemini-2.5-flash",
		"gemini-2.5-flash-lite":                       "gemini-2.5-flash-lite",
		"gemini-2.5-flash-preview-04-17":             "gemini-2.5-flash-preview-04-17",
		"gemini-2.5-flash-preview-05-20":             "gemini-2.5-flash-preview-05-20",
		"gemini-2.5-flash-preview-09-2025":           "gemini-2.5-flash-preview-09-2025",
		"gemini-2.5-flash-lite-preview-09-2025":      "gemini-2.5-flash-lite-preview-09-2025",
		"gemini-2.5-flash-image-preview":             "gemini-2.5-flash-image-preview",
		"gemini-2.5-flash-preview-04-17-thinking":    "gemini-2.5-flash-preview-04-17-thinking",
		// Gemma Series
		"gemma-3n-e2b-it":  "gemma-3n-e2b-it",
		"gemma-3n-e4b-it":  "gemma-3n-e4b-it",
		"gemma-3-1b-it":    "gemma-3-1b-it",
		"gemma-3-4b-it":    "gemma-3-4b-it",
		"gemma-3-12b-it":   "gemma-3-12b-it",
	}

	if actualName, exists := modelMap[modelID]; exists {
		return actualName
	}

	// Fallback to the original model ID if not found in map
	return modelID
}

// GoogleRequest represents the request format for Google's Gemini API
type GoogleRequest struct {
	Contents []GoogleContent `json:"contents"`
}

type GoogleContent struct {
	Parts []GooglePart `json:"parts"`
}

type GooglePart struct {
	Text string `json:"text"`
}

// GoogleResponse represents the response format from Google's Gemini API
type GoogleResponse struct {
	Candidates    []GoogleCandidate   `json:"candidates"`
	UsageMetadata GoogleUsageMetadata `json:"usageMetadata"`
}

type GoogleCandidate struct {
	Content GoogleResponseContent `json:"content"`
}

type GoogleResponseContent struct {
	Parts []GoogleResponsePart `json:"parts"`
}

type GoogleResponsePart struct {
	Text string `json:"text"`
}

type GoogleUsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

func (p *GoogleProvider) CreateChatCompletion(ctx context.Context, req *models.ChatCompletionsRequest) (*models.ChatCompletionsResponse, error) {
	// Check if API key is available
	if p.apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY environment variable is not set")
	}

	// Get the actual model name to use with Google API
	actualModelName := p.getActualModelName(req.Model)

	// Convert messages to Google format
	var contents []GoogleContent
	for _, msg := range req.Messages {
		content := GoogleContent{
			Parts: []GooglePart{
				{Text: msg.Content},
			},
		}
		contents = append(contents, content)
	}

	// Create request payload
	googleReq := GoogleRequest{
		Contents: contents,
	}

	// Marshal request to JSON
	requestBody, err := json.Marshal(googleReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Build API URL
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", actualModelName)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", p.apiKey)

	// Send request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var googleResp GoogleResponse
	if err := json.Unmarshal(responseBody, &googleResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Extract response text
	var responseText string
	if len(googleResp.Candidates) > 0 && len(googleResp.Candidates[0].Content.Parts) > 0 {
		for _, part := range googleResp.Candidates[0].Content.Parts {
			responseText += part.Text
		}
	}

	if responseText == "" {
		return nil, fmt.Errorf("no text response received from model")
	}

	// Create OpenAI-compatible response
	return &models.ChatCompletionsResponse{
		ID:      "chatcmpl-" + fmt.Sprintf("%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []struct {
			Index   int                 `json:"index"`
			Message *models.ChatMessage `json:"message,omitempty"`
			Delta   *models.ChatMessage `json:"delta,omitempty"`
		}{
			{
				Index: 0,
				Message: &models.ChatMessage{
					Role:    "assistant",
					Content: responseText,
				},
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     googleResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: googleResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      googleResp.UsageMetadata.TotalTokenCount,
		},
	}, nil
}

func (p *GoogleProvider) CalculateTokens(messages []models.ChatMessage) (int, error) {
	// Simple approximation - roughly 4 characters per token
	totalChars := 0
	for _, msg := range messages {
		totalChars += len(msg.Content)
	}
	return totalChars / 4, nil
}
