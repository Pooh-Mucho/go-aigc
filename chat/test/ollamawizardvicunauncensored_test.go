package test

import "testing"

// Tests for wizard-vicuna-uncensored 7b

func Test_Ollama_WizardVicunaUncensored_7bq2(t *testing.T) {
	OllamaTest(t, "wizard-vicuna-uncensored_7bq2", "wizard-vicuna-uncensored:7b-q2_k")
}

func Test_Ollama_WizardVicunaUncensored_7bq4(t *testing.T) {
	OllamaTest(t, "wizard-vicuna-uncensored_7bq4", "wizard-vicuna-uncensored:7b-q4_0")
}

func Test_Ollama_WizardVicunaUncensored_7bq8(t *testing.T) {
	OllamaTest(t, "wizard-vicuna-uncensored_7bq8", "wizard-vicuna-uncensored:7b-q8_0")
}
