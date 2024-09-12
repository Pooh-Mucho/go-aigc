package chat

import (
	"github.com/Pooh-Mucho/go-aigc"
)

type ModelRequest struct {
	Messages          []Message
	Tools             []Tool
	MaxTokens         aigc.Nullable[int32]
	Temperature       aigc.Nullable[float64]
	TopP              aigc.Nullable[float64]
	ToolChoice        *ToolChoice
	ParallelToolCalls aigc.Nullable[bool]
}

func (r *ModelRequest) Copy() *ModelRequest {
	var z = &ModelRequest{
		MaxTokens:         r.MaxTokens,
		Temperature:       r.Temperature,
		TopP:              r.TopP,
		ParallelToolCalls: r.ParallelToolCalls,
	}

	if len(r.Messages) > 0 {
		z.Messages = make([]Message, len(r.Messages))
		for i, _ := range r.Messages {
			z.Messages[i] = r.Messages[i].Copy()
		}
	}
	if len(r.Tools) > 0 {
		z.Tools = make([]Tool, len(r.Tools))
		copy(z.Tools, r.Tools)
	}
	if r.ToolChoice != nil {
		z.ToolChoice = new(ToolChoice)
		*z.ToolChoice = *r.ToolChoice
	}

	return z
}
