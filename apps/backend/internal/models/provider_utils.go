package models

// GetProviderFromModel returns the provider ID based on the model name
func GetProviderFromModel(modelName string) string {
	allModels := GetAllModels()
	for _, model := range allModels {
		if model.ID == modelName {
			if len(model.Providers) > 0 {
				return model.Providers[0].ProviderID
			}
		}
	}
	return ""
}
