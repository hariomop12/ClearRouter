package services

import "github.com/hariomop12/clearrouter/apps/backend/internal/models"

func GetAllModels() []models.ModelDefinition {
	all := []models.ModelDefinition{}
	all = append(all, models.OpenAIModels...)
	all = append(all, models.GoogleModels...)
	return all
}
