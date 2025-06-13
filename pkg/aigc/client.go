package aigc

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"friedbot/pkg/config"
)

var client *Client

var host string

type Client struct {
	*http.Client

	apiKey string
}

func InitClient() error {
	var err error
	aiSettings := config.GetAISettings()
	if aiSettings.APIKey == "" {
		return errors.New("ai.api_key is empty")
	}
	host = aiSettings.Host
	_, err = url.Parse(host)
	if err != nil {
		return fmt.Errorf("parse host failed: %w", err)
	}
	var transport *http.Transport
	if aiSettings.Proxy != "" {
		proxy, err := url.Parse(aiSettings.Proxy)
		if err != nil {
			return fmt.Errorf("parse proxy failed: %w", err)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	} else {
		transport = &http.Transport{}
	}
	client = &Client{
		Client: &http.Client{
			Transport: transport,
		},
		apiKey: aiSettings.APIKey,
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

func GetCompletionReason(req *Request) (msg, reason string, err error) {
	req.Model = modelTypeDeepSeekReasoning
	return req.Post(pathTypeChatCompletions)
}

func GetStreamChat(req *Request) (*Stream, error) {
	req.Model = modelTypeDeepSeekChat
	req.Stream = true
	return req.PostStream(pathTypeChatCompletions)
}

func GetStreamReason(req *Request) (*Stream, error) {
	req.Model = modelTypeDeepSeekReasoning
	req.Stream = true
	return req.PostStream(pathTypeChatCompletions)
}
