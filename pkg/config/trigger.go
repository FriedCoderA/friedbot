package config

import "github.com/spf13/viper"

type TriggerSettings struct {
	Temperature float64
}

func GetTriggerSettings() *TriggerSettings {
	return &TriggerSettings{
		Temperature: viper.GetFloat64("trigger.temperature"),
	}
}
