package chat

type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleTool      MessageRole = "tool"
)

type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeImage      ContentType = "image"
	ContentTypeToolCall   ContentType = "tool_call"
	ContentTypeToolResult ContentType = "tool_result"
)

type ImageMediaType string

const (
	ImagePng  ImageMediaType = "image/png"
	ImageJpeg ImageMediaType = "image/jpeg"
	ImageWebp ImageMediaType = "image/webp"
	ImageAvif ImageMediaType = "image/avif"
)

type ContentBlock struct {
	Type ContentType `json:"type"`

	// OpenAI compatible
	Refusal string `json:"refusal,omitempty"`

	// For text content
	Text string `json:"text,omitempty"`

	// For image content
	MediaType ImageMediaType `json:"media_type,omitempty"`
	Data      []byte         `json:"data,omitempty"`
	ImageUrl  string         `json:"image_url,omitempty"`

	// For tool use content
	ToolCallId string `json:"tool_call_id,omitempty"`
	// For tool call
	ToolName  string         `json:"tool_name,omitempty"`
	Arguments map[string]any `json:"arguments,omitempty"`
	// For tool result content. Must be string, []any or map[string]any
	Result any `json:"result,omitempty"`
}

type Message struct {
	Role     MessageRole    `json:"role,omitempty"`
	Name     string         `json:"name,omitempty"` // OpenAI compatible
	Contents []ContentBlock `json:"contents,omitempty"`
}

func (m *Message) Copy() Message {
	var z = Message{
		Role: m.Role,
		Name: m.Name,
	}
	if len(m.Contents) > 0 {
		z.Contents = make([]ContentBlock, len(m.Contents))
		copy(z.Contents, m.Contents)
	}
	return z
}
