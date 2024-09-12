package chat

type FinishReason string

type FinishReasonType uint32

const (
	FinishReasonUnknown       FinishReasonType = 0
	FinishReasonStop          FinishReasonType = 1
	FinishReasonLength        FinishReasonType = 2
	FinishReasonToolCalls     FinishReasonType = 3
	FinishReasonContentFilter FinishReasonType = 4
)

func (r FinishReason) Type() FinishReasonType {
	switch r {
	// OpenAI finish_reason. "stop" and "length" also for Llama3
	case "stop":
		return FinishReasonStop
	case "length":
		return FinishReasonLength
	case "tool_calls":
		return FinishReasonToolCalls
	case "content_filter":
		return FinishReasonContentFilter
	}
	switch r {
	// Anthropic stop_reason
	case "end_turn":
		return FinishReasonStop
	case "max_tokens":
		return FinishReasonLength
	case "stop_sequence":
		return FinishReasonStop
	case "tool_use":
		return FinishReasonToolCalls
	}
	return FinishReasonUnknown
}

type TokenUsage struct {
	InputTokens  int
	OutputTokens int
}

type ModelResponse struct {
	Id                  string
	Messages            []Message
	FinishReason        FinishReason
	Usage               TokenUsage
	ContentFilterResult string
}
