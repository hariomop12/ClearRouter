package models

var GoogleModels = []ModelDefinition{
	// Gemini 1.5 Series
	{
		ID:         "gemini-1.5-flash",
		Name:       "Gemini 1.5 Flash",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-1.5-flash",
				InputPrice:  0.0375 / 1e6, // $0.0375 per 1M tokens
				OutputPrice: 0.15 / 1e6,   // $0.15 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-1.5-flash-8b",
		Name:       "Gemini 1.5 Flash 8B",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-1.5-flash-8b",
				InputPrice:  0.0375 / 1e6, // $0.0375 per 1M tokens
				OutputPrice: 0.15 / 1e6,   // $0.15 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
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
	// Gemini 2.0 Series
	{
		ID:         "gemini-2.0-flash",
		Name:       "Gemini 2.0 Flash",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.0-flash",
				InputPrice:  0.1 / 1e6, // $0.1 per 1M tokens
				OutputPrice: 0.4 / 1e6, // $0.4 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   false,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.0-flash-lite",
		Name:       "Gemini 2.0 Flash Lite",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.0-flash-lite",
				InputPrice:  0.075 / 1e6, // $0.075 per 1M tokens
				OutputPrice: 0.3 / 1e6,   // $0.3 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Gemini 2.5 Series
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
	{
		ID:         "gemini-2.5-pro-preview-05-06",
		Name:       "Gemini 2.5 Pro Preview (05-06)",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-pro-preview-05-06",
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
	{
		ID:         "gemini-2.5-pro-preview-06-05",
		Name:       "Gemini 2.5 Pro Preview (06-05)",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-pro-preview-06-05",
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
	{
		ID:         "gemini-2.5-flash",
		Name:       "Gemini 2.5 Flash",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash",
				InputPrice:  0.3 / 1e6,  // $0.3 per 1M tokens
				OutputPrice: 2.5 / 1e6,  // $2.5 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash-lite",
		Name:       "Gemini 2.5 Flash Lite",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash-lite",
				InputPrice:  0.1 / 1e6, // $0.1 per 1M tokens
				OutputPrice: 0.4 / 1e6, // $0.4 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash-preview-04-17",
		Name:       "Gemini 2.5 Flash Preview (04-17)",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash-preview-04-17",
				InputPrice:  0.15 / 1e6, // $0.15 per 1M tokens
				OutputPrice: 0.6 / 1e6,  // $0.6 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash-preview-05-20",
		Name:       "Gemini 2.5 Flash Preview (05-20)",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash-preview-05-20",
				InputPrice:  0.15 / 1e6, // $0.15 per 1M tokens
				OutputPrice: 0.6 / 1e6,  // $0.6 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash-preview-09-2025",
		Name:       "Gemini 2.5 Flash Preview (09-2025)",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash-preview-09-2025",
				InputPrice:  0.3 / 1e6,  // $0.3 per 1M tokens
				OutputPrice: 2.5 / 1e6,  // $2.5 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash-lite-preview-09-2025",
		Name:       "Gemini 2.5 Flash Lite Preview (09-2025)",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash-lite-preview-09-2025",
				InputPrice:  0.1 / 1e6, // $0.1 per 1M tokens
				OutputPrice: 0.4 / 1e6, // $0.4 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash-image-preview",
		Name:       "Gemini 2.5 Flash Image Preview",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash-image-preview",
				InputPrice:  0.3 / 1e6, // $0.3 per 1M tokens
				OutputPrice: 30 / 1e6,  // $30 per 1M tokens
				ContextSize: 32800,
				MaxOutput:   8200,
				Streaming:   true,
				Vision:      true,
				Tools:       false,
			},
		},
	},
	{
		ID:         "gemini-2.5-flash-preview-04-17-thinking",
		Name:       "Gemini 2.5 Flash Preview Thinking (04-17)",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemini-2.5-flash-preview-04-17-thinking",
				InputPrice:  0.15 / 1e6, // $0.15 per 1M tokens
				OutputPrice: 0.6 / 1e6,  // $0.6 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	// Gemma Series
	{
		ID:         "gemma-3n-e2b-it",
		Name:       "Gemma 3n E2B IT",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemma-3n-e2b-it",
				InputPrice:  0.075 / 1e6, // $0.075 per 1M tokens
				OutputPrice: 0.3 / 1e6,   // $0.3 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       false,
			},
		},
	},
	{
		ID:         "gemma-3n-e4b-it",
		Name:       "Gemma 3n E4B IT",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemma-3n-e4b-it",
				InputPrice:  0.075 / 1e6, // $0.075 per 1M tokens
				OutputPrice: 0.3 / 1e6,   // $0.3 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       false,
			},
		},
	},
	{
		ID:         "gemma-3-1b-it",
		Name:       "Gemma 3 1B IT",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemma-3-1b-it",
				InputPrice:  0.075 / 1e6, // $0.075 per 1M tokens
				OutputPrice: 0.3 / 1e6,   // $0.3 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       false,
			},
		},
	},
	{
		ID:         "gemma-3-4b-it",
		Name:       "Gemma 3 4B IT",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemma-3-4b-it",
				InputPrice:  0.075 / 1e6, // $0.075 per 1M tokens
				OutputPrice: 0.3 / 1e6,   // $0.3 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       false,
			},
		},
	},
	{
		ID:         "gemma-3-12b-it",
		Name:       "Gemma 3 12B IT",
		Family:     "google",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "google",
				ModelName:   "gemma-3-12b-it",
				InputPrice:  0.075 / 1e6, // $0.075 per 1M tokens
				OutputPrice: 0.3 / 1e6,   // $0.3 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       false,
			},
		},
	},
}
