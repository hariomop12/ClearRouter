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

type MistralProvider struct {
	apiKey     string
	httpClient *http.Client
	definition ProviderDefinition
}

// Mistral API request/response structures
type MistralRequest struct {
	Model     string           `json:"model"`
	Messages  []MistralMessage `json:"messages"`
	Stream    bool             `json:"stream"`
	MaxTokens *int             `json:"max_tokens,omitempty"`
}

type MistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MistralResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int             `json:"index"`
		Message *MistralMessage `json:"message,omitempty"`
		Delta   *MistralMessage `json:"delta,omitempty"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewMistralProvider() *MistralProvider {
	return &MistralProvider{
		apiKey: os.Getenv("MISTRAL_API_KEY"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		definition: ProviderDefinition{
			ID:           "mistral",
			Name:         "Mistral AI",
			Description:  "Mistral's open-source and commercial AI models",
			Streaming:    true,
			Cancellation: true,
			JSONOutput:   true,
			BaseURL:      "https://api.mistral.ai",
		},
	}
}

func (p *MistralProvider) GetName() string {
	return p.definition.Name
}

func (p *MistralProvider) GetDefinition() ProviderDefinition {
	return p.definition
}

func (p *MistralProvider) GetSupportedModels() []string {
	return []string{
		"mistral-large-latest",
		"mixtral-8x7b-instruct",
		"mistral-7b-instruct",
		"pixtral-large-latest",
		"mistral-medium",
		"mistral-small",
		"mixtral-8x22b-instruct",
	}
}

func (p *MistralProvider) CreateChatCompletion(ctx context.Context, req *models.ChatCompletionsRequest) (*models.ChatCompletionsResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("MISTRAL_API_KEY environment variable is not set")
	}

	// Convert our request to Mistral format
	mistralReq := MistralRequest{
		Model:     req.Model,
		Stream:    req.Stream,
		MaxTokens: req.MaxTokens,
		Messages:  make([]MistralMessage, len(req.Messages)),
	}

	for i, msg := range req.Messages {
		mistralReq.Messages[i] = MistralMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Marshal request
	jsonData, err := json.Marshal(mistralReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
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
		return nil, fmt.Errorf("mistral API returned status %d", resp.StatusCode)
	}

	// Parse response
	var mistralResp MistralResponse
	if err := json.NewDecoder(resp.Body).Decode(&mistralResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our standard response format
	response := &models.ChatCompletionsResponse{
		ID:      mistralResp.ID,
		Object:  mistralResp.Object,
		Created: mistralResp.Created,
		Model:   mistralResp.Model,
		Choices: make([]struct {
			Index   int                 `json:"index"`
			Message *models.ChatMessage `json:"message,omitempty"`
			Delta   *models.ChatMessage `json:"delta,omitempty"`
		}, len(mistralResp.Choices)),
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     mistralResp.Usage.PromptTokens,
			CompletionTokens: mistralResp.Usage.CompletionTokens,
			TotalTokens:      mistralResp.Usage.TotalTokens,
		},
	}

	for i, choice := range mistralResp.Choices {
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

func (p *MistralProvider) CalculateTokens(messages []models.ChatMessage) (int, error) {
	// TODO: Implement proper token counting for Mistral
	// For now, use a simple approximation
	total := 0
	for _, msg := range messages {
		total += len(strings.Split(msg.Content, " ")) * 4 // Rough approximation
	}
	return total, nil
}
