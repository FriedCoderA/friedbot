package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func InitConfig() error {
	// 获取当前项目目录
	workPath, _ := os.Getwd()
	// 设置文件名和文件后缀
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	// 配置文件所在的文件夹
	viper.AddConfigPath(workPath + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read in config failed: %v, path=%s", err, workPath)
	}
	return nil
}
