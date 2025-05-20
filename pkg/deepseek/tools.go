package deepseek

type Tool struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function *Function
}

type Function struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Params      map[string]any `json:"parameters"`
	Args        map[string]any `json:"arguments"`
}
