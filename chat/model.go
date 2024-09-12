package chat

import (
	"context"
	"fmt"
	"github.com/Pooh-Mucho/go-aigc"
	"strings"
)

type Model interface {
	GetModelId() string
	Complete(ctx context.Context, request *ModelRequest) (*ModelResponse, error)
}

func NewModel(modelId aigc.ModelId, options ...aigc.ModelOptionFunc) (Model, error) {
	var opts aigc.ModelOptions

	for _, fn := range options {
		fn(&opts)
	}

	if strings.Index(string(modelId), "gpt") >= 0 {
		if opts.VendorId == aigc.Vendors.Microsoft {
			return newAzureGptModel(string(modelId), &opts)
		}
		if opts.VendorId == aigc.Vendors.OpenAI {
			return newOpenAIGptModel(string(modelId), &opts)
		}
		if opts.VendorId == "" {
			opts.VendorId = aigc.Vendors.OpenAI
			return newOpenAIGptModel(string(modelId), &opts)
		}
	}

	// anthropic.claude-3-5-sonnet-20240620-v1:0
	// claude-3-5-sonnet-20240620
	if strings.Index(string(modelId), "claude") >= 0 {
		if opts.VendorId == aigc.Vendors.Anthropic {
			return newAnthropicClaudeModel(string(modelId), &opts)
		}
		if opts.VendorId == "" {
			opts.VendorId = aigc.Vendors.Anthropic
			return newAnthropicClaudeModel(string(modelId), &opts)
		}
		if opts.VendorId == aigc.Vendors.Amazon {
			return newBedrockClaudeModel(string(modelId), &opts)
		}
	}

	if strings.Index(string(modelId), "llama3") >= 0 {
		if opts.VendorId == aigc.Vendors.Amazon {
			return newBedrockLlama3Model(string(modelId), &opts)
		}
	}

	if strings.Index(string(modelId), "qwen") >= 0 {
		if opts.VendorId == aigc.Vendors.Alibaba {
			return newDashScopeQwenModel(string(modelId), &opts)
		}
	}

	if strings.Index(string(modelId), "voyage") >= 0 {

	}

	if opts.VendorId == aigc.Vendors.Ollama {
		return newOllamaChatModel(string(modelId), &opts)
	}

	if opts.VendorId == aigc.Vendors.PoohMucho {
		return newPoohMuchoChatModel(string(modelId), &opts)
	}

	if opts.VendorId != "" {
		return nil, fmt.Errorf("model can not be created, vendor:%s, model: %s", opts.VendorId, modelId)
	}
	return nil, fmt.Errorf("model can not be created: %s", modelId)
}
