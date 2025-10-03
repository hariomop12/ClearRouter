package models

var MistralModels = []ModelDefinition{
	// Mistral Large Latest
	{
		ID:         "mistral-large-latest",
		Name:       "Mistral Large Latest",
		Family:     "mistral",
		JSONOutput: false,
		Providers: []ProviderModel{
			{
				ProviderID:  "mistral",
				ModelName:   "mistral-large-latest",
				InputPrice:  4.0 / 1e6,  // $4.0 per 1M tokens
				OutputPrice: 12.0 / 1e6, // $12.0 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Mixtral 8x7B Instruct
	{
		ID:         "mixtral-8x7b-instruct",
		Name:       "Mixtral 8x7B Instruct",
		Family:     "mistral",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "mistral",
				ModelName:   "mistralai/mixtral-8x7b-instruct-v0.1",
				InputPrice:  0.7 / 1e6, // $0.7 per 1M tokens
				OutputPrice: 0.7 / 1e6, // $0.7 per 1M tokens
				ContextSize: 32768,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Mistral 7B Instruct
	{
		ID:         "mistral-7b-instruct",
		Name:       "Mistral 7B Instruct",
		Family:     "mistral",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "mistral",
				ModelName:   "mistralai/mistral-7b-instruct-v0.1",
				InputPrice:  0.2 / 1e6, // $0.2 per 1M tokens
				OutputPrice: 0.2 / 1e6, // $0.2 per 1M tokens
				ContextSize: 8192,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Pixtral Large Latest (Vision Model)
	{
		ID:         "pixtral-large-latest",
		Name:       "Pixtral Large Latest",
		Family:     "mistral",
		JSONOutput: false,
		Providers: []ProviderModel{
			{
				ProviderID:  "mistral",
				ModelName:   "pixtral-large-latest",
				InputPrice:  4.0 / 1e6,  // $4.0 per 1M tokens
				OutputPrice: 12.0 / 1e6, // $12.0 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	// Mistral Medium
	{
		ID:         "mistral-medium",
		Name:       "Mistral Medium",
		Family:     "mistral",
		JSONOutput: false,
		Providers: []ProviderModel{
			{
				ProviderID:  "mistral",
				ModelName:   "mistral-medium-latest",
				InputPrice:  2.7 / 1e6, // $2.7 per 1M tokens
				OutputPrice: 8.1 / 1e6, // $8.1 per 1M tokens
				ContextSize: 32768,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Mistral Small
	{
		ID:         "mistral-small",
		Name:       "Mistral Small",
		Family:     "mistral",
		JSONOutput: false,
		Providers: []ProviderModel{
			{
				ProviderID:  "mistral",
				ModelName:   "mistral-small-latest",
				InputPrice:  1.0 / 1e6, // $1.0 per 1M tokens
				OutputPrice: 3.0 / 1e6, // $3.0 per 1M tokens
				ContextSize: 32768,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Mixtral 8x22B Instruct
	{
		ID:         "mixtral-8x22b-instruct",
		Name:       "Mixtral 8x22B Instruct",
		Family:     "mistral",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "mistral",
				ModelName:   "mistralai/mixtral-8x22b-instruct-v0.1",
				InputPrice:  2.0 / 1e6, // $2.0 per 1M tokens
				OutputPrice: 6.0 / 1e6, // $6.0 per 1M tokens
				ContextSize: 65536,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
}
