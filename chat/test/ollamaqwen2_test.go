package test

import (
	"testing"
)

// Tests for qwen2 0.5b

func Test_Ollama_Qwen2_05bq2(t *testing.T) {
	OllamaTest(t, "Qwen2_0.5b_q2", "qwen2:0.5b-instruct-q2_k")
}

func Test_Ollama_Qwen2_05bq4(t *testing.T) {
	OllamaTest(t, "Qwen2_0.5b_q4", "qwen2:0.5b-instruct-q4_0")
}

func Test_Ollama_Qwen2_05bq8(t *testing.T) {
	OllamaTest(t, "Qwen2_0.5b_q8", "qwen2:0.5b-instruct-q8_0")
}

// Tests for qwen2 1.5b

func Test_Ollama_Qwen2_15bq2(t *testing.T) {
	OllamaTest(t, "Qwen2_1.5b_q2", "qwen2:1.5b-instruct-q2_k")
}

func Test_Ollama_Qwen2_15bq4(t *testing.T) {
	OllamaTest(t, "Qwen2_1.5b_q4", "qwen2:1.5b-instruct-q4_0")
}

func Test_Ollama_Qwen2_15bq8(t *testing.T) {
	OllamaTest(t, "Qwen2_1.5b_q8", "qwen2:1.5b-instruct-q8_0")
}

// Tests for qwen2 7b

func Test_Ollama_Qwen2_7bq2(t *testing.T) {
	OllamaTest(t, "Qwen2_7b_q2", "qwen2:7b-instruct-q2_k")
}

func Test_Ollama_Qwen2_7bq4(t *testing.T) {
	OllamaTest(t, "Qwen2_7b_q4", "qwen2:7b-instruct-q4_0")
}

func Test_Ollama_Qwen2_7bq8(t *testing.T) {
	OllamaTest(t, "Qwen2_7b_q8", "qwen2:7b-instruct-q8_0")
}
