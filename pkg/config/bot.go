package config

import "github.com/spf13/viper"

const (
	botQQ       = "bot.qq"
	botPassword = "bot.password"
	botName     = "bot.name"
)

type BotSettings struct {
	QQ       int64
	Password string
	Name     string
}

func GetBotSettings() *BotSettings {
	return &BotSettings{
		QQ:       viper.GetInt64(botQQ),
		Password: viper.GetString(botPassword),
		Name:     viper.GetString(botName),
	}
}
