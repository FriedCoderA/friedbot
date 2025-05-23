package aigc

type Response struct {
	ID        string   `json:"id"`
	Choices   []Choice `json:"choices"`
	CreatedAt int      `json:"created"`
}

type Choice struct {
	FinishReason string           `json:"finish_reason"`
	Index        int              `json:"index"`
	Message      *ResponseMessage `json:"message"`
}
