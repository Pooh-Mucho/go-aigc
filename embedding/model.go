package embedding

import (
	"context"
	"fmt"
	"github.com/Pooh-Mucho/go-aigc"
)

type VectorDistanceType string

const (
	CosineDistance    VectorDistanceType = "cosine"
	EuclideanDistance VectorDistanceType = "euclidean"
)

type Model interface {
	GetModelId() string
	GetDistanceType() VectorDistanceType
	Distance(vector1, vector2 []float32) (float32, error)
	Embedding(ctx context.Context, request *ModelRequest) (*ModelResponse, error)
}

func NewModel(modelId aigc.ModelId, options ...aigc.ModelOptionFunc) (Model, error) {
	var opts aigc.ModelOptions

	for _, fn := range options {
		fn(&opts)
	}

	// text-embedding-ada-002, text-embedding-3-small, text-embedding-3-large
	switch modelId {
	case Models.OpenAITextEmbeddingAda_002,
		Models.OpenAITextEmbedding3Small,
		Models.OpenAITextEmbedding3Large:
		if opts.VendorId == aigc.Vendors.Microsoft {
			return newAzureOpenAIEmbeddingModel(string(modelId), &opts)
		}
		if opts.VendorId == aigc.Vendors.OpenAI {
			return newOpenAIEmbeddingModel(string(modelId), &opts)
		}
		if opts.VendorId == "" {
			opts.VendorId = aigc.Vendors.OpenAI
			return newOpenAIEmbeddingModel(string(modelId), &opts)
		}
	}

	if opts.VendorId == aigc.Vendors.Ollama {
		return newOllamaEmbeddingModel(string(modelId), &opts)
	}

	if opts.VendorId != "" {
		return nil, fmt.Errorf("model can not be created, vendor:%s, model: %s", opts.VendorId, modelId)
	}
	return nil, fmt.Errorf("model can not be created: %s", modelId)

}
