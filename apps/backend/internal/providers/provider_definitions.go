package providers

import (
	"context"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
)

// ProviderDefinition defines capabilities and metadata for a provider
type ProviderDefinition struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Streaming    bool   `json:"streaming"`
	Cancellation bool   `json:"cancellation"`
	JSONOutput   bool   `json:"jsonOutput"`
	BaseURL      string `json:"baseUrl,omitempty"`
}

// Provider interface defines methods that each provider must implement
type Provider interface {
	GetName() string
	GetDefinition() ProviderDefinition
	GetSupportedModels() []string
	CreateChatCompletion(ctx context.Context, req *models.ChatCompletionsRequest) (*models.ChatCompletionsResponse, error)
}

// Supported providers
var Providers = []ProviderDefinition{
	{
		ID:           "google",
		Name:         "Google AI",
		Description:  "Google's Gemini models with support for text and multimodal tasks",
		Streaming:    true,
		Cancellation: true,
		JSONOutput:   true,
	},
	{
		ID:           "openai",
		Name:         "OpenAI",
		Description:  "OpenAI's GPT models and DALL-E for text and image generation",
		Streaming:    true,
		Cancellation: true,
		JSONOutput:   true,
	},
}
