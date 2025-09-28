package models

// ProviderModel defines model-specific settings for a provider
type ProviderModel struct {
	ProviderID  string  `json:"provider_id"`
	ModelName   string  `json:"model_name"`
	InputPrice  float64 `json:"input_price,omitempty"`
	OutputPrice float64 `json:"output_price,omitempty"`
	ContextSize int     `json:"context_size,omitempty"`
	MaxOutput   int     `json:"max_output,omitempty"`
	Streaming   bool    `json:"streaming"`
	Vision      bool    `json:"vision,omitempty"`
	Tools       bool    `json:"tools,omitempty"`
}

// ModelDefinition defines a model and its supported providers
type ModelDefinition struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Family     string          `json:"family"`
	Providers  []ProviderModel `json:"providers"`
	Status     string          `json:"status"`
	JSONOutput bool            `json:"json_output"`
}

// Available models configuration
var AvailableModels = []ModelDefinition{
	{
		ID:     "gemini-pro",
		Name:   "Gemini Pro",
		Family: "google",
		Status: "active",
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-pro",
				InputPrice:  0.000001,
				OutputPrice: 0.000002,
				ContextSize: 32768,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:     "gemini-pro-vision",
		Name:   "Gemini Pro Vision",
		Family: "google",
		Status: "active",
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-pro-vision",
				InputPrice:  0.000001,
				OutputPrice: 0.000002,
				ContextSize: 32768,
				Streaming:   false,
				Vision:      true,
				Tools:       true,
			},
		},
	},
}

// GetModelByID returns a model definition by its ID
func GetModelByID(modelID string) *ModelDefinition {
	for _, model := range AvailableModels {
		if model.ID == modelID {
			return &model
		}
	}
	return nil
}

// GetProviderForModel returns the first available provider for a model
func GetProviderForModel(modelID string) *ProviderModel {
	model := GetModelByID(modelID)
	if model != nil && len(model.Providers) > 0 {
		return &model.Providers[0]
	}
	return nil
}
