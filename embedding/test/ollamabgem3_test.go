package test

import "github.com/Pooh-Mucho/go-aigc/embedding"
import "testing"

func Test_Ollama_BgeM3_Llama_Documents(t *testing.T) {
	Tests.EmbeddingLlamaDocuments(t, embedding.Models.BaaiBgeM3, WithOllama)
}

func Test_Ollama_BgeM3_Llama_Query_Animal(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.BaaiBgeM3, LlamaQueries.Animal, WithOllama)
}

func Test_Ollama_BgeM3_Llama_Query_LoadWeight_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.BaaiBgeM3, LlamaQueries.LoadWeight_ZH, WithOllama)
}

func Test_Ollama_BgeM3_Llama_Query_Size_ZH(t *testing.T) {
	Tests.EmbeddingLlamaQuery(t, embedding.Models.BaaiBgeM3, LlamaQueries.Size_ZH, WithOllama)
}
