package test

import "github.com/Pooh-Mucho/go-aigc/embedding"
import "testing"

func Test_Ollama_MxbaiEmbedLarge_Llama_Documents(t *testing.T) {
	Tests.EmbeddingLlamaDocuments(t, embedding.Models.MxbaiEmbedLarge, WithOllama)
}

func Test_Ollama_MxbaiEmbedLarge_Llama_Query_Animal(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.MxbaiEmbedLarge, LlamaQueries.Animal, WithOllama)
}

func Test_Ollama_MxbaiEmbedLarge_Llama_Query_LoadWeight_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.MxbaiEmbedLarge, LlamaQueries.LoadWeight_ZH, WithOllama)
}

func Test_Ollama_MxbaiEmbedLarge_Llama_Query_Size_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.MxbaiEmbedLarge, LlamaQueries.Size_ZH, WithOllama)
}
