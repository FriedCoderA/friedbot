package aigc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
)

const EOF = "[DONE]"

type Stream struct {
	closed bool
	data   chan string
}

func NewStream(body io.ReadCloser) *Stream {
	stream := &Stream{
		data: make(chan string),
	}
	go func() {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				slog.Warn("failed to close response body", "err", err)
			}
		}(body)

		reader := bufio.NewReader(body)
		for {
			if stream.closed {
				return
			}
			// 读取事件流数据
			line, err := reader.ReadBytes('\n')
			if err != nil {
				stream.data <- EOF
				return
			}
			// 解析SSE格式
			if bytes.HasPrefix(line, []byte("data: ")) {
				var chunk struct {
					Choices []struct {
						Delta struct {
							Content string `json:"content"`
						} `json:"delta"`
					} `json:"choices"`
				}
				if err := json.Unmarshal(line[6:], &chunk); err != nil {
					continue
				}
				if content := chunk.Choices[0].Delta.Content; content != "" {
					stream.data <- content
				}
			}
		}
	}()
	return stream
}

func (s *Stream) Close() {
	s.closed = true
	close(s.data)
}

func (s *Stream) Range(f func(string) bool) {
	for {
		select {
		case chunk, ok := <-s.data:
			if !ok {
				s.Close()
				return
			}
			if chunk == EOF {
				s.Close()
				return
			}
			if !f(chunk) {
				return
			}
		}
	}
}
