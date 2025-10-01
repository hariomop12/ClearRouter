package services

import (
	"context"
	"fmt"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"github.com/hariomop12/clearrouter/apps/backend/internal/providers"
)

type ChatProvider interface {
	GetName() string
	GetDefinition() providers.ProviderDefinition
	CreateChatCompletion(ctx context.Context, req *models.ChatCompletionsRequest) (*models.ChatCompletionsResponse, error)
	CalculateTokens(messages []models.ChatMessage) (int, error)
}

type ProviderService struct {
	providers map[string]ChatProvider
}

func NewProviderService() *ProviderService {
	return &ProviderService{
		providers: make(map[string]ChatProvider),
	}
}

func (s *ProviderService) RegisterProvider(provider ChatProvider) {
	s.providers[provider.GetDefinition().ID] = provider
}

func (s *ProviderService) GetProvider(name string) (ChatProvider, error) {
	provider, ok := s.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

// GetProviderForModel finds the provider that handles a specific model
func (s *ProviderService) GetProviderForModel(modelName string) (ChatProvider, error) {
	allModels := models.GetAllModels()
	for _, m := range allModels {
		if m.ID == modelName {
			if len(m.Providers) > 0 {
				providerID := m.Providers[0].ProviderID
				return s.GetProvider(providerID)
			}
		}
	}
	return nil, fmt.Errorf("no provider found for model %s", modelName)
}

// GetModelInfo returns the model pricing info
func (s *ProviderService) GetModelInfo(provider, model string) (*models.ProviderModel, error) {
	allModels := models.GetAllModels()
	for _, m := range allModels {
		if m.ID == model {
			for _, p := range m.Providers {
				if p.ProviderID == provider {
					return &p, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("model %s not found for provider %s", model, provider)
}

// CalculateCost calculates the total cost based on token usage
func (s *ProviderService) CalculateCost(model *models.ProviderModel, inputTokens, outputTokens int) (float64, float64) {
	inputCost := float64(inputTokens) * model.InputPrice
	outputCost := float64(outputTokens) * model.OutputPrice
	return inputCost, outputCost
}
