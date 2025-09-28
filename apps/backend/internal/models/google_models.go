package models

var GoogleModels = []ModelDefinition{
	{
		ID:         "gemini-1.5-flash",
		Name:       "Gemini 1.5 Flash",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-1.5-flash",
				InputPrice:  0.04 / 1e6, // $0.04 per 1M tokens
				OutputPrice: 0.15 / 1e6, // $0.15 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-1.5-pro",
		Name:       "Gemini 1.5 Pro",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-1.5-pro",
				InputPrice:  2.5 / 1e6,  // $2.50 per 1M tokens
				OutputPrice: 10.0 / 1e6, // $10.00 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.0-flash",
		Name:       "Gemini 2.0 Flash",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.0-flash",
				InputPrice:  0.10 / 1e6, // $0.10 per 1M tokens
				OutputPrice: 0.40 / 1e6, // $0.40 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash",
		Name:       "Gemini 2.5 Flash",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash",
				InputPrice:  0.30 / 1e6, // $0.30 per 1M tokens
				OutputPrice: 2.50 / 1e6, // $2.50 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-pro",
		Name:       "Gemini 2.5 Pro",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-pro",
				InputPrice:  1.25 / 1e6, // $1.25 per 1M tokens
				OutputPrice: 10.0 / 1e6, // $10.00 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
}
