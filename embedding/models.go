package embedding

import "github.com/Pooh-Mucho/go-aigc"

var Models = struct {
	// OpenAI Embedding models
	OpenAITextEmbeddingAda_002 aigc.ModelId
	OpenAITextEmbedding3Small  aigc.ModelId
	OpenAITextEmbedding3Large  aigc.ModelId

	// BAAI Embedding models
	// https://huggingface.co/BAAI
	// https://github.com/FlagOpen/FlagEmbedding/blob/master/README_zh.md
	BaaiBgeM3           aigc.ModelId
	BaaiBgeRerankerV2M3 aigc.ModelId

	// NOMIC Embedding models
	// https://huggingface.co/nomic-ai
	// https://huggingface.co/nomic-ai/nomic-embed-text-v1.5
	NomicEmbedText    aigc.ModelId
	NomicEmbedTextV1  aigc.ModelId
	NomicEmbedTextV15 aigc.ModelId

	// Mixedbread Embedding models
	// https://huggingface.co/mixedbread-ai
	// https://huggingface.co/mixedbread-ai/mxbai-embed-large-v1
	MxbaiEmbedLarge   aigc.ModelId
	MxbaiEmbedLargeV1 aigc.ModelId
}{
	// OpenAI Embedding models
	OpenAITextEmbeddingAda_002: "text-embedding-ada-002",
	OpenAITextEmbedding3Small:  "text-embedding-3-small",
	OpenAITextEmbedding3Large:  "text-embedding-3-large",

	// BAAI Embedding models
	BaaiBgeM3:           "bge-m3",
	BaaiBgeRerankerV2M3: "bge-reranker-v2-m3",

	// NOMIC Embedding models
	NomicEmbedText:    "nomic-embed-text",
	NomicEmbedTextV1:  "nomic-embed-text-v1",
	NomicEmbedTextV15: "nomic-embed-text-v1.5",

	// Mixedbread Embedding models
	MxbaiEmbedLarge:   "mxbai-embed-large",
	MxbaiEmbedLargeV1: "mxbai-embed-large-v1",
}
