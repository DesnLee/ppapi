package config

import (
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

func init() {
	LoadConfig()
}

func LoadConfig() {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(filename))
	// 读取配置文件
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(root)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
}
