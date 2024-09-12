package chat

import (
	"context"
	"errors"
	"fmt"
)

const (
	toolExecutorDefaultRoundtrips = 10
	toolExecutorMaxRoundtrips     = 20
)

type ToolCallResult struct {
	Content ContentBlock
	Result  any
}

type Roundtrip struct {
	Request         *ModelRequest
	Response        *ModelResponse
	ToolCallResults []ToolCallResult
}

type ToolExecutor struct {
	Model          Model
	InitialRequest *ModelRequest
	Dispatcher     func(callId string, toolName string, args map[string]any) (result any, err error)

	roundtrips    []Roundtrip
	maxRoundtrips int
}

func (e *ToolExecutor) GetMaxRoundtrips() int {
	if e.maxRoundtrips > 0 {
		return e.maxRoundtrips
	}
	return toolExecutorDefaultRoundtrips
}

func (e *ToolExecutor) SetMaxRoundtrips(value int) {
	if value <= 0 {
		panic(errors.New("[ToolExecutor.SetMaxRoundtrips] value must be greater than 0"))
	}
	if value > toolExecutorMaxRoundtrips {
		panic(fmt.Errorf("[ToolExecutor.SetMaxRoundtrips] value must be less than or equal to %d",
			toolExecutorMaxRoundtrips))
	}
	e.maxRoundtrips = value
}

func (e *ToolExecutor) Roundtrips() int {
	return len(e.roundtrips)
}

func (e *ToolExecutor) GetRoundtrip(index int) Roundtrip {
	return e.roundtrips[index]
}

func (e *ToolExecutor) LastRoundtrip() Roundtrip {
	if len(e.roundtrips) == 0 {
		panic(errors.New("[ToolExecutor.LastRoundtrip] not executed"))
	}
	return e.roundtrips[len(e.roundtrips)-1]
}

func (e *ToolExecutor) LastRequest() *ModelRequest {
	if len(e.roundtrips) == 0 {
		return nil
	}
	return e.roundtrips[len(e.roundtrips)-1].Request
}

func (e *ToolExecutor) LastResponse() *ModelResponse {
	if len(e.roundtrips) == 0 {
		return nil
	}
	return e.roundtrips[len(e.roundtrips)-1].Response
}

func (e *ToolExecutor) Execute(ctx context.Context) (finished bool, err error) {
	var (
		request  *ModelRequest
		response *ModelResponse
		last     *Roundtrip
	)

	// Check settings
	if e.Model == nil {
		panic(errors.New("[ToolExecutor.Execute] Model is not set"))
	}
	if e.InitialRequest == nil {
		panic(errors.New("[ToolExecutor.Execute] InitialRequest is not set"))
	}
	if e.Dispatcher == nil {
		var lastRequest = e.LastRequest()
		if lastRequest == nil {
			lastRequest = e.InitialRequest
		}
		for _, tool := range lastRequest.Tools {
			if tool.Function == nil {
				return false, fmt.Errorf(
					"[ToolExecutor.Execute] Dispatcher or tool.Function is not set, can not call '%s'",
					tool.Name)
			}
		}
	}

	if len(e.roundtrips) == 0 {
		e.roundtrips = append(e.roundtrips, Roundtrip{Request: e.InitialRequest})
	}
	last = &e.roundtrips[len(e.roundtrips)-1]

	if len(e.roundtrips) >= e.GetMaxRoundtrips() && last.Response != nil {
		return false, errors.New("[ToolExecutor.Execute] max roundtrips reached")
	}

	if last.Response != nil {
		if last.Response.FinishReason.Type() != FinishReasonToolCalls {
			return true, errors.New("[ToolExecutor.Execute] tool calls are already completed")
		}
		if len(last.ToolCallResults) == 0 {
			last.ToolCallResults, err = e.callTools(last.Response)
			if err != nil {
				return false, fmt.Errorf("[ToolExecutor.Execute] %w", err)
			}
			if len(last.ToolCallResults) == 0 {
				return false, errors.New("[ToolExecutor.Execute] no tool calls")
			}
		}
		request, err = e.getRequest()
		if err != nil {
			return false, fmt.Errorf("[ToolExecutor.Execute] %w", err)
		}
		e.roundtrips = append(e.roundtrips, Roundtrip{Request: request})
		last = &e.roundtrips[len(e.roundtrips)-1]
	} else {
		request = last.Request
	}

	response, err = e.Model.Complete(ctx, request)
	if err != nil {
		return false, fmt.Errorf("[ToolExecutor.Execute] %w", err)
	}

	last.Response = response

	if response.FinishReason.Type() != FinishReasonToolCalls {
		return true, nil
	}

	return false, nil
}

func (e *ToolExecutor) getRequest() (*ModelRequest, error) {
	if len(e.roundtrips) == 0 {
		return e.InitialRequest, nil
	}

	var request *ModelRequest
	var last = e.roundtrips[len(e.roundtrips)-1]
	var callResult Message

	if last.Response.FinishReason.Type() != FinishReasonToolCalls {
		return nil, errors.New("[ToolExecutor.getRequest] last response is not a tool call")
	}
	if len(last.ToolCallResults) == 0 {
		return nil, errors.New("[ToolExecutor.getRequest] no tool call results")
	}

	request = last.Request.Copy()
	for i := range last.Response.Messages {
		request.Messages = append(request.Messages, last.Response.Messages[i].Copy())
	}
	callResult.Role = RoleTool
	for _, result := range last.ToolCallResults {
		var content = ContentBlock{Type: ContentTypeToolResult}
		content.ToolCallId = result.Content.ToolCallId
		content.ToolName = result.Content.ToolName
		content.Result = result.Result
		callResult.Contents = append(callResult.Contents, content)
	}
	request.Messages = append(request.Messages, callResult)

	return request, nil
}

func (e *ToolExecutor) callTools(response *ModelResponse) ([]ToolCallResult, error) {
	var err error
	var results []ToolCallResult

	for _, message := range response.Messages {
		for _, content := range message.Contents {
			if content.Type != ContentTypeToolCall {
				continue
			}

			var result = ToolCallResult{Content: content}
			var ok = false

			var request = e.LastRequest()
			for _, tool := range request.Tools {
				if tool.Name != content.ToolName {
					continue
				}
				if tool.Function == nil {
					break
				}
				result.Result, err = tool.Function(content.Arguments)
				if err != nil {
					return nil, fmt.Errorf("[ToolExecutor.callTools] %w", err)
				}
				ok = true
				break
			}

			if !ok {
				if e.Dispatcher == nil {
					return nil, fmt.Errorf("[ToolExecutor.callTools] unknown tool: %s", content.ToolName)
				}
				result.Result, err = e.Dispatcher(content.ToolCallId, content.ToolName, content.Arguments)
				if err != nil {
					return nil, fmt.Errorf("[ToolExecutor.callTools] %w", err)
				}
			}

			results = append(results, result)
		}
	}

	return results, nil
}
