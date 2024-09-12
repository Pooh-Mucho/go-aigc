package test

import "testing"

// Tests for Llama 3.1 8b

func Test_Ollama_Llama31_8bq4(t *testing.T) {
	OllamaTest(t, "Llama31_8bq4", "llama3.1:8b-instruct-q4_0")
}

func Test_Ollama_Llama31_8bq8(t *testing.T) {
	OllamaTest(t, "Llama31_8bq8", "llama3.1:8b-instruct-q8_0")
}
