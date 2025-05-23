package aigc

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	Messages           Messages           `json:"messages"`
	Model              modelType          `json:"model"`
	FrequencyPenalty   float32            `json:"frequency_penalty,omitempty"`
	MaxTokens          int                `json:"max_tokens,omitempty"`
	PresencePenalty    float32            `json:"presence_penalty,omitempty"`
	ResponseFormatType responseFormatType `json:"response_format_type,omitempty"`
	Stop               []string           `json:"stop,omitempty"`
	Stream             bool               `json:"stream,omitempty"`
	Temperature        float32            `json:"temperature,omitempty"`
	TopP               float32            `json:"top_p,omitempty"`
	Tools              []*ToolMessage     `json:"tools,omitempty"`
	ToolChoice         toolChoiceType     `json:"tool_choice,omitempty"`
}

func (r *Request) Post(path string) (msg, reasoning string, err error) {
	body, err := json.Marshal(r)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %v", err)
	}
	resp, err := client.Send(path, body, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %v", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}
	var res *Response
	if err := json.Unmarshal(data, &res); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if len(res.Choices) == 0 {
		return "", "", fmt.Errorf("no choices returned")
	}
	return res.Choices[0].Message.Content, res.Choices[0].Message.ReasoningContent, nil
}

func (r *Request) PostStream(path string) (*Stream, error) {
	// 创建管道

	// 序列化请求
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 发送请求
	resp, err := client.Send(path, body, true)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// 处理非200响应
	if resp.StatusCode != http.StatusOK {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				slog.Warn("failed to close response body", "err", err)
			}
		}(resp.Body)
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, body)
	}

	stream := NewStream(resp.Body)
	return stream, nil
}
