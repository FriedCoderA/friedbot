package aigc

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

var client *Client

var host string

type Client struct {
	*http.Client

	apiKey string
}

func InitClient() error {
	apiKey := viper.GetString("ai.api_key")
	if apiKey == "" {
		return errors.New("ai.api_key is empty")
	}
	host = viper.GetString("ai.host")
	if host == "" {
		return errors.New("ai.host is empty")
	}
	client = &Client{
		Client: http.DefaultClient,
		apiKey: apiKey,
	}
	return nil
}

func (c *Client) Send(path string, body []byte, stream bool) (*http.Response, error) {
	req, err := http.NewRequest(
		"POST",
		host+path,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	}
	return c.Client.Do(req)
}

func GetCompletionChat(req *Request) (string, error) {
	req.Model = modelTypeDeepSeekChat
	msg, _, err := req.Post(pathTypeChatCompletions)
	return msg, err
}

func GetCompletionReason(req *Request) (msg, reasoning string, err error) {
	req.Model = modelTypeDeepSeekReasoning
	return req.Post(pathTypeChatCompletions)
}

func GetStreamChat(req *Request) (*Stream, error) {
	req.Model = modelTypeDeepSeekChat
	req.Stream = true
	return req.PostStream(pathTypeChatCompletions)
}
