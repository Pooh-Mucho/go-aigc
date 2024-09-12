package test

import (
	"github.com/Pooh-Mucho/go-aigc/embedding"
	"testing"
)

func Test_OpenAI_TextEmbedding3Small_Llama_Documents(t *testing.T) {
	Tests.EmbeddingLlamaDocuments(t, embedding.Models.OpenAITextEmbedding3Small, WithOpenAI)
}

func Test_OpenAI_TextEmbedding3Small_Llama_Query_Animal(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.OpenAITextEmbedding3Small, LlamaQueries.Animal, WithOpenAI)
}

func Test_OpenAI_TextEmbedding3Small_Llama_Query_LoadWeight_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.OpenAITextEmbedding3Small, LlamaQueries.LoadWeight_ZH, WithOpenAI)
}

func Test_OpenAI_TextEmbedding3Small_Llama_Query_Size_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.OpenAITextEmbedding3Small, LlamaQueries.Size_ZH, WithOpenAI)
}
