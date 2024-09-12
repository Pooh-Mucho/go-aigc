package chat

import "github.com/Pooh-Mucho/go-aigc"

var Models = struct {
	// OpenAI GPT models
	OpenAIGpt4oMini                  aigc.ModelId
	OpenAIGpt4oMini_20240718         aigc.ModelId
	OpenAIGpt4o                      aigc.ModelId
	OpenAIGpt4o_20240513             aigc.ModelId
	OpenAIGpt4o_20240806             aigc.ModelId
	OpenAIGpt4Turbo                  aigc.ModelId
	OpenAIGpt4Turbo_20240409         aigc.ModelId
	OpenAIGpt4Turbo_Preview_20240125 aigc.ModelId
	OpenAIGpt4Turbo_Preview_20231106 aigc.ModelId
	OpenAIGpt4                       aigc.ModelId
	OpenAIGpt4_20230613              aigc.ModelId
	OpenAIGpt4_32k_20230613          aigc.ModelId
	OpenAIGpt4_20230314              aigc.ModelId
	OpenAIGpt35_Turbo                aigc.ModelId
	OpenAIGpt35_Turbo_20240125       aigc.ModelId
	OpenAIGpt35_Turbo_20231106       aigc.ModelId
	OpenAIGpt35_Turbo_16k_20230613   aigc.ModelId
	OpenAIGpt35_Turbo_20230613       aigc.ModelId

	// Anthropic Claude models
	AnthropicClaude35Sonnet          aigc.ModelId
	AnthropicClaude3Opus             aigc.ModelId
	AnthropicClaude3Sonnet           aigc.ModelId
	AnthropicClaude3Haiku            aigc.ModelId
	AnthropicClaude35Sonnet_20240620 aigc.ModelId
	AnthropicClaude3Opus_20240229    aigc.ModelId
	AnthropicClaude3Sonnet_20240229  aigc.ModelId
	AnthropicClaude3Haiku_20240307   aigc.ModelId

	// Bedrock Anthropic Claude models
	BedrockAnthropicClaude35Sonnet          aigc.ModelId
	BedrockAnthropicClaude3Opus             aigc.ModelId
	BedrockAnthropicClaude3Sonnet           aigc.ModelId
	BedrockAnthropicClaude3Haiku            aigc.ModelId
	BedrockAnthropicInstant                 aigc.ModelId
	BedrockAnthropicClaude35Sonnet_20240620 aigc.ModelId
	BedrockAnthropicClaude3Opus_20240229    aigc.ModelId
	BedrockAnthropicClaude3Sonnet_20240229  aigc.ModelId
	BedrockAnthropicClaude3Haiku_20240307   aigc.ModelId

	// Bedrock Llama3 models
	BedrockLlama31_405B aigc.ModelId
	BedrockLlama31_70B  aigc.ModelId
	BedrockLlama31_8B   aigc.ModelId

	// Qwen models
	QwenMax            aigc.ModelId
	QwenMax_20240428   aigc.ModelId
	QwenMax_20240403   aigc.ModelId
	QwenMax_20240107   aigc.ModelId
	QwenMaxLongContext aigc.ModelId
	QwenPlus           aigc.ModelId
	QwenPlus_20240806  aigc.ModelId
	QwenPlus_20240723  aigc.ModelId
	QwenPlus_20240624  aigc.ModelId
	QwenPlus_20240206  aigc.ModelId
	QwenTurbo          aigc.ModelId
	QwenTurbo_20240624 aigc.ModelId
	QwenTurbo_20240206 aigc.ModelId
	QwenVLMax          aigc.ModelId
	QwenVLMax_20240809 aigc.ModelId
	Qwen2_72B_Instruct aigc.ModelId
	Qwen2_57B_Instruct aigc.ModelId
	Qwen2_7B_Instruct  aigc.ModelId
}{
	// OpenAI GPT models
	OpenAIGpt4oMini:                  "gpt-4o-mini",
	OpenAIGpt4oMini_20240718:         "gpt-4o-mini-2024-07-18",
	OpenAIGpt4o:                      "gpt-4o",
	OpenAIGpt4o_20240513:             "gpt-4o-2024-05-13",
	OpenAIGpt4o_20240806:             "gpt-4o-2024-08-06",
	OpenAIGpt4Turbo:                  "gpt-4-turbo",
	OpenAIGpt4Turbo_20240409:         "gpt-4-turbo-2024-04-09",
	OpenAIGpt4Turbo_Preview_20240125: "gpt-4-0125-preview",
	OpenAIGpt4Turbo_Preview_20231106: "gpt-4-1106-preview",
	OpenAIGpt4:                       "gpt-4",
	OpenAIGpt4_20230613:              "gpt-4-0613",
	OpenAIGpt4_32k_20230613:          "gpt-4-32k-0613",
	OpenAIGpt4_20230314:              "gpt-4-0314",
	OpenAIGpt35_Turbo:                "gpt-3.5-turbo",
	OpenAIGpt35_Turbo_20240125:       "gpt-3.5-turbo-0125",
	OpenAIGpt35_Turbo_20231106:       "gpt-3.5-turbo-1106",
	OpenAIGpt35_Turbo_16k_20230613:   "gpt-3.5-turbo-16k-0613",
	OpenAIGpt35_Turbo_20230613:       "gpt-3.5-turbo-0613",

	// Anthropic Claude models
	AnthropicClaude35Sonnet:          "claude-3-5-sonnet-20240620",
	AnthropicClaude3Opus:             "claude-3-opus-20240229",
	AnthropicClaude3Sonnet:           "claude-3-sonnet-20240229",
	AnthropicClaude3Haiku:            "claude-3-haiku-20240307",
	AnthropicClaude35Sonnet_20240620: "claude-3-5-sonnet-20240620",
	AnthropicClaude3Opus_20240229:    "claude-3-opus-20240229",
	AnthropicClaude3Sonnet_20240229:  "claude-3-sonnet-20240229",
	AnthropicClaude3Haiku_20240307:   "claude-3-haiku-20240307",

	// Bedrock Anthropic Claude models
	BedrockAnthropicClaude35Sonnet:          "anthropic.claude-3-5-sonnet-20240620-v1:0",
	BedrockAnthropicClaude3Opus:             "anthropic.claude-3-opus-20240229-v1:0",
	BedrockAnthropicClaude3Sonnet:           "anthropic.claude-3-sonnet-20240229-v1:0",
	BedrockAnthropicClaude3Haiku:            "anthropic.claude-3-haiku-20240307-v1:0",
	BedrockAnthropicInstant:                 "anthropic.claude-instant-v1",
	BedrockAnthropicClaude35Sonnet_20240620: "anthropic.claude-3-5-sonnet-20240620-v1:0",
	BedrockAnthropicClaude3Opus_20240229:    "anthropic.claude-3-opus-20240229-v1:0",
	BedrockAnthropicClaude3Sonnet_20240229:  "anthropic.claude-3-sonnet-20240229-v1:0",
	BedrockAnthropicClaude3Haiku_20240307:   "anthropic.claude-3-haiku-20240307-v1:0",

	// Bedrock Llama3 models
	BedrockLlama31_405B: "meta.llama3-1-405b-instruct-v1:0",
	BedrockLlama31_70B:  "meta.llama3-1-70b-instruct-v1:0",
	BedrockLlama31_8B:   "meta.llama3-1-8b-instruct-v1:0",

	// Qwen models on Aliyun DashScope
	QwenMax:            "qwen-max",
	QwenMax_20240428:   "qwen-max-0428",
	QwenMax_20240403:   "qwen-max-0403",
	QwenMax_20240107:   "qwen-max-0107",
	QwenMaxLongContext: "qwen-max-longcontext",
	QwenPlus:           "qwen-plus",
	QwenPlus_20240806:  "qwen-plus-0806",
	QwenPlus_20240723:  "qwen-plus-0723",
	QwenPlus_20240624:  "qwen-plus-0624",
	QwenPlus_20240206:  "qwen-plus-0206",
	QwenTurbo:          "qwen-turbo",
	QwenTurbo_20240624: "qwen-turbo-0624",
	QwenTurbo_20240206: "qwen-turbo-0206",
	QwenVLMax:          "qwen-vl-max",
	QwenVLMax_20240809: "qwen-vl-max-0809",
	Qwen2_72B_Instruct: "qwen2-72b-instruct",
	Qwen2_57B_Instruct: "qwen2-57b-a14b-instruct",
	Qwen2_7B_Instruct:  "qwen2-7b-instruct",
}
