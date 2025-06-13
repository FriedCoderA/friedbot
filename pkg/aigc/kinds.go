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
	modelTypeDeepSeekChat         = "deepseek-chat"
	modelTypeDeepSeekReasoning    = "deepseek-reasoner"
	modelType70bDeepSeekReasoning = "deepseek-r1-distill-llama-70b"
)

type responseFormatType = any

const (
	ResponseFormatTypeText = "text"
)

var ResponseFormatTypeJSON = struct {
	Type string `json:"type"`
}{Type: "json_object"}

type toolChoiceType = string

const (
	ToolChoiceTypeNone     = "none"
	ToolChoiceTypeAuto     = "auto"
	ToolChoiceTypeRequired = "required"
)

const (
	pathTypeChatCompletions = "/v1/chat/completions"
)
