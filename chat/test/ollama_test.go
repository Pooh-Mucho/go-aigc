package test

import (
	"github.com/Pooh-Mucho/go-aigc"
	"testing"
)

func OllamaTest(t *testing.T, modelName string, modelId aigc.ModelId) {
	t.Run("Test_Ollama_"+modelName+"_Hello", func(t *testing.T) {
		Tests.Hello(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Chinese_Poetry", func(t *testing.T) {
		Tests.Chinese_Poetry(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Emoji", func(t *testing.T) {
		Tests.Emoji(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Multi_Contents", func(t *testing.T) {
		Tests.Multi_Contents(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Template_System_Injections", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			Tests.Template_System_Injections(t, modelId, WithOllama)
		}
	})
	t.Run("Test_Ollama_"+modelName+"_System_Injections", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			Tests.Template_System_Injections(t, modelId, WithOllama)
		}
	})
	t.Run("Test_Ollama_"+modelName+"_Tool_Random_Number", func(t *testing.T) {
		Tests.Tool_Random_Number(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Tool_Add_Single", func(t *testing.T) {
		Tests.Tool_Add_Single(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Tool_Add_Parallel", func(t *testing.T) {
		Tests.Tool_Add_Parallel(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Tool_Sum", func(t *testing.T) {
		Tests.Tool_Sum(t, modelId, WithOllama)
	})
	t.Run("Test_Ollama_"+modelName+"_Tool_File", func(t *testing.T) {
		Tests.Tool_File(t, modelId, WithOllama)
	})
}
