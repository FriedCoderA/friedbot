package aigc

type Tool struct {
	ID       string    `json:"id,omitempty"`
	Type     string    `json:"type,omitempty"`
	Function *Function `json:"function,omitempty"`
}

type Function struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Params      map[string]any `json:"parameters,omitempty"`
	Args        map[string]any `json:"arguments,omitempty"`
}
