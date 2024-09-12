package test

import "github.com/Pooh-Mucho/go-aigc/embedding"
import "testing"

func Test_Ollama_NomicEmbedText_Llama_Documents(t *testing.T) {
	Tests.EmbeddingLlamaDocuments(t, embedding.Models.NomicEmbedText, WithOllama)
}

func Test_Ollama_NomicEmbedText_Llama_Query_Animal(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.NomicEmbedText, LlamaQueries.Animal, WithOllama)
}

func Test_Ollama_NomicEmbedText_Llama_Query_LoadWeight_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.NomicEmbedText, LlamaQueries.LoadWeight_ZH, WithOllama)
}

func Test_Ollama_NomicEmbedText_Llama_Query_Size_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.NomicEmbedText, LlamaQueries.Size_ZH, WithOllama)
}
