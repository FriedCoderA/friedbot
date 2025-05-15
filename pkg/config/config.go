package config

import (
	"fmt"
	"log/slog"
	"os"

	"friedbot/pkg/xslog"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() error {
	workPath, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(workPath + "/config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config failed: %v, path=%s", err, workPath)
	}

	// 设置配置监听
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("config file changed")
		xslog.UpdateLogLevel() // 更新日志级别
	})
	return nil
}
