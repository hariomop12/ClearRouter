package models

var OpenAIModels = []ModelDefinition{
	// GPT-4o Series
	{
		ID:         "gpt-4o-mini",
		Name:       "GPT-4o Mini",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4o-mini",
				InputPrice:  0.15 / 1e6, // $0.15 per 1M tokens
				OutputPrice: 0.6 / 1e6,  // $0.6 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   16384,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4o",
		Name:       "GPT-4o",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4o",
				InputPrice:  2.5 / 1e6,  // $2.5 per 1M tokens
				OutputPrice: 10.0 / 1e6, // $10.0 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   16384,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	// GPT-4 Series
	{
		ID:         "gpt-4",
		Name:       "GPT-4",
		Family:     "openai",
		JSONOutput: false,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4",
				InputPrice:  30.0 / 1e6, // $30.0 per 1M tokens
				OutputPrice: 60.0 / 1e6, // $60.0 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4-turbo",
		Name:       "GPT-4 Turbo",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4-turbo",
				InputPrice:  10.0 / 1e6, // $10.0 per 1M tokens
				OutputPrice: 30.0 / 1e6, // $30.0 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4.1",
		Name:       "GPT-4.1",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4.1",
				InputPrice:  2.0 / 1e6, // $2.0 per 1M tokens
				OutputPrice: 8.0 / 1e6, // $8.0 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4.1-mini",
		Name:       "GPT-4.1 Mini",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4.1-mini",
				InputPrice:  0.4 / 1e6, // $0.4 per 1M tokens
				OutputPrice: 1.6 / 1e6, // $1.6 per 1M tokens
				ContextSize: 1000000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-4.1-nano",
		Name:       "GPT-4.1 Nano",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-4.1-nano",
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
	// GPT-3.5 Series
	{
		ID:         "gpt-3.5-turbo",
		Name:       "GPT-3.5 Turbo",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-3.5-turbo",
				InputPrice:  0.5 / 1e6, // $0.5 per 1M tokens
				OutputPrice: 1.5 / 1e6, // $1.5 per 1M tokens
				ContextSize: 16385,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// o1 Series (Reasoning Models)
	{
		ID:         "o1",
		Name:       "o1",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "o1",
				InputPrice:  15.0 / 1e6, // $15.0 per 1M tokens
				OutputPrice: 60.0 / 1e6, // $60.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      true,
				Tools:       false,
			},
		},
	},
	{
		ID:         "o1-mini",
		Name:       "o1 Mini",
		Family:     "openai",
		JSONOutput: false,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "o1-mini",
				InputPrice:  1.1 / 1e6, // $1.1 per 1M tokens
				OutputPrice: 4.4 / 1e6, // $4.4 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   8192,
				Streaming:   false,
				Vision:      false,
				Tools:       false,
			},
		},
	},
	{
		ID:         "o3",
		Name:       "o3",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "o3",
				InputPrice:  2.0 / 1e6, // $2.0 per 1M tokens
				OutputPrice: 8.0 / 1e6, // $8.0 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   false,
				Vision:      true,
				Tools:       false,
			},
		},
	},
	{
		ID:         "o3-mini",
		Name:       "o3 Mini",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "o3-mini",
				InputPrice:  1.1 / 1e6, // $1.1 per 1M tokens
				OutputPrice: 4.4 / 1e6, // $4.4 per 1M tokens
				ContextSize: 200000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       false,
			},
		},
	},
	// GPT-5 Series
	{
		ID:         "gpt-5",
		Name:       "GPT-5",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-5",
				InputPrice:  1.25 / 1e6, // $1.25 per 1M tokens
				OutputPrice: 10.0 / 1e6, // $10.0 per 1M tokens
				ContextSize: 400000,
				MaxOutput:   128000,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-5-mini",
		Name:       "GPT-5 Mini",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-5-mini",
				InputPrice:  0.25 / 1e6, // $0.25 per 1M tokens
				OutputPrice: 2.0 / 1e6,  // $2.0 per 1M tokens
				ContextSize: 400000,
				MaxOutput:   128000,
				Streaming:   true,
				Vision:      true,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-5-nano",
		Name:       "GPT-5 Nano",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-5-nano",
				InputPrice:  0.05 / 1e6, // $0.05 per 1M tokens
				OutputPrice: 0.4 / 1e6,  // $0.4 per 1M tokens
				ContextSize: 400000,
				MaxOutput:   128000,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-5-chat-latest",
		Name:       "GPT-5 Chat Latest",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "gpt-5-chat-latest",
				InputPrice:  1.25 / 1e6, // $1.25 per 1M tokens
				OutputPrice: 10.0 / 1e6, // $10.0 per 1M tokens
				ContextSize: 400000,
				MaxOutput:   128000,
				Streaming:   true,
				Vision:      true,
				Tools:       false,
			},
		},
	},
	// GPT OSS Series
	{
		ID:         "gpt-oss-120b",
		Name:       "GPT OSS 120B",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "groq",
				ModelName:   "openai/gpt-oss-120b",
				InputPrice:  0.15 / 1e6, // $0.15 per 1M tokens
				OutputPrice: 0.75 / 1e6, // $0.75 per 1M tokens
				ContextSize: 131072,
				MaxOutput:   32766,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	{
		ID:         "gpt-oss-20b",
		Name:       "GPT OSS 20B",
		Family:     "openai",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "openai",
				ModelName:   "openai/gpt-oss-20b",
				InputPrice:  0.1 / 1e6, // $0.1 per 1M tokens
				OutputPrice: 0.5 / 1e6, // $0.5 per 1M tokens
				ContextSize: 131072,
				MaxOutput:   32766,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
}
