package config

import "github.com/spf13/viper"

type BotSettings struct {
	GroupWhiteList []string
	UserBlackList  []string
}

func GetBotSettings() *BotSettings {
	return &BotSettings{
		GroupWhiteList: viper.GetStringSlice("bot.group_white_list"),
		UserBlackList:  viper.GetStringSlice("bot.user_black_list"),
	}
}
