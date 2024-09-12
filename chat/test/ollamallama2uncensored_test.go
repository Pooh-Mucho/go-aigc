package test

import "testing"

// Tests for llama2-uncensored 7b

func Test_Ollama_Llama2Uncensored_7bq2(t *testing.T) {
	OllamaTest(t, "Llama2Uncensored_7bq2", "llama2-uncensored:7b-chat-q2_k")
}

func Test_Ollama_Llama2Uncensored_7bq4(t *testing.T) {
	OllamaTest(t, "Llama2Uncensored_7bq4", "llama2-uncensored:7b-chat-q4_0")
}

func Test_Ollama_Llama2Uncensored_7bq8(t *testing.T) {
	OllamaTest(t, "Llama2Uncensored_7bq8", "llama2-uncensored:7b-chat-q8_0")
}
