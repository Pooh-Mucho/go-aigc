package test

import (
	"testing"
)

// Tests for gemma2 2b

func Test_Ollama_Gemma2_2bq2(t *testing.T) {
	OllamaTest(t, "gemma2_2bq2", "gemma2:2b-instruct-q2_k")
}

func Test_Ollama_Gemma2_2bq4(t *testing.T) {
	OllamaTest(t, "gemma2_2bq4", "gemma2:2b-instruct-q4_0")
}

func Test_Ollama_Gemma2_2bq8(t *testing.T) {
	OllamaTest(t, "gemma2_2bq8", "gemma2:2b-instruct-q8_0")
}

// Tests for gemma2 9b

func Test_Ollama_Gemma2_9bq2(t *testing.T) {
	OllamaTest(t, "gemma2_9bq2", "gemma2:9b-instruct-q2_k")
}

func Test_Ollama_Gemma2_9bq4(t *testing.T) {
	OllamaTest(t, "gemma2_9bq4", "gemma2:9b-instruct-q4_0")
}

func Test_Ollama_Gemma2_9bq8(t *testing.T) {
	OllamaTest(t, "gemma2_9bq8", "gemma2:9b-instruct-q8_0")
}
