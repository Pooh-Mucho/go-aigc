package test

import (
	"github.com/Pooh-Mucho/go-aigc"
	"github.com/Pooh-Mucho/go-aigc/chat"
	"testing"
)

func Test_Bedrock_Sonnet35_Hello(t *testing.T) {
	Tests.Hello(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Chinese_Poetry(t *testing.T) {
	Tests.Chinese_Poetry(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Emoji(t *testing.T) {
	Tests.Emoji(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Multi_Contents(t *testing.T) {
	Tests.Multi_Contents(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Template_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.Template_System_Injections(t, chat.Models.AnthropicClaude35Sonnet_20240620,
			WithAWS, aigc.WithRegion("us-west-2"))
	}
}

func Test_Bedrock_Sonnet35_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.Template_System_Injections(t, chat.Models.OpenAIGpt4o, WithAzure)
	}
}

func Test_Bedrock_Sonnet35_Tool_Random_Number(t *testing.T) {
	Tests.Tool_Random_Number(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Tool_Add_Single(t *testing.T) {
	Tests.Tool_Add_Single(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Tool_Add_Parallel(t *testing.T) {
	Tests.Tool_Add_Parallel(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Tool_Sum(t *testing.T) {
	Tests.Tool_Sum(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}

func Test_Bedrock_Sonnet35_Tool_File(t *testing.T) {
	Tests.Tool_File(t, chat.Models.AnthropicClaude35Sonnet_20240620,
		WithAWS, aigc.WithRegion("us-west-2"))
}
