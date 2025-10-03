package models

var AnthropicModels = []ModelDefinition{
	// Claude 3.7 Sonnet
	{
		ID:         "claude-3-7-sonnet",
		Name:       "Claude 3.7 Sonnet",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-7-sonnet-latest",
				InputPrice:  3.0 / 1e6,  // $3.0 per 1M tokens
				OutputPrice: 15.0 / 1e6, // $15.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude 3.5 Haiku (2024-10-22)
	{
		ID:         "claude-3-5-haiku-20241022",
		Name:       "Claude 3.5 Haiku (2024-10-22)",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-5-haiku-20241022",
				InputPrice:  0.8 / 1e6, // $0.8 per 1M tokens
				OutputPrice: 4.0 / 1e6, // $4.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude 3.7 Sonnet (2025-02-19)
	{
		ID:         "claude-3-7-sonnet-20250219",
		Name:       "Claude 3.7 Sonnet (2025-02-19)",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-7-sonnet-20250219",
				InputPrice:  3.0 / 1e6,  // $3.0 per 1M tokens
				OutputPrice: 15.0 / 1e6, // $15.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude 3.5 Sonnet (2024-10-22)
	{
		ID:         "claude-3-5-sonnet-20241022",
		Name:       "Claude 3.5 Sonnet (2024-10-22)",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-5-sonnet-20241022",
				InputPrice:  3.0 / 1e6,  // $3.0 per 1M tokens
				OutputPrice: 15.0 / 1e6, // $15.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude Sonnet 4 (2025-05-14)
	{
		ID:         "claude-sonnet-4-20250514",
		Name:       "Claude Sonnet 4 (2025-05-14)",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-sonnet-4-20250514",
				InputPrice:  3.0 / 1e6,  // $3.0 per 1M tokens
				OutputPrice: 15.0 / 1e6, // $15.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude Sonnet 4.5
	{
		ID:         "claude-sonnet-4-5",
		Name:       "Claude Sonnet 4.5",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-sonnet-4-5",
				InputPrice:  3.0 / 1e6,  // $3.0 per 1M tokens
				OutputPrice: 15.0 / 1e6, // $15.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude Opus 4 (2025-05-14)
	{
		ID:         "claude-opus-4-20250514",
		Name:       "Claude Opus 4 (2025-05-14)",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-opus-4-20250514",
				InputPrice:  15.0 / 1e6, // $15.0 per 1M tokens
				OutputPrice: 75.0 / 1e6, // $75.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude Opus 4.1
	{
		ID:         "claude-opus-4-1-20250805",
		Name:       "Claude Opus 4.1",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-opus-4-1-20250805",
				InputPrice:  15.0 / 1e6, // $15.0 per 1M tokens
				OutputPrice: 75.0 / 1e6, // $75.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   32000,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	// Claude 3.5 Sonnet (Old)
	{
		ID:         "claude-3-5-sonnet-20240620",
		Name:       "Claude 3.5 Sonnet (Old)",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-5-sonnet-20240620",
				InputPrice:  3.0 / 1e6,  // $3.0 per 1M tokens
				OutputPrice: 15.0 / 1e6, // $15.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	// Claude 3.5 Sonnet Latest
	{
		ID:         "claude-3-5-sonnet",
		Name:       "Claude 3.5 Sonnet",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-5-sonnet-latest",
				InputPrice:  3.0 / 1e6,  // $3.0 per 1M tokens
				OutputPrice: 15.0 / 1e6, // $15.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude 3.5 Haiku Latest
	{
		ID:         "claude-3-5-haiku",
		Name:       "Claude 3.5 Haiku",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-5-haiku-latest",
				InputPrice:  0.8 / 1e6, // $0.8 per 1M tokens
				OutputPrice: 4.0 / 1e6, // $4.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// Claude 3 Opus
	{
		ID:         "claude-3-opus",
		Name:       "Claude 3 Opus",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-opus-20240229",
				InputPrice:  15.0 / 1e6, // $15.0 per 1M tokens
				OutputPrice: 75.0 / 1e6, // $75.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	// Claude 3 Haiku
	{
		ID:         "claude-3-haiku",
		Name:       "Claude 3 Haiku",
		Family:     "anthropic",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "anthropic",
				ModelName:   "claude-3-haiku-20240307",
				InputPrice:  0.25 / 1e6, // $0.25 per 1M tokens
				OutputPrice: 1.25 / 1e6, // $1.25 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
}
