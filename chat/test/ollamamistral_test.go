package test

import "testing"

// Tests for mistral 7b

func Test_Ollama_Mistral_7bq2(t *testing.T) {
	OllamaTest(t, "mistral_7b_q2", "mistral:7b-instruct-q2_k")
}

func Test_Ollama_Mistral_7bq4(t *testing.T) {
	OllamaTest(t, "mistral_7b_q4", "mistral:7b-instruct-q4_0")
}

func Test_Ollama_Mistral_7bq8(t *testing.T) {
	OllamaTest(t, "mistral_7b_q8", "mistral:7b-instruct-q8_0")
}
