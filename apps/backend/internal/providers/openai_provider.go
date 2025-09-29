package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
)

type OpenAIProvider struct {
	apiKey     string
	httpClient *http.Client
	definition ProviderDefinition
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: os.Getenv("OPENAI_API_KEY"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		definition: ProviderDefinition{
			ID:           "openai",
			Name:         "OpenAI",
			Description:  "OpenAI's GPT models for text generation and chat",
			Streaming:    true,
			Cancellation: true,
			JSONOutput:   true,
		},
	}
}

func (p *OpenAIProvider) GetName() string {
	return p.definition.Name
}

func (p *OpenAIProvider) GetDefinition() ProviderDefinition {
	return p.definition
}

func (p *OpenAIProvider) GetSupportedModels() []string {
	return []string{
		// GPT-4o Series
		"gpt-4o-mini",
		"gpt-4o",
		// GPT-4 Series
		"gpt-4",
		"gpt-4-turbo",
		"gpt-4.1",
		"gpt-4.1-mini",
		"gpt-4.1-nano",
		// GPT-3.5 Series
		"gpt-3.5-turbo",
		// o1 Series (Reasoning Models)
		"o1",
		"o1-mini",
		"o3",
		"o3-mini",
		// GPT-5 Series
		"gpt-5",
		"gpt-5-mini",
		"gpt-5-nano",
		"gpt-5-chat-latest",
		// GPT OSS Series
		"gpt-oss-120b",
		"gpt-oss-20b",
	}
}

func (p *OpenAIProvider) CreateChatCompletion(ctx context.Context, req *models.ChatCompletionsRequest) (*models.ChatCompletionsResponse, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Convert our request to OpenAI format
	openaiReq := struct {
		Model     string               `json:"model"`
		Messages  []models.ChatMessage `json:"messages"`
		MaxTokens *int                 `json:"max_tokens,omitempty"`
		Stream    bool                 `json:"stream"`
	}{
		Model:     req.Model,
		Messages:  req.Messages,
		MaxTokens: req.MaxTokens,
		Stream:    req.Stream,
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("error response from OpenAI (status %d)", resp.StatusCode)
		}
		return nil, fmt.Errorf("error from OpenAI: %s", errResp.Error.Message)
	}

	var openaiResp models.ChatCompletionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &openaiResp, nil
}

func (p *OpenAIProvider) CalculateTokens(messages []models.ChatMessage) (int, error) {
	// TODO: Implement proper token counting using tiktoken
	// For now, use a simple approximation
	total := 0
	for _, msg := range messages {
		total += len(strings.Split(msg.Content, " ")) * 4 // Rough approximation
	}
	return total, nil
}
