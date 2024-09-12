package chat

import "github.com/Pooh-Mucho/go-aigc"

type ToolParameters struct {
	Properties []aigc.JsonSchemaProperty `json:"properties"`
	Required   []string                  `json:"required,omitempty"`
}

type Tool struct {
	// The name of the tool
	Name string `json:"name,omitempty"`
	// Type        string
	// The description of the tool
	Description string `json:"description,omitempty"`
	// The input parameters schema
	Parameters ToolParameters `json:"parameters"`
	// OpenAI compatible.
	// https://platform.openai.com/docs/guides/function-calling
	Strict bool `json:"strict,omitempty"`
	// The function to call
	Function func(map[string]interface{}) (any, error) `json:"-"`
}

// Anthropic:
// 	 "tool_choice":{"type":"auto"}
//	 "tool_choice":{"type":"any"}
//	 "tool_choice":{"type":"tool", "name":"get_weather"}}
// OpenAI:
// 	 "tool_choice":"none"
// 	 "tool_choice":"auto"
// 	 "tool_choice":"required"
// 	 "tool_choice":{"type": "function", "function": {"name": "get_weather"}}
// 	 "tool_choice":{"type": "file_search"}
//

type ToolChoiceType string

const (
	ToolChoiceTypeAuto     ToolChoiceType = "auto"
	ToolChoiceTypeRequired ToolChoiceType = "required"
	ToolChoiceTypeRestrict ToolChoiceType = "restrict"
)

type ToolChoice struct {
	Type ToolChoiceType
	Name string
}

var (
	ToolChoiceAuto     = ToolChoice{Type: ToolChoiceTypeAuto}
	ToolChoiceRequired = ToolChoice{Type: ToolChoiceTypeRequired}
)
