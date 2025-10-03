package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
)

type AnthropicProvider struct {
	apiKey     string
	httpClient *http.Client
	definition ProviderDefinition
}

// Anthropic API request/response structures
type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
	Stream    bool               `json:"stream,omitempty"`
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func NewAnthropicProvider() *AnthropicProvider {
	return &AnthropicProvider{
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		definition: ProviderDefinition{
			ID:           "anthropic",
			Name:         "Anthropic",
			Description:  "Anthropic's Claude models for advanced reasoning and conversation",
			Streaming:    true,
			Cancellation: true,
			JSONOutput:   true,
			BaseURL:      "https://api.anthropic.com",
		},
	}
}

func (p *AnthropicProvider) GetName() string {
	return p.definition.Name
}

func (p *AnthropicProvider) GetDefinition() ProviderDefinition {
	return p.definition
}

func (p *AnthropicProvider) GetSupportedModels() []string {
	return []string{
		"claude-3-7-sonnet",
		"claude-3-5-haiku-20241022",
		"claude-3-7-sonnet-20250219",
		"claude-3-5-sonnet-20241022",
		"claude-sonnet-4-20250514",
		"claude-sonnet-4-5",
		"claude-opus-4-20250514",
		"claude-opus-4-1-20250805",
		"claude-3-5-sonnet-20240620",
		"claude-3-5-sonnet",
		"claude-3-5-haiku",
		"claude-3-opus",
		"claude-3-haiku",
	}
}

func (p *AnthropicProvider) CreateChatCompletion(ctx context.Context, req *models.ChatCompletionsRequest) (*models.ChatCompletionsResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	// Set default max tokens if not provided
	maxTokens := 1024
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}

	// Convert our request to Anthropic format
	anthropicReq := AnthropicRequest{
		Model:     req.Model,
		MaxTokens: maxTokens,
		Stream:    req.Stream,
		Messages:  make([]AnthropicMessage, len(req.Messages)),
	}

	for i, msg := range req.Messages {
		anthropicReq.Messages[i] = AnthropicMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Marshal request
	jsonData, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers (following the curl example)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Make request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic API returned status %d", resp.StatusCode)
	}

	// Parse response
	var anthropicResp AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our standard response format
	response := &models.ChatCompletionsResponse{
		ID:      anthropicResp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   anthropicResp.Model,
		Choices: []struct {
			Index   int                 `json:"index"`
			Message *models.ChatMessage `json:"message,omitempty"`
			Delta   *models.ChatMessage `json:"delta,omitempty"`
		}{
			{
				Index: 0,
				Message: &models.ChatMessage{
					Role:    anthropicResp.Role,
					Content: "",
				},
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
	}

	// Combine all content text into a single message
	var contentText string
	for _, content := range anthropicResp.Content {
		if content.Type == "text" {
			contentText += content.Text
		}
	}
	response.Choices[0].Message.Content = contentText

	return response, nil
}

func (p *AnthropicProvider) CalculateTokens(messages []models.ChatMessage) (int, error) {
	// TODO: Implement proper token counting for Anthropic
	// For now, use a simple approximation
	total := 0
	for _, msg := range messages {
		total += len(strings.Split(msg.Content, " ")) * 4 // Rough approximation
	}
	return total, nil
}
