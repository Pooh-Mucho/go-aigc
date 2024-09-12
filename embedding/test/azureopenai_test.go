package test

import (
	"github.com/Pooh-Mucho/go-aigc/embedding"
	"testing"
)

func Test_AzureOpenAI_Ada002_Embedding_Llama_Documents(t *testing.T) {
	Tests.EmbeddingLlamaDocuments(t, embedding.Models.OpenAITextEmbeddingAda_002, WithAzure)
}

func Test_AzureOpenAI_Ada002_Embedding_Llama_Query_Animal(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.OpenAITextEmbeddingAda_002, LlamaQueries.Animal, WithAzure)
}

func Test_AzureOpenAI_Ada002_Embedding_Llama_Query_LoadWeight_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.OpenAITextEmbeddingAda_002, LlamaQueries.LoadWeight_ZH, WithAzure)
}

func Test_AzureOpenAI_Ada002_Embedding_Llama_Query_Size_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.OpenAITextEmbeddingAda_002, LlamaQueries.Size_ZH, WithAzure)
}
