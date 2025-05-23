package aigc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
)

type Stream struct {
	data chan string
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
			// 读取事件流数据
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					stream.data <- fmt.Sprintf("[ERROR] %v", err)
				}
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
			if !f(chunk) {
				return
			}
		}
	}
}
