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

type DeepSeekProvider struct {
	apiKey     string
	httpClient *http.Client
	definition ProviderDefinition
}

// DeepSeek API request/response structures
type DeepSeekRequest struct {
	Model     string            `json:"model"`
	Messages  []DeepSeekMessage `json:"messages"`
	Stream    bool              `json:"stream"`
	MaxTokens *int              `json:"max_tokens,omitempty"`
}

type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int              `json:"index"`
		Message *DeepSeekMessage `json:"message,omitempty"`
		Delta   *DeepSeekMessage `json:"delta,omitempty"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewDeepSeekProvider() *DeepSeekProvider {
	return &DeepSeekProvider{
		apiKey: os.Getenv("DEEPSEEK_API_KEY"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		definition: ProviderDefinition{
			ID:           "deepseek",
			Name:         "DeepSeek",
			Description:  "DeepSeek's advanced AI models for code and chat",
			Streaming:    true,
			Cancellation: true,
			JSONOutput:   true,
			BaseURL:      "https://api.deepseek.com",
		},
	}
}

func (p *DeepSeekProvider) GetName() string {
	return p.definition.Name
}

func (p *DeepSeekProvider) GetDefinition() ProviderDefinition {
	return p.definition
}

func (p *DeepSeekProvider) GetSupportedModels() []string {
	return []string{
		"deepseek-chat",
		"deepseek-coder",
		"deepseek-v2-chat",
		"deepseek-v2.5-chat",
		"deepseek-v3",
		"deepseek-r1",
		"deepseek-r1-distill",
	}
}

func (p *DeepSeekProvider) CreateChatCompletion(ctx context.Context, req *models.ChatCompletionsRequest) (*models.ChatCompletionsResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("DEEPSEEK_API_KEY environment variable is not set")
	}

	// Convert our request to DeepSeek format
	deepSeekReq := DeepSeekRequest{
		Model:     req.Model,
		Stream:    req.Stream,
		MaxTokens: req.MaxTokens,
		Messages:  make([]DeepSeekMessage, len(req.Messages)),
	}

	for i, msg := range req.Messages {
		deepSeekReq.Messages[i] = DeepSeekMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Marshal request
	jsonData, err := json.Marshal(deepSeekReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Make request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deepSeek API returned status %d", resp.StatusCode)
	}

	// Parse response
	var deepSeekResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&deepSeekResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our standard response format
	response := &models.ChatCompletionsResponse{
		ID:      deepSeekResp.ID,
		Object:  deepSeekResp.Object,
		Created: deepSeekResp.Created,
		Model:   deepSeekResp.Model,
		Choices: make([]struct {
			Index   int                 `json:"index"`
			Message *models.ChatMessage `json:"message,omitempty"`
			Delta   *models.ChatMessage `json:"delta,omitempty"`
		}, len(deepSeekResp.Choices)),
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     deepSeekResp.Usage.PromptTokens,
			CompletionTokens: deepSeekResp.Usage.CompletionTokens,
			TotalTokens:      deepSeekResp.Usage.TotalTokens,
		},
	}

	for i, choice := range deepSeekResp.Choices {
		response.Choices[i].Index = choice.Index
		if choice.Message != nil {
			response.Choices[i].Message = &models.ChatMessage{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			}
		}
		if choice.Delta != nil {
			response.Choices[i].Delta = &models.ChatMessage{
				Role:    choice.Delta.Role,
				Content: choice.Delta.Content,
			}
		}
	}

	return response, nil
}

func (p *DeepSeekProvider) CalculateTokens(messages []models.ChatMessage) (int, error) {
	// TODO: Implement proper token counting for DeepSeek
	// For now, use a simple approximation
	total := 0
	for _, msg := range messages {
		total += len(strings.Split(msg.Content, " ")) * 4 // Rough approximation
	}
	return total, nil
}
