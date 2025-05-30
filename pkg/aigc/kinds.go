package aigc

type RoleType = string

const (
	RoleTypeUser      = "user"
	RoleTypeAssistant = "assistant"
	RoleTypeSystem    = "system"
	RoleTypeTool      = "tool"
)

type modelType = string

const (
	modelTypeDeepSeekChat      = "deepseek-chat"
	modelTypeDeepSeekReasoning = "deepseek-reasoning"
)

type responseFormatType = string

const (
	ResponseFormatTypeText = "text"
	ResponseFormatTypeJSON = "json_object"
)

type toolChoiceType = string

const (
	ToolChoiceTypeNone     = "none"
	ToolChoiceTypeAuto     = "auto"
	ToolChoiceTypeRequired = "required"
)

const (
	pathTypeChatCompletions = "/v1/chat/completions"
)
