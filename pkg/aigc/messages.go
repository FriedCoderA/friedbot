package aigc

type Message interface{}

type Messages []Message

type BaseMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type SystemMessage struct {
	*BaseMessage
	Name string `json:"name,omitempty"`
}

type UserMessage struct {
	*BaseMessage
	Name string `json:"name,omitempty"`
}

type AssistantMessage struct {
	*BaseMessage
	Name             string `json:"name,omitempty"`
	Prefix           bool   `json:"prefix,omitempty"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

type ToolMessage struct {
	*BaseMessage
	ToolCallID string `json:"tool_call_id,omitempty"`
}

type ResponseMessage struct {
	*BaseMessage
	ReasoningContent string `json:"reasoning_content,omitempty"`
	ToolCalls        []Tool `json:"tool_calls,omitempty"`
}

func NewSystemMessage(content, name string) *SystemMessage {
	return &SystemMessage{
		BaseMessage: &BaseMessage{
			Role:    RoleTypeSystem,
			Content: content,
		},
		Name: name,
	}
}

func NewUserMessage(content, name string) *UserMessage {
	return &UserMessage{
		BaseMessage: &BaseMessage{
			Role:    RoleTypeUser,
			Content: content,
		},
		Name: name,
	}
}

func NewAssistantMessage(content, name string, prefix bool, reasoningContent string) *AssistantMessage {
	return &AssistantMessage{
		BaseMessage: &BaseMessage{
			Role:    RoleTypeAssistant,
			Content: content,
		},
		Name:             name,
		Prefix:           prefix,
		ReasoningContent: reasoningContent,
	}
}

func NewToolMessage(content, toolCallID string) *ToolMessage {
	return &ToolMessage{
		BaseMessage: &BaseMessage{
			Role:    RoleTypeTool,
			Content: content,
		},
		ToolCallID: toolCallID,
	}
}
