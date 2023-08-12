package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"ppapi.desnlee.com/cmd"
	"ppapi.desnlee.com/docs"
)

// @title           Pocket Purse API Docs
// @description     Pocket Purse API Docs with Swagger

// @contact.name   DesnLee
// @contact.url    https://desnlee.com
// @contact.email  jiakun.ui@gmail.com

// @host      localhost:9999
// @BasePath  /api/v1

// @securityDefinitions.basic  Bearer

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	defer log.Println("服务已关闭")

	docs.SwaggerInfo.Version = "1.0"

	// 读取配置文件
	readEnv()

	// 启动
	cmd.Run()
}

// 读取配置文件
func readEnv() {
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
}
