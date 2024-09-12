package test

import (
	"github.com/Pooh-Mucho/go-aigc/chat"
	"testing"
)

func Test_AzureOpenAI_Gpt4o_Hello(t *testing.T) {
	Tests.Hello(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Chinese_Poetry(t *testing.T) {
	Tests.Chinese_Poetry(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Emoji(t *testing.T) {
	Tests.Emoji(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Multi_Contents(t *testing.T) {
	Tests.Multi_Contents(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Template_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.Template_System_Injections(t, chat.Models.OpenAIGpt4o, WithAzure)
	}
}

func Test_AzureOpenAI_Gpt4o_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.Template_System_Injections(t, chat.Models.OpenAIGpt4o, WithAzure)
	}
}

func Test_AzureOpenAI_Gpt4o_Tool_Random_Number(t *testing.T) {
	Tests.Tool_Random_Number(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Tool_Add_Single(t *testing.T) {
	Tests.Tool_Add_Single(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Tool_Add_Parallel(t *testing.T) {
	Tests.Tool_Add_Parallel(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Tool_Sum(t *testing.T) {
	Tests.Tool_Sum(t, chat.Models.OpenAIGpt4o, WithAzure)
}

func Test_AzureOpenAI_Gpt4o_Tool_File(t *testing.T) {
	Tests.Tool_File(t, chat.Models.OpenAIGpt4o, WithAzure)
}
