package models

var OpenAIModels = []ModelDefinition{
	{
		ID:         "gpt-3.5-turbo",
		Name:       "GPT-3.5 Turbo",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-3.5-turbo",
				InputPrice:  0.5 / 1e6, // $0.0005 per 1K tokens
				OutputPrice: 1.5 / 1e6, // $0.0015 per 1K tokens
				ContextSize: 16385,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4",
		Name:       "GPT-4",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4",
				InputPrice:  30.0 / 1e6, // $0.03 per 1K tokens
				OutputPrice: 60.0 / 1e6, // $0.06 per 1K tokens
				ContextSize: 8192,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4-turbo-preview",
		Name:       "GPT-4 Turbo",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4-turbo-preview",
				InputPrice:  10.0 / 1e6, // $0.01 per 1K tokens
				OutputPrice: 30.0 / 1e6, // $0.03 per 1K tokens
				ContextSize: 128000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4-vision-preview",
		Name:       "GPT-4 Vision",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4-vision-preview",
				InputPrice:  10.0 / 1e6, // $0.01 per 1K tokens
				OutputPrice: 30.0 / 1e6, // $0.03 per 1K tokens
				ContextSize: 128000,
				MaxOutput:   4096,
				Streaming:   false,
				Vision:      true,
				Tools:       true,
			},
		},
	},
}
