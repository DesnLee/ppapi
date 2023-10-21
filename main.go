package main

import (
	"log"

	"ppapi.desnlee.com/cmd"
	_ "ppapi.desnlee.com/config"
	"ppapi.desnlee.com/docs"
	"ppapi.desnlee.com/internal/database"
)

//	@title			Pocket Purse API Docs
//	@description	Pocket Purse API Docs with Swagger

//	@contact.name	DesnLee
//	@contact.url	https://desnlee.com
//	@contact.email	jiakun.ui@gmail.com

//	@host		localhost:9999
//	@BasePath	/api/v1

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	defer log.Println("服务已关闭")
	defer database.Close()

	docs.SwaggerInfo.Version = "1.0"

	// 读取配置文件
	// config.LoadConfig()

	// 启动
	cmd.Run()
}
