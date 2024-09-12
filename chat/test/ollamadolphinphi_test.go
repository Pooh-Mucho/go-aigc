package test

import "testing"

// Tests for dolphin-phi v2.6 2.7b

func Test_Ollama_DolphinPhi_v26_27bq2(t *testing.T) {
	OllamaTest(t, "dolphinphi_v26_27bq2", "dolphin-phi:2.7b-v2.6-q2_k")
}

func Test_Ollama_DolphinPhi_v26_27bq4(t *testing.T) {
	OllamaTest(t, "dolphinphi_v26_27bq4", "dolphin-phi:2.7b-v2.6-q4_0")
}

func Test_Ollama_DolphinPhi_v26_27bq8(t *testing.T) {
	OllamaTest(t, "dolphinphi_v26_27bq8", "dolphin-phi:2.7b-v2.6-q8_0")
}
