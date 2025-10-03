package models

var DeepSeekModels = []ModelDefinition{
	// DeepSeek Chat
	{
		ID:         "deepseek-chat",
		Name:       "DeepSeek Chat",
		Family:     "deepseek",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "deepseek",
				ModelName:   "deepseek-chat",
				InputPrice:  0.14 / 1e6, // $0.14 per 1M tokens
				OutputPrice: 0.28 / 1e6, // $0.28 per 1M tokens
				ContextSize: 32768,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// DeepSeek Coder
	{
		ID:         "deepseek-coder",
		Name:       "DeepSeek Coder",
		Family:     "deepseek",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "deepseek",
				ModelName:   "deepseek-coder",
				InputPrice:  0.14 / 1e6, // $0.14 per 1M tokens
				OutputPrice: 0.28 / 1e6, // $0.28 per 1M tokens
				ContextSize: 32768,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// DeepSeek V2 Chat
	{
		ID:         "deepseek-v2-chat",
		Name:       "DeepSeek V2 Chat",
		Family:     "deepseek",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "deepseek",
				ModelName:   "deepseek-v2-chat",
				InputPrice:  0.07 / 1e6, // $0.07 per 1M tokens
				OutputPrice: 0.28 / 1e6, // $0.28 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   4096,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// DeepSeek V2.5 Chat
	{
		ID:         "deepseek-v2-5-chat",
		Name:       "DeepSeek V2.5 Chat",
		Family:     "deepseek",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "deepseek",
				ModelName:   "deepseek-v2.5-chat",
				InputPrice:  0.07 / 1e6, // $0.07 per 1M tokens
				OutputPrice: 0.28 / 1e6, // $0.28 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// DeepSeek V3
	{
		ID:         "deepseek-v3",
		Name:       "DeepSeek V3",
		Family:     "deepseek",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "deepseek",
				ModelName:   "deepseek-v3",
				InputPrice:  0.07 / 1e6, // $0.07 per 1M tokens
				OutputPrice: 0.28 / 1e6, // $0.28 per 1M tokens
				ContextSize: 128000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// DeepSeek R1
	{
		ID:         "deepseek-r1",
		Name:       "DeepSeek R1",
		Family:     "deepseek",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "deepseek",
				ModelName:   "deepseek-r1",
				InputPrice:  0.55 / 1e6, // $0.55 per 1M tokens (reasoning model)
				OutputPrice: 2.19 / 1e6, // $2.19 per 1M tokens
				ContextSize: 64000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
	// DeepSeek R1 Distill
	{
		ID:         "deepseek-r1-distill",
		Name:       "DeepSeek R1 Distill",
		Family:     "deepseek",
		JSONOutput: true,
		Providers: []ProviderModel{
			{
				ProviderID:  "deepseek",
				ModelName:   "deepseek-r1-distill",
				InputPrice:  0.14 / 1e6, // $0.14 per 1M tokens
				OutputPrice: 0.28 / 1e6, // $0.28 per 1M tokens
				ContextSize: 64000,
				MaxOutput:   8192,
				Streaming:   true,
				Vision:      false,
				Tools:       true,
			},
		},
	},
}
