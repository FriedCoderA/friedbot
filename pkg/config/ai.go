package config

import "github.com/spf13/viper"

type AISettings struct {
	APIKey string
	Host   string
	Proxy  string
}

func GetAISettings() *AISettings {
	return &AISettings{
		APIKey: viper.GetString("ai.api_key"),
		Host:   viper.GetString("ai.host"),
		Proxy:  viper.GetString("ai.proxy"),
	}
}
