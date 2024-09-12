package test

import (
	"testing"
)

// phi3 3.8b q2_K is so bad, we don't test it.
// func Test_Ollama_Phi3_38bq2(t *testing.T) {
// 	OllamaTest(t, "Phi3_3.8bq2", "phi3:3.8b-mini-128k-instruct-q2_K")
// }

func Test_Ollama_Phi3_38b_128k_q4(t *testing.T) {
	OllamaTest(t, "Phi3_3.8b_128k_q4", "phi3:3.8b-mini-128k-instruct-q4_0")
}

func Test_Ollama_Phi3_38b_128k_q8(t *testing.T) {
	OllamaTest(t, "Phi3_3.8b_128k_q8", "phi3:3.8b-mini-128k-instruct-q8_0")
}

func Test_Ollama_Phi3_38b_4k_q4(t *testing.T) {
	OllamaTest(t, "Phi3_3.8b_4k_q4", "phi3:3.8b-mini-4k-instruct-q4_0")
}

func Test_Ollama_Phi3_38b_4k_q8(t *testing.T) {
	OllamaTest(t, "Phi3_3.8b_4k_q8", "phi3:3.8b-mini-4k-instruct-q8_0")
}

// phi3.5 3.8b q2_K is so bad, we don't test it.
// func Test_Ollama_Phi35_38b_128k_q2(t *testing.T) {
// 	OllamaTest(t, "Phi3.5_3.8b_128k_q2", "phi3.5:3.8b-mini-instruct-q2_k")
// }

func Test_Ollama_Phi35_38b_128k_q4(t *testing.T) {
	OllamaTest(t, "Phi3.5_3.8b_128k_q4", "phi3.5:3.8b-mini-instruct-q4_0")
}

func Test_Ollama_Phi35_38b_128k_q8(t *testing.T) {
	OllamaTest(t, "Phi3.5_3.8b_128k_q8", "phi3.5:3.8b-mini-instruct-q8_0")
}
