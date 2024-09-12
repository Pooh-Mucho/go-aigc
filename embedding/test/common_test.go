package test

import (
	"context"
	"github.com/Pooh-Mucho/go-aigc"
	"github.com/Pooh-Mucho/go-aigc/embedding"
	"slices"
	"testing"
)

var LlamaDocuments = []string{
	"Llamas are members of the camelid family meaning they're pretty closely related to vicuñas and camels",
	"Llamas were first domesticated and used as pack animals 4,000 to 5,000 years ago in the Peruvian highlands",
	"Llamas can grow as much as 6 feet tall though the average llama between 5 feet 6 inches and 5 feet 9 inches tall",
	"Llamas weigh between 280 and 450 pounds and can carry 25 to 30 percent of their body weight",
	"Llamas are vegetarians and have very efficient digestive systems",
	"Llamas live to be about 20 years old, though some only live for 15 years and others live to be 30 years old",
}

type DocumentEmbedding struct {
	Document  string
	Embedding []float32
}

type DocumentSimilarity struct {
	DocumentEmbedding
	Similarity float32
}

var LlamaQueries = struct {
	Animal        string
	LoadWeight_ZH string
	Size_ZH       string
}{
	Animal:        "What animals are llamas related to?",
	LoadWeight_ZH: "羊驼能拉多重的东西?",
	Size_ZH:       "羊驼能长多大?",
}

var Tests = struct {
	EmbeddingLlamaDocuments func(t *testing.T, modelId aigc.ModelId, options ...aigc.ModelOptionFunc)
	EmbeddingLlamaQuery     func(t *testing.T, modelId aigc.ModelId, query string, options ...aigc.ModelOptionFunc)
}{
	EmbeddingLlamaDocuments: func(t *testing.T, modelId aigc.ModelId, options ...aigc.ModelOptionFunc) {
		model, err := embedding.NewModel(modelId, options...)
		if err != nil {
			t.Fatal(err)
		}
		var vectors = make([]DocumentEmbedding, len(LlamaDocuments))
		for i, doc := range LlamaDocuments[:2] {
			response, err := model.Embedding(context.Background(), &embedding.ModelRequest{Document: doc})
			if err != nil {
				t.Fatal(err)
			}
			vectors[i].Document = doc
			vectors[i].Embedding = response.Embedding
			mag, err := embedding.VectorMagnitude(vectors[i].Embedding)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("magnitude:", mag, "dimension:", len(vectors[i].Embedding))
		}
	},

	EmbeddingLlamaQuery: func(t *testing.T, modelId aigc.ModelId, query string, options ...aigc.ModelOptionFunc) {
		model, err := embedding.NewModel(modelId, options...)
		if err != nil {
			t.Fatal(err)
		}

		similarities := make([]DocumentSimilarity, len(LlamaDocuments))
		for i, doc := range LlamaDocuments {
			response, err := model.Embedding(context.Background(), &embedding.ModelRequest{Document: doc})
			if err != nil {
				t.Fatal(err)
			}
			similarities[i].Document = doc
			similarities[i].Embedding = response.Embedding
		}
		// query := "What animals are llamas related to?"
		response, err := model.Embedding(context.Background(), &embedding.ModelRequest{Document: query})
		if err != nil {
			t.Fatal(err)
		}
		queryEmbedding := response.Embedding
		for i, _ := range similarities {
			distance, err := model.Distance(queryEmbedding, similarities[i].Embedding)
			if err != nil {
				t.Fatal(err)
			}
			similarities[i].Similarity = distance
		}

		slices.SortFunc(similarities, func(i, j DocumentSimilarity) int {
			if i.Similarity < j.Similarity {
				return 1
			}
			if i.Similarity > j.Similarity {
				return -1
			}
			return 0
		})

		for _, s := range similarities {
			t.Logf("similarity: %f, document: %s", s.Similarity, s.Document)
		}
	},
}
