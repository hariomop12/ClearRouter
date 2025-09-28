package models

// GetAllModels returns all available models from different providers
func GetAllModels() []ModelDefinition {
	var all []ModelDefinition
	all = append(all, OpenAIModels...)
	all = append(all, GoogleModels...)
	return all
}
