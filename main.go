package main

import (
	"log"

	"ppapi.desnlee.com/cmd"
)

func main() {
	defer log.Println("服务已关闭")

	// 关闭数据库连接
	// defer database.Close()

	// 初始化服务器
	cmd.RunServer()
}
