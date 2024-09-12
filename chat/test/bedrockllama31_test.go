package test

import (
	"github.com/Pooh-Mucho/go-aigc"
	"github.com/Pooh-Mucho/go-aigc/chat"
	"testing"
)

func Test_Bedrock_Llama31_70B_Hello(t *testing.T) {
	Tests.Hello(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Chinese_Poetry(t *testing.T) {
	Tests.Chinese_Poetry(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Emoji(t *testing.T) {
	Tests.Emoji(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Multi_Contents(t *testing.T) {
	Tests.Multi_Contents(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Template_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.Template_System_Injections(t, chat.Models.BedrockLlama31_70B,
			WithAWS, aigc.WithRegion("us-west-2"))
	}
}

func Test_Bedrock_Llama31_70B_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.System_Injections(t, chat.Models.BedrockLlama31_70B,
			WithAWS, aigc.WithRegion("us-west-2"))
	}
}

func Test_Bedrock_Llama31_70B_Tool_Random_Number(t *testing.T) {
	Tests.Tool_Random_Number(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Tool_Add_Single(t *testing.T) {
	Tests.Tool_Add_Single(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Tool_Add_Parallel(t *testing.T) {
	Tests.Tool_Add_Parallel(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Tool_Sum(t *testing.T) {
	Tests.Tool_Sum(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Llama31_70B_Tool_File(t *testing.T) {
	Tests.Tool_File(t, chat.Models.BedrockLlama31_70B,
		WithAWS, aigc.WithRegion("us-west-2"))
}
