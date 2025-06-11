package config

import "github.com/spf13/viper"

type BotSettings struct {
	QQ             int64
	Address        string
	GroupWhiteList []string
	UserBlackList  []string
}

func GetBotSettings() *BotSettings {
	return &BotSettings{
		QQ:             viper.GetInt64("bot.qq"),
		Address:        viper.GetString("bot.address"),
		GroupWhiteList: viper.GetStringSlice("bot.group_white_list"),
		UserBlackList:  viper.GetStringSlice("bot.user_black_list"),
	}
}
