package onebot

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"friedbot/pkg/config"
)

type request struct {
	path string
	body any
}

func (r *request) Post() (*http.Response, error) {
	host := config.GetBotSettings().Address
	if host == "" {
		return nil, errors.New("config error: bot address is empty")
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	body, err := json.Marshal(r.body)
	if err != nil {
		return nil, err
	}
	return client.Post("http://"+host+r.path, "application/json", bytes.NewReader(body))
}
