package test

import "testing"

// Tests for dolphin mistral v2.6 7b

func Test_Ollama_DolphinMistral_v26_7bq2(t *testing.T) {
	OllamaTest(t, "dolphinmistral_v26_7bq2", "dolphin-mistral:7b-v2.6-dpo-laser-q2_k")
}

func Test_Ollama_DolphinMistral_v26_7bq4(t *testing.T) {
	OllamaTest(t, "dolphinmistral_v26_7bq4", "dolphin-mistral:7b-v2.6-dpo-laser-q4_0")
}

func Test_Ollama_DolphinMistral_v26_7bq8(t *testing.T) {
	OllamaTest(t, "dolphinmistral_v26_7bq8", "dolphin-mistral:7b-v2.6-dpo-laser-q8_0")
}

// Tests for dolphin mistral v2.8 7b

func Test_Ollama_DolphinMistral_v28_7bq2(t *testing.T) {
	OllamaTest(t, "dolphinmistral_v28_7bq2", "dolphin-mistral:7b-v2.8-q2_k")
}

func Test_Ollama_DolphinMistral_v28_7bq4(t *testing.T) {
	OllamaTest(t, "dolphinmistral_v28_7bq4", "dolphin-mistral:7b-v2.8-q4_0")
}

func Test_Ollama_DolphinMistral_v28_7bq8(t *testing.T) {
	OllamaTest(t, "dolphinmistral_v28_7bq8", "dolphin-mistral:7b-v2.8-q8_0")
}
