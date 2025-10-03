package services

import "github.com/hariomop12/clearrouter/apps/backend/internal/models"

func GetAllModels() []models.ModelDefinition {
	all := []models.ModelDefinition{}
	all = append(all, models.OpenAIModels...)
	all = append(all, models.GoogleModels...)
	all = append(all, models.AnthropicModels...)
	all = append(all, models.MistralModels...)
	all = append(all, models.DeepSeekModels...)
	return all
}
