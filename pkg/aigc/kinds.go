package aigc

type roleType = string

const (
	roleTypeUser      = "user"
	roleTypeAssistant = "assistant"
	roleTypeSystem    = "system"
	roleTypeTool      = "tool"
)

type modelType = string

const (
	modelTypeDeepSeekChat      = "deepseek-chat"
	modelTypeDeepSeekReasoning = "deepseek-reasoning"
)

type responseFormatType = string

const (
	responseFormatTypeText = "text"
	responseFormatTypeJSON = "json_object"
)

type toolChoiceType = string

const (
	toolChoiceTypeNone     = "none"
	toolChoiceTypeAuto     = "auto"
	toolChoiceTypeRequired = "required"
)

type pathType = string

const (
	pathTypeChatCompletions = "/v1/chat/completions"
)
